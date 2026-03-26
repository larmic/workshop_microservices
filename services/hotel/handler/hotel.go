package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Hotel struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

var hotels = []Hotel{
	{ID: "H1", Name: "Manhattan Plaza Hotel", Price: 250, Currency: "USD", Unit: "night"},
	{ID: "H2", Name: "London Bridge Inn", Price: 180, Currency: "GBP", Unit: "night"},
	{ID: "H3", Name: "Paris Étoile Hotel", Price: 195, Currency: "EUR", Unit: "night"},
}

func HotelsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}
