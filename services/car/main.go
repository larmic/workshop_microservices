package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

type Car struct {
	ID       string  `json:"id"`
	Model    string  `json:"model"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

var cars = []Car{
	{ID: "C1", Model: "VW Golf", Price: 45, Currency: "EUR", Unit: "day"},
	{ID: "C2", Model: "BMW 3er", Price: 85, Currency: "EUR", Unit: "day"},
	{ID: "C3", Model: "Mercedes E-Klasse", Price: 120, Currency: "EUR", Unit: "day"},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /cars", carsHandler)
	mux.HandleFunc("GET /openapi", openapiHandler)

	log.Println("CarService starting on port 8083...")
	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func carsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func openapiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(openapiSpec)
}
