package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

type Hotel struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

var hotels = []Hotel{
	{ID: "H1", Name: "Grand Hotel Berlin", Price: 150, Currency: "EUR", Unit: "night"},
	{ID: "H2", Name: "Seaside Resort Hamburg", Price: 200, Currency: "EUR", Unit: "night"},
	{ID: "H3", Name: "City Lodge München", Price: 95, Currency: "EUR", Unit: "night"},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /hotels", hotelsHandler)
	mux.HandleFunc("GET /openapi", openapiHandler)

	log.Println("HotelService starting on port 8082...")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func hotelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

func openapiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(openapiSpec)
}
