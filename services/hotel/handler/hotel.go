package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Hotel struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit"`
}

type BookingRequest struct {
	HotelID      string `json:"hotelId"`
	CustomerName string `json:"customerName"`
}

type Booking struct {
	BookingID    string `json:"bookingId"`
	HotelID      string `json:"hotelId"`
	CustomerName string `json:"customerName"`
	Status       string `json:"status"`
}

var hotels = []Hotel{
	{ID: "H1", Name: "Manhattan Plaza Hotel", Price: 250, Currency: "USD", Unit: "night"},
	{ID: "H2", Name: "London Bridge Inn", Price: 180, Currency: "GBP", Unit: "night"},
	{ID: "H3", Name: "Paris Étoile Hotel", Price: 195, Currency: "EUR", Unit: "night"},
}

func HotelsHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
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
		BookingID:    newBookingID("H"),
		HotelID:      req.HotelID,
		CustomerName: req.CustomerName,
		Status:       "CONFIRMED",
	}

	log.Printf("[%s] Hotel booking confirmed: bookingId=%s hotelId=%s customer=%q",
		hostname, booking.BookingID, booking.HotelID, booking.CustomerName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	log.Printf("[%s] %s %s from %s", hostname, r.Method, r.URL.Path, r.RemoteAddr)

	id := r.PathValue("id")
	log.Printf("[%s] Hotel booking cancelled: bookingId=%s", hostname, id)

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

	log.Printf("[%s] event=CompensationRequested service=hotel phase=received   eventId=%s sagaId=%s bookingId=%s",
		hostname, ev.EventID, ev.SagaID, ev.BookingID)

	w.WriteHeader(http.StatusAccepted)

	go func() {
		log.Printf("[%s] event=CompensationRequested service=hotel phase=processing eventId=%s sagaId=%s bookingId=%s",
			hostname, ev.EventID, ev.SagaID, ev.BookingID)
		log.Printf("[%s] event=CompensationRequested service=hotel phase=done       eventId=%s sagaId=%s bookingId=%s",
			hostname, ev.EventID, ev.SagaID, ev.BookingID)
	}()
}

func newBookingID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
