package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

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

type consulServiceEntry struct {
	Service struct {
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

var config Config
var httpClient *http.Client

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func resolveServiceURL(serviceName string) (string, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s?passing=true", config.ConsulURL, serviceName)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("consul request failed: %w", err)
	}
	defer resp.Body.Close()

	var entries []consulServiceEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return "", fmt.Errorf("consul response parse failed: %w", err)
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("no healthy instances found for %s", serviceName)
	}

	entry := entries[rand.Intn(len(entries))]
	return fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port), nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	config = Config{
		ConsulURL:     getEnv("CONSUL_URL", "http://localhost:8500"),
		Timeout:       5000,
		ActiveProfile: getEnv("ACTIVE_PROFILE", "dev"),
	}

	httpClient = &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /info", infoHandler)
	mux.HandleFunc("GET /booking/offers", bookingOffersHandler)
	mux.HandleFunc("GET /openapi", openapiHandler)

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func bookingOffersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	flightURL, err := resolveServiceURL("flight-service")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resolve flight-service: %v", err), http.StatusInternalServerError)
		return
	}

	hotelURL, err := resolveServiceURL("hotel-service")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resolve hotel-service: %v", err), http.StatusInternalServerError)
		return
	}

	carURL, err := resolveServiceURL("car-service")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resolve car-service: %v", err), http.StatusInternalServerError)
		return
	}

	flights, err := fetchJSON(fmt.Sprintf("%s/flights", flightURL))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch flights: %v", err), http.StatusInternalServerError)
		return
	}

	hotels, err := fetchJSON(fmt.Sprintf("%s/hotels", hotelURL))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch hotels: %v", err), http.StatusInternalServerError)
		return
	}

	cars, err := fetchJSON(fmt.Sprintf("%s/cars", carURL))
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

func fetchJSON(url string) (json.RawMessage, error) {
	resp, err := httpClient.Get(url)
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

func openapiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(openapiSpec)
}
