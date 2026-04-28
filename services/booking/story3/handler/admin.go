package handler

import (
	"encoding/json"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/booking/story3/circuitbreaker"
)

func CircuitStateHandler(cb *circuitbreaker.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cb.Snapshot())
	}
}

func CircuitEventsHandler(cb *circuitbreaker.CircuitBreaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cb.RecentEvents())
	}
}
