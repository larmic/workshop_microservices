package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/booking/story2/consul"
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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func InfoHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}
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
		json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		})
	}
}

func OpenapiHandler(spec []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
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

	return json.RawMessage(body), nil
}
