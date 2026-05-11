package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
)

type Config struct {
	Service   string `json:"service"`
	ConsulURL string `json:"consulUrl"`
	Timeout   int    `json:"timeout"`
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

func BookingOffersHandler(resolver *consul.Resolver, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		flightURL, err := resolver.ResolveServiceURL("flight-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve flight-service: %v", err), http.StatusInternalServerError)
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

		flights, err := fetchJSON(client, fmt.Sprintf("%s/flights", flightURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch flights: %v", err), http.StatusInternalServerError)
			return
		}

		hotels, err := fetchJSON(client, fmt.Sprintf("%s/hotels", hotelURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch hotels: %v", err), http.StatusInternalServerError)
			return
		}

		cars, err := fetchJSON(client, fmt.Sprintf("%s/cars", carURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch cars: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		}); err != nil {
			log.Printf("encode booking offers failed: %v", err)
		}
	}
}

func fetchJSON(client *http.Client, url string) (json.RawMessage, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream %s returned HTTP %d: %s",
			url, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.RawMessage(body), nil
}

func CreateBookingHandler(resolver *consul.Resolver, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		flightURL, err := resolver.ResolveServiceURL("flight-service")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to resolve flight-service: %v", err), http.StatusInternalServerError)
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

		flight, err := postJSON(client, fmt.Sprintf("%s/bookings", flightURL),
			map[string]string{"flightId": req.FlightID, "customerName": req.CustomerName})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to book flight: %v", err), http.StatusInternalServerError)
			return
		}

		hotel, err := postJSON(client, fmt.Sprintf("%s/bookings", hotelURL),
			map[string]string{"hotelId": req.HotelID, "customerName": req.CustomerName})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to book hotel: %v", err), http.StatusInternalServerError)
			return
		}

		car, err := postJSON(client, fmt.Sprintf("%s/bookings", carURL),
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

		log.Printf("Aggregated booking confirmed: bookingId=%s customer=%q",
			booking.BookingID, booking.CustomerName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(booking)
	}
}

func postJSON(client *http.Client, url string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
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
	rand.Read(b)
	return "B-" + hex.EncodeToString(b)
}
