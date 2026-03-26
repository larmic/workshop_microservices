package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story1/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/env"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

func main() {
	config := handler.Config{
		FlightServiceURL: env.GetEnv("FLIGHT_SERVICE_URL", "http://localhost:8081"),
		HotelServiceURL:  env.GetEnv("HOTEL_SERVICE_URL", "http://localhost:8082"),
		CarServiceURL:    env.GetEnv("CAR_SERVICE_URL", "http://localhost:8083"),
		Timeout:          5000,
		ActiveProfile:    env.GetEnv("ACTIVE_PROFILE", "dev"),
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /info", handler.InfoHandler(config))
	mux.HandleFunc("GET /booking/offers", handler.BookingOffersHandler(config, httpClient))
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", middleware.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
