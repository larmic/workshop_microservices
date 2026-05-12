package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/team-neusta-skills/workshop_microservices/booking/story7/bulkhead"
	"github.com/team-neusta-skills/workshop_microservices/booking/story7/circuitbreaker"
	"github.com/team-neusta-skills/workshop_microservices/booking/story7/saga"
	"github.com/team-neusta-skills/workshop_microservices/shared/consul"
	"github.com/team-neusta-skills/workshop_microservices/shared/tracing"
)

type Config struct {
	Service   string `json:"service"`
	ConsulURL string `json:"consulUrl"`
	Timeout   int    `json:"timeout"`
}

type Breakers struct {
	Flight *circuitbreaker.CircuitBreaker
	Hotel  *circuitbreaker.CircuitBreaker
	Car    *circuitbreaker.CircuitBreaker
}

func (b Breakers) All() []*circuitbreaker.CircuitBreaker {
	return []*circuitbreaker.CircuitBreaker{b.Flight, b.Hotel, b.Car}
}

type Bulkheads struct {
	Flight *bulkhead.Bulkhead
	Hotel  *bulkhead.Bulkhead
	Car    *bulkhead.Bulkhead
}

func (b Bulkheads) All() []*bulkhead.Bulkhead {
	return []*bulkhead.Bulkhead{b.Flight, b.Hotel, b.Car}
}

type BookingOffers struct {
	Flights json.RawMessage `json:"flights"`
	Hotels  json.RawMessage `json:"hotels"`
	Cars    json.RawMessage `json:"cars"`
}

type BookingRequest struct {
	FlightID     string `json:"flightId"`
	HotelID      string `json:"hotelId"`
	CarID        string `json:"carId"`
	CustomerName string `json:"customerName"`
}

type Booking struct {
	BookingID    string          `json:"bookingId"`
	CustomerName string          `json:"customerName"`
	Flight       json.RawMessage `json:"flight"`
	Hotel        json.RawMessage `json:"hotel"`
	Car          json.RawMessage `json:"car"`
}

var emptyJSONArray = json.RawMessage("[]")

// consulNameFor liefert den Consul-Service-Namen pro Saga-Schritt.
var consulNameFor = map[string]string{
	"flight": "flight-service",
	"hotel":  "hotel-service",
	"car":    "car-service",
}

func BookingOffersHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers, bulkheads Bulkheads) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		flights := fetchOffersWithBHCB(r.Context(), w, bulkheads.Flight, breakers.Flight, "flight",
			resolver, client, "flight-service", "/flights")
		hotels := fetchOffersWithBHCB(r.Context(), w, bulkheads.Hotel, breakers.Hotel, "hotel",
			resolver, client, "hotel-service", "/hotels")
		cars := fetchOffersWithBHCB(r.Context(), w, bulkheads.Car, breakers.Car, "car",
			resolver, client, "car-service", "/cars")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BookingOffers{
			Flights: flights,
			Hotels:  hotels,
			Cars:    cars,
		})
	}
}

