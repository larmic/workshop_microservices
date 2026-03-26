package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/flight/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /flights", handler.FlightsHandler)
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))

	consulURL := os.Getenv("CONSUL_URL")
	serviceName := os.Getenv("SERVICE_NAME")
	serviceAddress := os.Getenv("SERVICE_ADDRESS")

	if consulURL != "" {
		cfg := consul.ServiceConfig{Name: serviceName, Address: serviceAddress, Port: 8080}
		if err := consul.Register(consulURL, cfg); err != nil {
			log.Printf("WARNING: Consul registration failed: %v", err)
		}
	}

	srv := &http.Server{Addr: ":8080", Handler: middleware.CORSMiddleware(mux)}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("FlightService starting on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down FlightService...")

	if consulURL != "" {
		consul.Deregister(consulURL, serviceName)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("WARNING: HTTP server shutdown error: %v", err)
	}
}
