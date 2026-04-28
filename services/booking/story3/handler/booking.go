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

func BookingOffersHandler(resolver *consul.Resolver, client *http.Client, cb *circuitbreaker.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		hotelURL, err := resolver.ResolveServiceURL("hotel-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve hotel-service: %v", err), http.StatusInternalServerError)
			return
		}

		carURL, err := resolver.ResolveServiceURL("car-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve car-service: %v", err), http.StatusInternalServerError)
			return
		}

		var flights json.RawMessage
		flightErr := cb.Execute(r.Context(), func(ctx context.Context) error {
			flightURL, resolveErr := resolver.ResolveServiceURL("flight-service")
			if resolveErr != nil {
				return resolveErr
			}
			data, fetchErr := fetchJSON(ctx, client, fmt.Sprintf("%s/flights", flightURL))
			if fetchErr != nil {
				return fetchErr
			}
			flights = data
			return nil
		})
		if flightErr != nil {
			flights = emptyJSONArray
			markFallback(w, flightErr)
		}

		hotels, err := fetchJSON(r.Context(), client, fmt.Sprintf("%s/hotels", hotelURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch hotels: %v", err), http.StatusInternalServerError)
			return
		}

		cars, err := fetchJSON(r.Context(), client, fmt.Sprintf("%s/cars", carURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch cars: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		})
	}
}

func CreateBookingHandler(resolver *consul.Resolver, client *http.Client, cb *circuitbreaker.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		hotelURL, err := resolver.ResolveServiceURL("hotel-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve hotel-service: %v", err), http.StatusInternalServerError)
			return
		}

		carURL, err := resolver.ResolveServiceURL("car-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve car-service: %v", err), http.StatusInternalServerError)
			return
		}

		var flight json.RawMessage
		flightErr := cb.Execute(r.Context(), func(ctx context.Context) error {
			flightURL, resolveErr := resolver.ResolveServiceURL("flight-service")
			if resolveErr != nil {
				return resolveErr
			}
			data, postErr := postJSON(ctx, client, fmt.Sprintf("%s/bookings", flightURL),
				map[string]string{"flightId": req.FlightID, "customerName": req.CustomerName})
			if postErr != nil {
				return postErr
			}
			flight = data
			return nil
		})
		if flightErr != nil {
			flight = nil
			markFallback(w, flightErr)
		}

		hotel, err := postJSON(r.Context(), client, fmt.Sprintf("%s/bookings", hotelURL),
			map[string]string{"hotelId": req.HotelID, "customerName": req.CustomerName})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to book hotel: %v", err), http.StatusInternalServerError)
			return
		}

		car, err := postJSON(r.Context(), client, fmt.Sprintf("%s/bookings", carURL),
			map[string]string{"carId": req.CarID, "customerName": req.CustomerName})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to book car: %v", err), http.StatusInternalServerError)
			return
		}

		booking := Booking{
			BookingID:    newBookingID(),
			CustomerName: req.CustomerName,
			Flight:       flight,
			Hotel:        hotel,
			Car:          car,
		}

		log.Printf("Aggregated booking confirmed: bookingId=%s customer=%q flightIncluded=%t",
			booking.BookingID, booking.CustomerName, flight != nil)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(booking)
	}
}

func markFallback(w http.ResponseWriter, err error) {
	w.Header().Set("X-Flight-Fallback", "true")
	if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
		w.Header().Set("X-Circuit-Open", "flight")
		log.Printf("flight call short-circuited (CB OPEN)")
	} else {
		log.Printf("flight call failed, applying fallback: %v", err)
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
