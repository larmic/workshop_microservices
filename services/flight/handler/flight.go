package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/team-neusta-skills/workshop_microservices/shared/tracing"
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
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
}

func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

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

	tracing.Logger(r.Context()).Info("flight booking confirmed",
		"host", hostname(),
		"bookingId", booking.BookingID,
		"flightId", booking.FlightID,
		"customer", booking.CustomerName,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	id := r.PathValue("id")
	tracing.Logger(r.Context()).Info("flight booking cancelled",
		"host", hostname(),
		"bookingId", id,
	)

	w.WriteHeader(http.StatusNoContent)
}

type CompensationEvent struct {
	EventID     string `json:"eventId"`
	SagaID      string `json:"sagaId"`
	BookingID   string `json:"bookingId"`
	Traceparent string `json:"traceparent,omitempty"`
}

// CompensationEventHandler nimmt ein CompensationRequested-Event entgegen,
// antwortet sofort mit 202 Accepted und führt die Stornierung asynchron
// in einer Goroutine aus (fire & forget aus Sicht des Senders).
//
// Trace-Kontext wandert über die `traceparent`-Property im Event-Body
// mit (Async-Grenze) und wird vor dem Goroutine-Start in den Kontext gelegt,
// damit die Stornierungs-Logzeilen dieselbe trace_id tragen wie die
// ursprüngliche Buchung.
func CompensationEventHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var ev CompensationEvent
	if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
		http.Error(w, "invalid event body", http.StatusBadRequest)
		return
	}
	if ev.EventID == "" || ev.SagaID == "" || ev.BookingID == "" {
		http.Error(w, "eventId, sagaId and bookingId are required", http.StatusBadRequest)
		return
	}

	asyncCtx := contextFromEvent(ev.Traceparent)
	logger := tracing.Logger(asyncCtx).With(
		"host", hostname(),
		"event", "CompensationRequested",
		"service", "flight",
		"eventId", ev.EventID,
		"sagaId", ev.SagaID,
		"bookingId", ev.BookingID,
	)

	logger.Info("compensation event", "phase", "received")

	w.WriteHeader(http.StatusAccepted)

	go func() {
		logger.Info("compensation event", "phase", "processing")
		logger.Info("compensation event", "phase", "done")
	}()
}

func logRequest(r *http.Request) {
	tracing.Logger(r.Context()).Info("request",
		"host", hostname(),
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}

// contextFromEvent baut einen Kontext für die asynchrone Verarbeitung. Wenn
// das Event eine `traceparent`-Property mitbringt, wird die Trace-ID
// übernommen — sonst wird ein frischer Kontext erzeugt, damit auch dann
// strukturierte Korrelation möglich ist (mit unterscheidbarer ID, nicht
// stillschweigend ohne).
func contextFromEvent(traceparent string) context.Context {
	tc, ok := tracing.Parse(traceparent)
	if !ok {
		tc = tracing.New()
	}
	return tracing.WithContext(context.Background(), tc)
}

func newBookingID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
