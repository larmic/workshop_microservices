package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/hotel/handler"
)

//go:embed api/openapi.yaml
var openapiSpec []byte

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
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("GET /hotels", handler.HotelsHandler)
	mux.HandleFunc("GET /openapi", handler.OpenapiHandler(openapiSpec))

	log.Println("HotelService starting on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
