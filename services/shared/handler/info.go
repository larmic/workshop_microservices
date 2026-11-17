package handler

import (
	"encoding/json"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/shared/tracing"
)

func InfoHandler(config any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracing.Logger(r.Context()).Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}
}
