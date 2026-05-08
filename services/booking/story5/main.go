package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story5/bulkhead"
	"github.com/team-neusta-skills/workshop_microservices/booking/story5/circuitbreaker"
	"github.com/team-neusta-skills/workshop_microservices/booking/story5/handler"
	"github.com/team-neusta-skills/workshop_microservices/booking/story5/saga"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	"github.com/team-neusta-skills/workshop_microservices/shared/env"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

func main() {
	config := handler.Config{
		Service:   "booking-5",
		ConsulURL: env.GetEnv("CONSUL_URL", "http://localhost:8500"),
		Timeout:   3000,
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	resolver := consul.NewResolver(config.ConsulURL, httpClient)

	cbConfig := func(name string) circuitbreaker.Config {
		return circuitbreaker.Config{
			Name:             name,
			FailureThreshold: 5,
			OpenTimeout:      30 * time.Second,
		}
	}
	breakers := handler.Breakers{
		Flight: circuitbreaker.New(cbConfig("flight")),
		Hotel:  circuitbreaker.New(cbConfig("hotel")),
		Car:    circuitbreaker.New(cbConfig("car")),
	}

	bhConfig := func(name string) bulkhead.Config {
		return bulkhead.Config{
			Name:          name,
			MaxConcurrent: 10,
		}
	}
	bulkheads := handler.Bulkheads{
		Flight: bulkhead.New(bhConfig("flight")),
		Hotel:  bulkhead.New(bhConfig("hotel")),
		Car:    bulkhead.New(bhConfig("car")),
	}

	sagaStore := saga.NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /info", sharedhandler.InfoHandler(config))
	mux.HandleFunc("GET /booking/offers", handler.BookingOffersHandler(resolver, httpClient, breakers, bulkheads))
	mux.HandleFunc("POST /booking/bookings", handler.CreateBookingHandler(resolver, httpClient, breakers, sagaStore))
	mux.HandleFunc("GET /booking/bookings/{id}", handler.GetSagaStatusHandler(sagaStore))
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))
	mux.HandleFunc("GET /admin/circuit-state", handler.CircuitStateHandler(breakers))
	mux.HandleFunc("GET /admin/circuit-events", handler.CircuitEventsHandler(breakers))
	mux.HandleFunc("GET /admin/bulkhead-state", handler.BulkheadStateHandler(bulkheads))
	mux.HandleFunc("POST /admin/bulkhead-reset", handler.BulkheadResetHandler(bulkheads))

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", middleware.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
