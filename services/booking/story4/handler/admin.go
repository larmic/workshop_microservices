package handler

import (
	"encoding/json"
	"net/http"

	"github.com/team-neusta-skills/workshop_microservices/booking/story4/circuitbreaker"
)

func CircuitStateHandler(breakers Breakers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := breakers.All()
		snapshots := make([]circuitbreaker.Snapshot, 0, len(all))
		for _, cb := range all {
			snapshots = append(snapshots, cb.Snapshot())
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snapshots)
	}
}

func CircuitEventsHandler(breakers Breakers) http.HandlerFunc {
	type namedEvents struct {
		Name   string                 `json:"name"`
		Events []circuitbreaker.Event `json:"events"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		all := breakers.All()
		out := make([]namedEvents, 0, len(all))
		for _, cb := range all {
			snap := cb.Snapshot()
			out = append(out, namedEvents{
				Name:   snap.Name,
				Events: cb.RecentEvents(),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}
