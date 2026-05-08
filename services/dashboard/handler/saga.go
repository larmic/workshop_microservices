package handler

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

func SagaStateHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get(bookingURL + "/admin/sagas")
		if err != nil {
			http.Error(w, "booking-story5 unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}

func SagaResetHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 2 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, bookingURL+"/admin/sagas-reset", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "booking-story5 unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
	}
}

// SagaTriggerHandler stößt eine Saga mit festen Demo-Werten an. Bewusst
// kein freier Body, damit das Dashboard-Panel nicht zum POST-Editor wird —
// für individuelle Buchungen nutzt man Swagger.
func SagaTriggerHandler(bookingURL string) http.HandlerFunc {
	client := &http.Client{Timeout: 15 * time.Second}
	body := []byte(`{"flightId":"LH123","hotelId":"H1","carId":"C1","customerName":"Dashboard Demo"}`)
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, bookingURL+"/booking/bookings", bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "booking-story5 unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}
}
