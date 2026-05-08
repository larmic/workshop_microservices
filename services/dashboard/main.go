package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/dashboard/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed static
var staticFS embed.FS

func main() {
	projectName := getEnv("COMPOSE_PROJECT_NAME", "services")
	composeFilesRaw := getEnv("COMPOSE_FILES", "/compose/docker-compose.yml,/compose/docker-compose.infra.yml,/compose/docker-compose.reference.yml")
	consulURL := getEnv("CONSUL_URL", "http://consul:8500")
	bookingStory3URL := getEnv("BOOKING_STORY3_URL", "http://booking-story3:8080")
	bookingStory4URL := getEnv("BOOKING_STORY4_URL", "http://booking-story4:8080")

	var composeArgs []string
	for _, f := range strings.Split(composeFilesRaw, ",") {
		composeArgs = append(composeArgs, "-f", strings.TrimSpace(f))
	}
	composeArgs = append(composeArgs, "-p", projectName)

	allowedServices := []string{"flight", "hotel", "car"}

	staticContent, _ := fs.Sub(staticFS, "static")

	resolver := consul.NewResolver(consulURL, &http.Client{Timeout: 2 * time.Second})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /api/services", handler.ListServicesHandler(composeArgs, allowedServices))
	mux.HandleFunc("POST /api/services/{name}/scale", handler.ScaleServiceHandler(composeArgs, allowedServices))
	mux.HandleFunc("GET /api/services/{name}/instances", handler.ListInstancesHandler(resolver, allowedServices))
	mux.HandleFunc("POST /api/services/{name}/chaos", handler.SetChaosHandler(resolver, allowedServices))
	mux.HandleFunc("GET /api/circuit-state", handler.CircuitStateHandler(bookingStory3URL))
	mux.HandleFunc("GET /api/bulkhead-state", handler.BulkheadStateHandler(bookingStory4URL))
	mux.HandleFunc("POST /api/bulkhead-reset", handler.BulkheadResetHandler(bookingStory4URL))
	mux.HandleFunc("GET /api/booking-story4/offers", handler.BookingOffersProxyHandler(bookingStory4URL))
	mux.HandleFunc("POST /api/booking-story4/burst", handler.BurstHandler(bookingStory4URL, 20))
	mux.Handle("GET /", http.FileServer(http.FS(staticContent)))

	srv := &http.Server{Addr: ":8080", Handler: middleware.CORSMiddleware(mux)}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Dashboard starting on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Dashboard...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("WARNING: HTTP server shutdown error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
