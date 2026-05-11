package handler

import (
	"io"
	"net/http"
	"time"
)

// ProxyHandler leitet einen Request unverändert (Method + Pfad) an den
// angegebenen Booking-Service weiter. Wird genutzt, damit das Dashboard
// nicht für jeden Story-Endpoint einen eigenen Handler braucht.
func ProxyHandler(targetBaseURL, method, path string) http.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), method, targetBaseURL+path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if ct := r.Header.Get("Content-Type"); ct != "" {
			req.Header.Set("Content-Type", ct)
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "upstream unreachable: "+err.Error(), http.StatusBadGateway)
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
