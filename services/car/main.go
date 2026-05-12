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

	"github.com/team-neusta-skills/workshop_microservices/car/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/chaos"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

func main() {
	consulURL := os.Getenv("CONSUL_URL")
	serviceName := os.Getenv("SERVICE_NAME")
	serviceAddress := os.Getenv("SERVICE_ADDRESS")

	info := map[string]string{
		"serviceName":    serviceName,
		"serviceAddress": serviceAddress,
	}

	chaosState := chaos.New()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /info", sharedhandler.InfoHandler(info))
	mux.HandleFunc("GET /cars", handler.CarsHandler)
	mux.HandleFunc("POST /bookings", handler.CreateBookingHandler)
	mux.HandleFunc("DELETE /bookings/{id}", handler.CancelBookingHandler)
	mux.HandleFunc("POST /events/compensation", handler.CompensationEventHandler)
	mux.HandleFunc("GET /openapi", sharedhandler.OpenapiHandler(openapiSpec))
	mux.HandleFunc("GET /admin/chaos", chaosState.GetHandler)
	mux.HandleFunc("POST /admin/chaos", chaosState.SetHandler)

	var serviceID string
	if consulURL != "" {
		cfg := consul.ServiceConfig{Name: serviceName, Address: serviceAddress, Port: 8080}
		var err error
		serviceID, err = consul.Register(consulURL, cfg)
		if err != nil {
			log.Printf("WARNING: Consul registration failed: %v", err)
		}
	}

	srv := &http.Server{Addr: ":8080", Handler: middleware.CORSMiddleware(chaosState.Middleware(mux))}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("CarService starting on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down CarService...")

	if serviceID != "" {
		consul.Deregister(consulURL, serviceID)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("WARNING: HTTP server shutdown error: %v", err)
	}
}
