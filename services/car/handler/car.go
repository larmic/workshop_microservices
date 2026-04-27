package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Car struct {
	ID       string  `json:"id"`
	Model    string  `json:"model"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

type BookingRequest struct {
	CarID        string `json:"carId"`
	CustomerName string `json:"customerName"`
}

type Booking struct {
	BookingID    string `json:"bookingId"`
	CarID        string `json:"carId"`
	CustomerName string `json:"customerName"`
	Status       string `json:"status"`
}

var cars = []Car{
	{ID: "C1", Model: "Ford Mustang (New York)", Price: 95, Currency: "USD", Unit: "day"},
	{ID: "C2", Model: "Mini Cooper (London)", Price: 65, Currency: "GBP", Unit: "day"},
	{ID: "C3", Model: "Renault Clio (Paris)", Price: 55, Currency: "EUR", Unit: "day"},
}

func CarsHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	booking := Booking{
		BookingID:    newBookingID("C"),
		CarID:        req.CarID,
		CustomerName: req.CustomerName,
		Status:       "CONFIRMED",
	}

	log.Printf("[%s] Car booking confirmed: bookingId=%s carId=%s customer=%q",
		hostname, booking.BookingID, booking.CarID, booking.CustomerName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func newBookingID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
