package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

type Config struct {
	FlightServiceURL string `json:"flightServiceUrl"`
	HotelServiceURL  string `json:"hotelServiceUrl"`
	CarServiceURL    string `json:"carServiceUrl"`
	Timeout          int    `json:"timeout"`
	ActiveProfile    string `json:"activeProfile"`
}

type BookingOffers struct {
	Flights json.RawMessage `json:"flights"`
	Hotels  json.RawMessage `json:"hotels"`
	Cars    json.RawMessage `json:"cars"`
}

var config Config
var httpClient *http.Client

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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
		FlightServiceURL: getEnv("FLIGHT_SERVICE_URL", "http://localhost:8081"),
		HotelServiceURL:  getEnv("HOTEL_SERVICE_URL", "http://localhost:8082"),
		CarServiceURL:    getEnv("CAR_SERVICE_URL", "http://localhost:8083"),
		Timeout:          5000,
		ActiveProfile:    getEnv("ACTIVE_PROFILE", "dev"),
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func bookingOffersHandler(w http.ResponseWriter, r *http.Request) {
	flights, err := fetchJSON(fmt.Sprintf("%s/flights", config.FlightServiceURL))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch flights: %v", err), http.StatusInternalServerError)
		return
	}

	hotels, err := fetchJSON(fmt.Sprintf("%s/hotels", config.HotelServiceURL))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch hotels: %v", err), http.StatusInternalServerError)
		return
	}

	cars, err := fetchJSON(fmt.Sprintf("%s/cars", config.CarServiceURL))
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
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(openapiSpec)
}
