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

func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)

	id := r.PathValue("id")
	log.Printf("[%s] Flight booking cancelled: bookingId=%s", hostname, id)

	w.WriteHeader(http.StatusNoContent)
}

type CompensationEvent struct {
	EventID   string `json:"eventId"`
	SagaID    string `json:"sagaId"`
	BookingID string `json:"bookingId"`
}

// CompensationEventHandler nimmt ein CompensationRequested-Event entgegen,
// antwortet sofort mit 202 Accepted und führt die Stornierung asynchron
// in einer Goroutine aus (fire & forget aus Sicht des Senders).
func CompensationEventHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)

	var ev CompensationEvent
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "invalid event body", http.StatusBadRequest)
		return
	}
	if ev.EventID == "" || ev.SagaID == "" || ev.BookingID == "" {
		http.Error(w, "eventId, sagaId and bookingId are required", http.StatusBadRequest)
		return
	}

	log.Printf("[%s] event=CompensationRequested service=flight phase=received   eventId=%s sagaId=%s bookingId=%s",
		hostname, ev.EventID, ev.SagaID, ev.BookingID)

	w.WriteHeader(http.StatusAccepted)

	go func() {
		log.Printf("[%s] event=CompensationRequested service=flight phase=processing eventId=%s sagaId=%s bookingId=%s",
			hostname, ev.EventID, ev.SagaID, ev.BookingID)
		log.Printf("[%s] event=CompensationRequested service=flight phase=done       eventId=%s sagaId=%s bookingId=%s",
			hostname, ev.EventID, ev.SagaID, ev.BookingID)
	}()
}

func newBookingID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
