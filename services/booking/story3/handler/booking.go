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

	"github.com/team-neusta-skills/workshop_microservices/booking/story3/circuitbreaker"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type Config struct {
	ConsulURL     string `json:"consulUrl"`
	Timeout       int    `json:"timeout"`
	ActiveProfile string `json:"activeProfile"`
}

type Breakers struct {
	Flight *circuitbreaker.CircuitBreaker
	Hotel  *circuitbreaker.CircuitBreaker
	Car    *circuitbreaker.CircuitBreaker
}

func (b Breakers) All() []*circuitbreaker.CircuitBreaker {
	return []*circuitbreaker.CircuitBreaker{b.Flight, b.Hotel, b.Car}
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

func BookingOffersHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		flights := fetchOffersWithCB(r.Context(), w, breakers.Flight, "flight",
			resolver, client, "flight-service", "/flights")
		hotels := fetchOffersWithCB(r.Context(), w, breakers.Hotel, "hotel",
			resolver, client, "hotel-service", "/hotels")
		cars := fetchOffersWithCB(r.Context(), w, breakers.Car, "car",
			resolver, client, "car-service", "/cars")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		})
	}
}

func CreateBookingHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		flight := bookWithCB(r.Context(), w, breakers.Flight, "flight",
			resolver, client, "flight-service",
			map[string]string{"flightId": req.FlightID, "customerName": req.CustomerName})
		hotel := bookWithCB(r.Context(), w, breakers.Hotel, "hotel",
			resolver, client, "hotel-service",
			map[string]string{"hotelId": req.HotelID, "customerName": req.CustomerName})
		car := bookWithCB(r.Context(), w, breakers.Car, "car",
			resolver, client, "car-service",
			map[string]string{"carId": req.CarID, "customerName": req.CustomerName})

		booking := Booking{
			BookingID:    newBookingID(),
			CustomerName: req.CustomerName,
			Flight:       flight,
			Hotel:        hotel,
			Car:          car,
		}

		log.Printf("Aggregated booking confirmed: bookingId=%s customer=%q flight=%t hotel=%t car=%t",
			booking.BookingID, booking.CustomerName,
			flight != nil, hotel != nil, car != nil)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(booking)
	}
}

func fetchOffersWithCB(
	ctx context.Context,
	w http.ResponseWriter,
	cb *circuitbreaker.CircuitBreaker,
	serviceLabel string,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	path string,
) json.RawMessage {
	var data json.RawMessage
	err := cb.Execute(ctx, func(ctx context.Context) error {
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
	if err != nil {
		markFallback(w, serviceLabel, err)
		return emptyJSONArray
	}
	return data
}

func bookWithCB(
	ctx context.Context,
	w http.ResponseWriter,
	cb *circuitbreaker.CircuitBreaker,
	serviceLabel string,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	payload map[string]string,
) json.RawMessage {
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
	if err != nil {
		markFallback(w, serviceLabel, err)
		return nil
	}
	return data
}

func markFallback(w http.ResponseWriter, serviceLabel string, err error) {
	w.Header().Add("X-Fallback", serviceLabel)
	if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
		w.Header().Add("X-Circuit-Open", serviceLabel)
		log.Printf("%s call short-circuited (CB OPEN)", serviceLabel)
	} else {
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
