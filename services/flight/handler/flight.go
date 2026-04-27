package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Flight struct {
	ID          string  `json:"id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
}

type BookingRequest struct {
	FlightID     string `json:"flightId"`
	CustomerName string `json:"customerName"`
}

type Booking struct {
	BookingID    string `json:"bookingId"`
	FlightID     string `json:"flightId"`
	CustomerName string `json:"customerName"`
	Status       string `json:"status"`
}

var flights = []Flight{
	{ID: "LH123", Origin: "Frankfurt", Destination: "New York", Price: 450, Currency: "EUR"},
	{ID: "LH456", Origin: "München", Destination: "London", Price: 180, Currency: "EUR"},
	{ID: "BA789", Origin: "Berlin", Destination: "Paris", Price: 120, Currency: "EUR"},
}

func FlightsHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
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
		BookingID:    newBookingID("F"),
		FlightID:     req.FlightID,
		CustomerName: req.CustomerName,
		Status:       "CONFIRMED",
	}

	log.Printf("[%s] Flight booking confirmed: bookingId=%s flightId=%s customer=%q",
		hostname, booking.BookingID, booking.FlightID, booking.CustomerName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func newBookingID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