// CreateBookingHandler orchestriert eine Saga über Flight → Hotel → Car.
// Schlägt ein Forward-Step fehl, werden die bereits gebuchten Schritte in
// umgekehrter Reihenfolge kompensiert (DELETE /bookings/{id}). Der
// Saga-Status wird im Store gehalten und ist über
// GET /booking/bookings/{id} abrufbar.
func CreateBookingHandler(resolver *consul.Resolver, client *http.Client, breakers Breakers, store *saga.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := tracing.Logger(ctx)
		logRequest(r)

		// Pre-Check: ist ein CB OPEN, fängt die Saga gar nicht erst an.
		// Der Pre-Check spart einen Forward-Step plus die anschließende
		// Kompensation, die ohnehin scheitern würde.
		for _, cb := range breakers.All() {
			snap := cb.Snapshot()
			if snap.State == "OPEN" {
				writePreCheckFailure(ctx, w, snap.Name)
				return
			}
		}

		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		sagaID := newBookingID()
		now := time.Now().UTC()
		s := saga.Saga{
			SagaID:       sagaID,
			CustomerName: req.CustomerName,
			Status:       saga.StatusPending,
			CreatedAt:    now,
		}
		store.Save(s)
		logger.Info("saga started", "sagaId", sagaID, "customer", req.CustomerName)

		type stepDef struct {
			label   string
			cb      *circuitbreaker.CircuitBreaker
			payload map[string]string
		}
		defs := []stepDef{
			{"flight", breakers.Flight, map[string]string{"flightId": req.FlightID, "customerName": req.CustomerName}},
			{"hotel", breakers.Hotel, map[string]string{"hotelId": req.HotelID, "customerName": req.CustomerName}},
			{"car", breakers.Car, map[string]string{"carId": req.CarID, "customerName": req.CustomerName}},
		}

		for _, def := range defs {
			raw, err := bookSingle(ctx, def.cb, resolver, client, consulNameFor[def.label], def.payload)
			if err != nil {
				s.Status = saga.StatusCompensating
				s.FailedAt = def.label
				s.Reason = err.Error()
				store.Save(s)
				logger.Warn("saga forward step failed — compensating",
					"sagaId", sagaID,
					"step", def.label,
					"err", err.Error(),
					"compensateSteps", len(s.Steps),
				)

				// Kompensation läuft in eigenem Background-Context: das
				// Aufräumen darf nicht abgebrochen werden, nur weil der
				// Client die Verbindung zumacht. Der Trace-Kontext wird
				// explizit kopiert, damit alle Compensation-Logzeilen
				// dieselbe trace_id tragen.
				compCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				if tc, ok := tracing.FromContext(ctx); ok {
					compCtx = tracing.WithContext(compCtx, tc)
				}
				compensate(compCtx, client, resolver, &s)
				cancel()

				s.Status = saga.StatusFailed
				store.Save(s)

				writeSagaFailure(w, s, fmt.Sprintf("saga failed at %s — see saga state for details", def.label))
				return
			}
			s.Steps = append(s.Steps, saga.Step{
				Service:   def.label,
				BookingID: extractBookingID(raw),
				Status:    saga.StepBooked,
				Detail:    raw,
			})
			store.Save(s)
		}

		s.Status = saga.StatusCompleted
		store.Save(s)
		logger.Info("saga completed", "sagaId", sagaID)

		booking := Booking{
			BookingID:    sagaID,
			CustomerName: req.CustomerName,
			Flight:       s.Steps[0].Detail,
			Hotel:        s.Steps[1].Detail,
			Car:          s.Steps[2].Detail,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(booking)
	}
}

func GetSagaStatusHandler(store *saga.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		id := r.PathValue("id")
		s, ok := store.Get(id)
		if !ok {
			http.Error(w, "saga not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s)
	}
}

// ListSagasHandler liefert alle Sagas, neueste zuerst. Für das
// Dashboard-Panel.
func ListSagasHandler(store *saga.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.List())
	}
}

// ResetSagasHandler leert den Saga-Store. Demo-Helfer.
func ResetSagasHandler(store *saga.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Reset()
		tracing.Logger(r.Context()).Info("saga store reset")
		w.WriteHeader(http.StatusNoContent)
	}
}

// compensate löst die Kompensation für jeden bereits gebuchten Schritt
// asynchron aus: Pro Schritt schickt Booking ein CompensationRequested-Event
// per HTTP-POST an Flight/Hotel/Car. Der Backend-Service antwortet sofort
// mit 202 Accepted und erledigt die Stornierung in einer Goroutine — Booking
// wartet NICHT auf die fachliche Verarbeitung.
//
// Story 7-Aspekt: Der Trace-Kontext wandert über die `traceparent`-Property
// im Event-Body mit, sodass die spätere Stornierungs-Logzeile im Backend
// dieselbe trace_id trägt wie die Buchung — auch über die Async-Grenze.
func compensate(ctx context.Context, client *http.Client, resolver *consul.Resolver, s *saga.Saga) {
	logger := tracing.Logger(ctx)
	dispatched := 0
	for i := len(s.Steps) - 1; i >= 0; i-- {
		step := &s.Steps[i]
		if step.Status != saga.StepBooked {
			continue
		}
		if step.BookingID == "" {
			step.Status = saga.StepCompensationFailed
			step.Reason = "no bookingId captured — cannot compensate"
			logger.Warn("compensation skipped: no bookingId",
				"sagaId", s.SagaID, "service", step.Service)
			continue
		}
		eventID := newEventID()
		stepLogger := logger.With(
			"sagaId", s.SagaID,
			"event", "CompensationRequested",
			"eventId", eventID,
			"service", step.Service,
			"bookingId", step.BookingID,
		)
		stepLogger.Info("compensation event", "phase", "publishing")
		if err := compensateSingle(ctx, client, resolver, consulNameFor[step.Service], eventID, s.SagaID, step.BookingID); err != nil {
			stepLogger.Error("compensation event", "phase", "dispatch-failed", "err", err.Error())
		} else {
			stepLogger.Info("compensation event", "phase", "dispatched")
		}
		step.Status = saga.StepCompensated
		dispatched++
	}
	logger.Info("compensation dispatched — booking finishes immediately",
		"sagaId", s.SagaID, "services", dispatched)
}

