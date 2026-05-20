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
	bookingRefStory1URL := getEnv("BOOKING_REF_STORY1_URL", "http://booking-ref-story1:8080")
	bookingRefStory2URL := getEnv("BOOKING_REF_STORY2_URL", "http://booking-ref-story2:8080")
	bookingRefStory3URL := getEnv("BOOKING_REF_STORY3_URL", "http://booking-ref-story3:8080")
	bookingRefStory4URL := getEnv("BOOKING_REF_STORY4_URL", "http://booking-ref-story4:8080")
	bookingRefStory5URL := getEnv("BOOKING_REF_STORY5_URL", "http://booking-ref-story5:8080")
	bookingRefStory6URL := getEnv("BOOKING_REF_STORY6_URL", "http://booking-ref-story6:8080")
	bookingRefStory7URL := getEnv("BOOKING_REF_STORY7_URL", "http://booking-ref-story7:8080")
	// Leer-String = Custom-Service ist nicht Teil des aktuellen Setups.
	// Nur die docker-compose.custom.yml setzt die ENV; bei einem reinen
	// Reference-Setup bleibt sie unset und der Custom-Teil wird komplett
	// ausgeblendet.
	bookingCustomURL := getEnv("BOOKING_CUSTOM_URL", "")
	traefikPingURL := getEnv("TRAEFIK_PING_URL", "http://traefik:8080/api/version")
	swaggerUIURL := getEnv("SWAGGER_UI_URL", "http://swagger-ui:8080/api/")
	consulStatusURL := getEnv("CONSUL_STATUS_URL", consulURL+"/v1/status/leader")
	feedbackWebhookURL := getEnv("FEEDBACK_WEBHOOK_URL", "")

	var composeArgs []string
	for _, f := range strings.Split(composeFilesRaw, ",") {
		composeArgs = append(composeArgs, "-f", strings.TrimSpace(f))
	}
	composeArgs = append(composeArgs, "-p", projectName)

	allowedServices := []string{"flight", "hotel", "car"}

	bookingURLs := map[string]string{
		"booking-ref-story1": bookingRefStory1URL,
		"booking-ref-story2": bookingRefStory2URL,
		"booking-ref-story3": bookingRefStory3URL,
		"booking-ref-story4": bookingRefStory4URL,
		"booking-ref-story5": bookingRefStory5URL,
		"booking-ref-story6": bookingRefStory6URL,
		"booking-ref-story7": bookingRefStory7URL,
	}
	if bookingCustomURL != "" {
		bookingURLs["booking-custom"] = bookingCustomURL
	}

	infraTargets := []handler.InfraTarget{
		{Name: "consul", URL: consulStatusURL},
		{Name: "traefik", URL: traefikPingURL},
		{Name: "swagger-ui", URL: swaggerUIURL},
	}

	staticContent, _ := fs.Sub(staticFS, "static")

	resolver := consul.NewResolver(consulURL, &http.Client{Timeout: 2 * time.Second})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /api/services", handler.ListServicesHandler(composeArgs, allowedServices))
	mux.HandleFunc("GET /api/health-overview", handler.HealthOverviewHandler(composeArgs, allowedServices, resolver, bookingURLs, infraTargets))
	mux.HandleFunc("POST /api/services/{name}/scale", handler.ScaleServiceHandler(composeArgs, allowedServices))
	mux.HandleFunc("GET /api/services/{name}/instances", handler.ListInstancesHandler(resolver, allowedServices, composeArgs))
	mux.HandleFunc("POST /api/services/{name}/chaos", handler.SetChaosHandler(resolver, allowedServices))
	mux.HandleFunc("GET /api/booking-ref-story1/offers", handler.ProxyHandler(bookingRefStory1URL, http.MethodGet, "/booking/offers"))
	mux.HandleFunc("GET /api/booking-ref-story2/offers", handler.ProxyHandler(bookingRefStory2URL, http.MethodGet, "/booking/offers"))
	mux.HandleFunc("POST /api/booking-ref-story2/bookings", handler.ProxyHandler(bookingRefStory2URL, http.MethodPost, "/booking/bookings"))
	mux.HandleFunc("GET /api/booking-ref-story3/offers", handler.ProxyHandler(bookingRefStory3URL, http.MethodGet, "/booking/offers"))
	mux.HandleFunc("GET /api/circuit-state", handler.CircuitStateHandler(bookingRefStory3URL))
	mux.HandleFunc("GET /api/bulkhead-state", handler.BulkheadStateHandler(bookingRefStory4URL))
	mux.HandleFunc("POST /api/bulkhead-reset", handler.BulkheadResetHandler(bookingRefStory4URL))
	mux.HandleFunc("GET /api/booking-ref-story4/offers", handler.BookingOffersProxyHandler(bookingRefStory4URL))
	mux.HandleFunc("POST /api/booking-ref-story4/burst", handler.BurstHandler(bookingRefStory4URL, 20))
	mux.HandleFunc("GET /api/saga-state", handler.SagaStateHandler(bookingRefStory5URL))
	mux.HandleFunc("POST /api/saga-reset", handler.SagaResetHandler(bookingRefStory5URL))
	mux.HandleFunc("POST /api/saga-trigger", handler.SagaTriggerHandler(bookingRefStory5URL))
	mux.HandleFunc("GET /api/saga6-state", handler.SagaStateHandler(bookingRefStory6URL))
	mux.HandleFunc("POST /api/saga6-reset", handler.SagaResetHandler(bookingRefStory6URL))
	mux.HandleFunc("POST /api/saga6-trigger", handler.SagaTriggerHandler(bookingRefStory6URL))
	mux.HandleFunc("GET /api/booking-custom-available", func(w http.ResponseWriter, r *http.Request) {
		if bookingCustomURL == "" {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		client := &http.Client{Timeout: 1500 * time.Millisecond}
		resp, err := client.Get(bookingCustomURL + "/health")
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_ = resp.Body.Close()
		w.WriteHeader(http.StatusOK)
	})
	if bookingCustomURL != "" {
		mux.HandleFunc("GET /api/booking-custom/offers", handler.ProxyHandler(bookingCustomURL, http.MethodGet, "/booking/offers"))
		mux.HandleFunc("POST /api/booking-custom/bookings", handler.ProxyHandler(bookingCustomURL, http.MethodPost, "/booking/bookings"))
		mux.HandleFunc("GET /api/booking-custom/circuit-state", handler.CircuitStateHandler(bookingCustomURL))
		mux.HandleFunc("GET /api/booking-custom/bulkhead-state", handler.BulkheadStateHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/bulkhead-reset", handler.BulkheadResetHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/burst", handler.BurstHandler(bookingCustomURL, 20))
		mux.HandleFunc("GET /api/booking-custom/saga-state", handler.SagaStateHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/saga-reset", handler.SagaResetHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/saga-trigger", handler.SagaTriggerHandler(bookingCustomURL))
		mux.HandleFunc("GET /api/booking-custom/saga6-state", handler.SagaStateHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/saga6-reset", handler.SagaResetHandler(bookingCustomURL))
		mux.HandleFunc("POST /api/booking-custom/saga6-trigger", handler.SagaTriggerHandler(bookingCustomURL))
	}
	mux.HandleFunc("POST /api/feedback", handler.FeedbackHandler(feedbackWebhookURL, nil))
	mux.HandleFunc("GET /api/feedback/status", handler.FeedbackStatusHandler(feedbackWebhookURL))
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
