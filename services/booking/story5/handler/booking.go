package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/booking/story5/bulkhead"
	"github.com/team-neusta-skills/workshop_microservices/booking/story5/circuitbreaker"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type Config struct {
	Service   string `json:"service"`
	ConsulURL string `json:"consulUrl"`
	Timeout   int    `json:"timeout"`
}

type Breakers struct {
	Flight *circuitbreaker.CircuitBreaker
	Hotel  *circuitbreaker.CircuitBreaker
	Car    *circuitbreaker.CircuitBreaker
}

func (b Breakers) All() []*circuitbreaker.CircuitBreaker {
	return []*circuitbreaker.CircuitBreaker{b.Flight, b.Hotel, b.Car}
}

type Bulkheads struct {
	Flight *bulkhead.Bulkhead
	Hotel  *bulkhead.Bulkhead
	Car    *bulkhead.Bulkhead
}

func (b Bulkheads) All() []*bulkhead.Bulkhead {
	return []*bulkhead.Bulkhead{b.Flight, b.Hotel, b.Car}
}

type BookingOffers struct {
	Flights json.RawMessage `json:"flights"`
	Hotels  json.RawMessage `json:"hotels"`
	Cars    json.RawMessage `json:"cars"`
}

type BookingRequest struct {
	FlightID     string `json:"flightId"`
	HotelID      string `json:"hotelId"`
	CarID        string `json:"carId"`
	CustomerName string `json:"customerName"`
}

type Booking struct {
	BookingID    string          `json:"bookingId"`
	CustomerName string          `json:"customerName"`
	Flight       json.RawMessage `json:"flight"`
	Hotel        json.RawMessage `json:"hotel"`
	Car          json.RawMessage `json:"car"`
}

var emptyJSONArray = json.RawMessage("[]")

func BookingOffersHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers, bulkheads Bulkheads) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		flights := fetchOffersWithBHCB(r.Context(), w, bulkheads.Flight, breakers.Flight, "flight",
			resolver, client, "flight-service", "/flights")
		hotels := fetchOffersWithBHCB(r.Context(), w, bulkheads.Hotel, breakers.Hotel, "hotel",
			resolver, client, "hotel-service", "/hotels")
		cars := fetchOffersWithBHCB(r.Context(), w, bulkheads.Car, breakers.Car, "car",
			resolver, client, "car-service", "/cars")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		})
	}
}

// CreateBookingHandler bucht alles oder nichts (Best-Effort) — bei einer
// Teilbuchung sähe der Kunde sonst eine inkonsistente Bestätigung. Echte
// atomare Buchung über mehrere Services wäre Saga (Story 5).
func CreateBookingHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Pre-Check: wenn ein CB OPEN ist, gar nicht erst anfangen zu buchen.
		// HALF_OPEN ist OK — das ist der Probe-Versuch, der zur Recovery führt.
		for _, cb := range breakers.All() {
			snap := cb.Snapshot()
			if snap.State == "OPEN" {
				writeBookingFailure(w, snap.Name, errCircuitOpenForBooking, nil)
				return
			}
		}

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var booked []string

		flight, err := bookSingle(r.Context(), breakers.Flight,
			resolver, client, "flight-service",
			map[string]string{"flightId": req.FlightID, "customerName": req.CustomerName})
		if err != nil {
			writeBookingFailure(w, "flight", err, booked)
			return
		}
		booked = append(booked, "flight")

		hotel, err := bookSingle(r.Context(), breakers.Hotel,
			resolver, client, "hotel-service",
			map[string]string{"hotelId": req.HotelID, "customerName": req.CustomerName})
		if err != nil {
			writeBookingFailure(w, "hotel", err, booked)
			return
		}
		booked = append(booked, "hotel")

		car, err := bookSingle(r.Context(), breakers.Car,
			resolver, client, "car-service",
			map[string]string{"carId": req.CarID, "customerName": req.CustomerName})
		if err != nil {
			writeBookingFailure(w, "car", err, booked)
			return
		}

		booking := Booking{
			BookingID:    newBookingID(),
			CustomerName: req.CustomerName,
			Flight:       flight,
			Hotel:        hotel,
			Car:          car,
		}

		log.Printf("Aggregated booking confirmed: bookingId=%s customer=%q",
			booking.BookingID, booking.CustomerName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(booking)
	}
}

