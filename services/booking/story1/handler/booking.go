package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	Service          string `json:"service"`
	FlightServiceURL string `json:"flightServiceUrl"`
	HotelServiceURL  string `json:"hotelServiceUrl"`
	CarServiceURL    string `json:"carServiceUrl"`
	Timeout          int    `json:"timeout"`
}

type BookingOffers struct {
	Flights json.RawMessage `json:"flights"`
	Hotels  json.RawMessage `json:"hotels"`
	Cars    json.RawMessage `json:"cars"`
}

func BookingOffersHandler(config Config, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		flights, err := fetchJSON(client, fmt.Sprintf("%s/flights", config.FlightServiceURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch flights: %v", err), http.StatusInternalServerError)
			return
		}

		hotels, err := fetchJSON(client, fmt.Sprintf("%s/hotels", config.HotelServiceURL))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch hotels: %v", err), http.StatusInternalServerError)
			return
		}

		cars, err := fetchJSON(client, fmt.Sprintf("%s/cars", config.CarServiceURL))
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