func compensateSingle(ctx context.Context, client *http.Client, resolver *consul.Resolver, consulName, eventID, sagaID, bookingID string) error {
	url, err := resolver.ResolveServiceURL(consulName)
	if err != nil {
		return err
	}
	traceparent := ""
	if tc, ok := tracing.FromContext(ctx); ok {
		traceparent = tc.Header()
	}
	payload, err := json.Marshal(map[string]string{
		"eventId":     eventID,
		"sagaId":      sagaID,
		"bookingId":   bookingID,
		"traceparent": traceparent,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/events/compensation", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	tracing.Inject(ctx, req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend returned %d (expected 202 Accepted): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func newEventID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "ev-" + hex.EncodeToString(b)
}

func extractBookingID(raw json.RawMessage) string {
	var probe struct {
		BookingID string `json:"bookingId"`
	}
	_ = json.Unmarshal(raw, &probe)
	return probe.BookingID
}

type SagaFailureResponse struct {
	Error string    `json:"error"`
	Saga  saga.Saga `json:"saga"`
}

func writeSagaFailure(w http.ResponseWriter, s saga.Saga, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_ = json.NewEncoder(w).Encode(SagaFailureResponse{Error: message, Saga: s})
}

func writePreCheckFailure(ctx context.Context, w http.ResponseWriter, service string) {
	now := time.Now().UTC()
	s := saga.Saga{
		Status:    saga.StatusFailed,
		FailedAt:  service,
		Reason:    "circuit breaker is open — saga not started",
		CreatedAt: now,
		UpdatedAt: now,
	}
	tracing.Logger(ctx).Warn("saga not started: circuit breaker open", "service", service)
	writeSagaFailure(w, s, "saga not started — circuit breaker is open for "+service)
}

func bookSingle(
	ctx context.Context,
	cb *circuitbreaker.CircuitBreaker,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	payload map[string]string,
) (json.RawMessage, error) {
	var data json.RawMessage
	err := cb.Execute(ctx, func(ctx context.Context) error {
		url, resolveErr := resolver.ResolveServiceURL(consulName)
		if resolveErr != nil {
			return resolveErr
		}
		raw, postErr := postJSON(ctx, client, url+"/bookings", payload)
		if postErr != nil {
			return postErr
		}
		data = raw
		return nil
	})
	return data, err
}

// fetchOffersWithBHCB legt den Bulkhead VOR den Circuit Breaker: ist der Pool
// voll, kürzen wir sofort ab — ohne die CB-Statistik zu vergiften, weil der
// Backend-Call ja gar nicht stattgefunden hat. So bleiben CB und Bulkhead in
// ihren Zuständigkeiten getrennt: CB misst Backend-Gesundheit, Bulkhead schützt
// vor Ressourcen-Erschöpfung.
func fetchOffersWithBHCB(
	ctx context.Context,
	w http.ResponseWriter,
	bh *bulkhead.Bulkhead,
	cb *circuitbreaker.CircuitBreaker,
	serviceLabel string,
	resolver *consul.Resolver,
	client *http.Client,
	consulName string,
	path string,
) json.RawMessage {
	var data json.RawMessage
	err := bh.Execute(ctx, func(ctx context.Context) error {
		return cb.Execute(ctx, func(ctx context.Context) error {
			url, resolveErr := resolver.ResolveServiceURL(consulName)
			if resolveErr != nil {
				return resolveErr
			}
			raw, fetchErr := fetchJSON(ctx, client, url+path)
			if fetchErr != nil {
				return fetchErr
			}
			data = raw
			return nil
		})
	})
	if err != nil {
		markFallback(ctx, w, serviceLabel, err)
		return emptyJSONArray
	}
	return data
}

func markFallback(ctx context.Context, w http.ResponseWriter, serviceLabel string, err error) {
	logger := tracing.Logger(ctx)
	w.Header().Add("X-Fallback", serviceLabel)
	switch {
	case errors.Is(err, bulkhead.ErrBulkheadFull):
		w.Header().Add("X-Bulkhead-Full", serviceLabel)
		logger.Warn("call rejected: bulkhead full", "service", serviceLabel)
	case errors.Is(err, circuitbreaker.ErrCircuitOpen):
		w.Header().Add("X-Circuit-Open", serviceLabel)
		logger.Warn("call short-circuited: CB OPEN", "service", serviceLabel)
	default:
		logger.Warn("call failed — applying fallback", "service", serviceLabel, "err", err.Error())
	}
}

func fetchJSON(ctx context.Context, client *http.Client, url string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	tracing.Inject(ctx, req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("backend returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.RawMessage(body), nil
}

func postJSON(ctx context.Context, client *http.Client, url string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	tracing.Inject(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("backend returned %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	return json.RawMessage(respBody), nil
}

func logRequest(r *http.Request) {
	tracing.Logger(r.Context()).Info("request",
		"method", r.Method,
		"path", r.URL.Path,
		"remote", r.RemoteAddr,
	)
}

func newBookingID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "B-" + hex.EncodeToString(b)
}
