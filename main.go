package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type PingPongResponse struct {
	Ping string `json:"ping"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ShortifyRequest struct {
	Url string `json:"url"`
}

type ShortifyResponse struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

func PingPongHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		resp := ErrorResponse{"Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := PingPongResponse{"pong"}

	if err := db.Raw("SELECT 1").Error; err != nil {
		log.Println("DB Connection failed:", err)
		return
	}
	log.Println("DB Connection successful")

	json.NewEncoder(w).Encode(resp)
}

func ShortifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		resp := ErrorResponse{"Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req ShortifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := ErrorResponse{"Invalid request"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !isValidURL(req.Url) {
		resp := ErrorResponse{"Invalid url"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	url := strings.TrimSuffix(req.Url, "/")
	code := GenerateCode(url)
	_, err := db.GetOrCreateUrl(code, url)

	if err != nil {
		log.Println(err)
		resp := ErrorResponse{"Internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	shortUrl := GenerateShortUrl(code)
	resp := ShortifyResponse{url, shortUrl}
	json.NewEncoder(w).Encode(resp)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	secPurpose := r.Header.Get("Sec-Purpose")

	if strings.Contains(secPurpose, "prefetch") {
		log.Printf("Ignoring prefetch/prerender request: %s", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		resp := ErrorResponse{"Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")

	if len(code) != 8 {
		resp := ErrorResponse{"Invalid url"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	url, err := db.GetUrl(code)
	if err != nil {
		resp := ErrorResponse{"Not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Redirect(w, r, *url, http.StatusFound)
}

var db *DB

func main() {
	var err error
	db, err = NewDB()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/ping", PingPongHandler)
	http.HandleFunc("/shortify", ShortifyHandler)
	http.HandleFunc("/", RedirectHandler)
	http.ListenAndServe(":8080", nil)
}