var errCircuitOpenForBooking = errors.New("circuit breaker is open")

type BookingFailure struct {
	Error            string   `json:"error"`
	FailedService    string   `json:"failedService"`
	Reason           string   `json:"reason"`
	PreviouslyBooked []string `json:"previouslyBooked,omitempty"`
	Hint             string   `json:"hint"`
}

func writeBookingFailure(w http.ResponseWriter, service string, err error, previouslyBooked []string) {
	failure := BookingFailure{
		Error:            "booking aborted (fail-fast, no partial booking)",
		FailedService:    service,
		Reason:           err.Error(),
		PreviouslyBooked: previouslyBooked,
		Hint:             "Echte atomare Buchung über mehrere Services erfordert Compensation/Saga (Story 5).",
	}
	if len(previouslyBooked) > 0 {
		log.Printf("Booking ABORTED at %s (err=%v) — bereits gebucht und NICHT zurückgerollt: %v (Saga in Story 5)",
			service, err, previouslyBooked)
	} else {
		log.Printf("Booking ABORTED at %s (err=%v) — keine Buchung durchgeführt", service, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(failure)
}

func bookSingle(
	ctx context.Context,
	cb *circuitbreaker.CircuitBreaker,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	payload map[string]string,
) (json.RawMessage, error) {
	var data json.RawMessage
	err := cb.Execute(ctx, func(ctx context.Context) error {
		url, resolveErr := resolver.ResolveServiceURL(consulName)
		if resolveErr != nil {
			return resolveErr
		}
		raw, postErr := postJSON(ctx, client, url+"/bookings", payload)
		if postErr != nil {
			return postErr
		}
		data = raw
		return nil
	})
	return data, err
}

// fetchOffersWithBHCB legt den Bulkhead VOR den Circuit Breaker: ist der Pool
// voll, kürzen wir sofort ab — ohne die CB-Statistik zu vergiften, weil der
// Backend-Call ja gar nicht stattgefunden hat. So bleiben CB und Bulkhead in
// ihren Zuständigkeiten getrennt: CB misst Backend-Gesundheit, Bulkhead schützt
// vor Ressourcen-Erschöpfung.
func fetchOffersWithBHCB(
	ctx context.Context,
	w http.ResponseWriter,
	bh *bulkhead.Bulkhead,
	cb *circuitbreaker.CircuitBreaker,
	serviceLabel string,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	path string,
) json.RawMessage {
	var data json.RawMessage
	err := bh.Execute(ctx, func(ctx context.Context) error {
		return cb.Execute(ctx, func(ctx context.Context) error {
			url, resolveErr := resolver.ResolveServiceURL(consulName)
			if resolveErr != nil {
				return resolveErr
			}
			raw, fetchErr := fetchJSON(ctx, client, url+path)
			if fetchErr != nil {
				return fetchErr
			}
			data = raw
			return nil
		})
	})
	if err != nil {
		markFallback(w, serviceLabel, err)
		return emptyJSONArray
	}
	return data
}

func markFallback(w http.ResponseWriter, serviceLabel string, err error) {
	w.Header().Add("X-Fallback", serviceLabel)
	switch {
	case errors.Is(err, bulkhead.ErrBulkheadFull):
		w.Header().Add("X-Bulkhead-Full", serviceLabel)
		log.Printf("%s call rejected (bulkhead full)", serviceLabel)
	case errors.Is(err, circuitbreaker.ErrCircuitOpen):
		w.Header().Add("X-Circuit-Open", serviceLabel)
		log.Printf("%s call short-circuited (CB OPEN)", serviceLabel)
	default:
		log.Printf("%s call failed, applying fallback: %v", serviceLabel, err)
	}
}

func fetchJSON(ctx context.Context, client *http.Client, url string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("backend returned %d: %s", resp.StatusCode, string(body))
	}

	return json.RawMessage(body), nil
}

func postJSON(ctx context.Context, client *http.Client, url string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("backend returned %d: %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

func newBookingID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "B-" + hex.EncodeToString(b)
}
