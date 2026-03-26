package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Flight struct {
	ID          string  `json:"id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
}

var flights = []Flight{
	{ID: "LH123", Origin: "Frankfurt", Destination: "New York", Price: 450, Currency: "EUR"},
	{ID: "LH456", Origin: "München", Destination: "London", Price: 180, Currency: "EUR"},
	{ID: "BA789", Origin: "Berlin", Destination: "Paris", Price: 120, Currency: "EUR"},
}

func FlightsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
}
