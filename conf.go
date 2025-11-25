package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSslMode string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config := Config{
		BaseUrl:   os.Getenv("BASE_URL"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBSslMode: os.Getenv("DB_SSL_MODE"),
	}
	return &config
}

var Conf = LoadConfig()
