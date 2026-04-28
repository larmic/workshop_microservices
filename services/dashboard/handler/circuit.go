package handler

import (
	"io"
	"net/http"
	"time"
)

func CircuitStateHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get(bookingURL + "/admin/circuit-state")
		if err != nil {
			http.Error(w, "booking-story3 unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}
