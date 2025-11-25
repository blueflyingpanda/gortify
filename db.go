package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type URL struct {
	Code      string    `gorm:"primaryKey;type:varchar(8)"`
	URL       string    `gorm:"not null;type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Clicks    int64     `gorm:"default:0"`
}

// DB wraps the GORM database connection
type DB struct {
	*gorm.DB
}

func (db *DB) GetOrCreateUrl(code string, url string) (*URL, error) {
	var urlRecord URL
	result := db.Where(URL{Code: code}).FirstOrCreate(&urlRecord, URL{Code: code, URL: url})
	if result.Error != nil {
		return nil, result.Error
	}
	return &urlRecord, nil
}

func (db *DB) GetUrl(code string) (*string, error) {
	var urlRecord URL
	result := db.Where(URL{Code: code}).First(&urlRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	// UpdateColumn does not update updatedAt and is a bit faster than Update
	go func() {
		updResult := db.Model(&urlRecord).UpdateColumn("clicks", gorm.Expr("clicks + 1"))
		if updResult.Error != nil {
			log.Println("Error updating clicks for code:", code, result.Error)
		}
	}()

	return &urlRecord.URL, nil
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Conf.DBHost, Conf.DBPort, Conf.DBUser, Conf.DBPass, Conf.DBName, Conf.DBSslMode)
}

// NewDB creates a new database connection and migrates tables
func NewDB() (*DB, error) {
	connectionString := getConnectionString()
	db, connectErr := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if connectErr != nil {
		panic(connectErr)
	}

	migrateErr := db.AutoMigrate(&URL{})
	if migrateErr != nil {
		panic(migrateErr)
	}

	return &DB{db}, nil
}
