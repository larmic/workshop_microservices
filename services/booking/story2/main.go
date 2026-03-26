package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story2/handler"
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
		Timeout:       5000,
		ActiveProfile: env.GetEnv("ACTIVE_PROFILE", "dev"),
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	resolver := consul.NewResolver(config.ConsulURL, httpClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /info", handler.InfoHandler(config))
	mux.HandleFunc("GET /booking/offers", handler.BookingOffersHandler(resolver, httpClient))
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", middleware.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
