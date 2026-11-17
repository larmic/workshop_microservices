package handler

import (
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/shared/tracing"
)

func OpenapiHandler(spec []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracing.Logger(r.Context()).Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		w.Write(spec)
	}
}
