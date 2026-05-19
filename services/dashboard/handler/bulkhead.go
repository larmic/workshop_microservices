package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func BulkheadStateHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get(bookingURL + "/admin/bulkhead-state")
		if err != nil {
			http.Error(w, "booking service unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}

func BulkheadResetHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, bookingURL+"/admin/bulkhead-reset", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "booking service unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
	}
}

func BookingOffersProxyHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get(bookingURL + "/booking/offers")
		if err != nil {
			http.Error(w, "booking service unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}

type BurstResult struct {
	Total      int   `json:"total"`
	Succeeded  int   `json:"succeeded"`
	Failed     int   `json:"failed"`
	DurationMs int64 `json:"durationMs"`
}

// BurstHandler feuert N parallele GET /booking/offers gegen den Booking-Service.
// Ohne paralleles Aufschütten ist der Bulkhead-Effekt im Workshop nicht
// sichtbar — ein einzelner Aufruf belegt nie mehr als einen Slot.
func BurstHandler(bookingURL string, parallel int) http.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		var succeeded, failed atomic.Int64
		start := time.Now()

		var wg sync.WaitGroup
		wg.Add(parallel)
		for i := 0; i < parallel; i++ {
			go func() {
				defer wg.Done()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, bookingURL+"/booking/offers", nil)
				if err != nil {
					failed.Add(1)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					failed.Add(1)
					return
				}
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				if resp.StatusCode >= 400 {
					failed.Add(1)
				} else {
					succeeded.Add(1)
				}
			}()
		}
		wg.Wait()

		result := BurstResult{
			Total:      parallel,
			Succeeded:  int(succeeded.Load()),
			Failed:     int(failed.Load()),
			DurationMs: time.Since(start).Milliseconds(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}
}
