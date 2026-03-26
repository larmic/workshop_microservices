package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Car struct {
	ID       string  `json:"id"`
	Model    string  `json:"model"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

var cars = []Car{
	{ID: "C1", Model: "Ford Mustang (New York)", Price: 95, Currency: "USD", Unit: "day"},
	{ID: "C2", Model: "Mini Cooper (London)", Price: 65, Currency: "GBP", Unit: "day"},
	{ID: "C3", Model: "Renault Clio (Paris)", Price: 55, Currency: "EUR", Unit: "day"},
}

func CarsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}
