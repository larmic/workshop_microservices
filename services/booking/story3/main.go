package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story3/circuitbreaker"
	"github.com/team-neusta-skills/workshop_microservices/booking/story3/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	"github.com/team-neusta-skills/workshop_microservices/shared/env"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

func main() {
	config := handler.Config{
		ConsulURL:     env.GetEnv("CONSUL_URL", "http://localhost:8500"),
		Timeout:       3000,
		ActiveProfile: env.GetEnv("ACTIVE_PROFILE", "dev"),
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	resolver := consul.NewResolver(config.ConsulURL, httpClient)

	cb := circuitbreaker.New(circuitbreaker.Config{
		Name:             "flight-service",
		FailureThreshold: 5,
		OpenTimeout:      30 * time.Second,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /info", sharedhandler.InfoHandler(config))
	mux.HandleFunc("GET /booking/offers", handler.BookingOffersHandler(resolver, httpClient, cb))
	mux.HandleFunc("POST /booking/bookings", handler.CreateBookingHandler(resolver, httpClient, cb))
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))
	mux.HandleFunc("GET /admin/circuit-state", handler.CircuitStateHandler(cb))
	mux.HandleFunc("GET /admin/circuit-events", handler.CircuitEventsHandler(cb))

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", middleware.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
