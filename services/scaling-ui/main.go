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

	"github.com/team-neusta-skills/workshop_microservices/scaling-ui/handler"
	sharedhandler "github.com/team-neusta-skills/workshop_microservices/shared/handler"
	"github.com/team-neusta-skills/workshop_microservices/shared/middleware"
)

//go:embed static
var staticFS embed.FS

func main() {
	projectName := getEnv("COMPOSE_PROJECT_NAME", "services")
	composeFilesRaw := getEnv("COMPOSE_FILES", "/compose/docker-compose.yml,/compose/docker-compose.infra.yml,/compose/docker-compose.reference.yml")

	var composeArgs []string
	for _, f := range strings.Split(composeFilesRaw, ",") {
		composeArgs = append(composeArgs, "-f", strings.TrimSpace(f))
	}
	composeArgs = append(composeArgs, "-p", projectName)

	allowedServices := []string{"flight", "hotel", "car"}

	staticContent, _ := fs.Sub(staticFS, "static")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", sharedhandler.HealthHandler)
	mux.HandleFunc("GET /api/services", handler.ListServicesHandler(composeArgs, allowedServices))
	mux.HandleFunc("POST /api/services/{name}/scale", handler.ScaleServiceHandler(composeArgs, allowedServices))
	mux.Handle("GET /", http.FileServer(http.FS(staticContent)))

	srv := &http.Server{Addr: ":8080", Handler: middleware.CORSMiddleware(mux)}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Scaling UI starting on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Scaling UI...")

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
