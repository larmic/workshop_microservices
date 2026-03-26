package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story2/consul"
	"github.com/team-neusta-skills/workshop_microservices/booking/story2/handler"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

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
	config := handler.Config{
		ConsulURL:     getEnv("CONSUL_URL", "http://localhost:8500"),
		Timeout:       5000,
		ActiveProfile: getEnv("ACTIVE_PROFILE", "dev"),
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}

	resolver := consul.NewResolver(config.ConsulURL, httpClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("GET /info", handler.InfoHandler(config))
	mux.HandleFunc("GET /booking/offers", handler.BookingOffersHandler(resolver, httpClient))
	mux.HandleFunc("GET /openapi", handler.OpenapiHandler(openapiSpec))

	log.Println("BookingService starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
