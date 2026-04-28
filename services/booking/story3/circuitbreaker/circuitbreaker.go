package circuitbreaker

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

func (s State) String() string {
	switch s {
	case Closed:
		return "CLOSED"
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

type Config struct {
	Name             string
	FailureThreshold int
	OpenTimeout      time.Duration
}

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Reason    string    `json:"reason"`
}

type Snapshot struct {
	Name               string     `json:"name"`
	State              string     `json:"state"`
	FailureCount       int        `json:"failureCount"`
	FailureThreshold   int        `json:"failureThreshold"`
	LastStateChangeAt  time.Time  `json:"lastStateChangeAt"`
	OpenUntil          *time.Time `json:"openUntil,omitempty"`
	OpenTimeoutSeconds int        `json:"openTimeoutSeconds"`
	TotalCalls         uint64     `json:"totalCalls"`
	TotalFailures      uint64     `json:"totalFailures"`
	TotalShortCircuits uint64     `json:"totalShortCircuits"`
}

const eventCapacity = 20

type CircuitBreaker struct {
	cfg               Config
	mu                sync.Mutex
	state             State
	failures          int
	openUntil         time.Time
	lastStateChangeAt time.Time
	probeInFlight     atomic.Bool

	totalCalls         atomic.Uint64
	totalFailures      atomic.Uint64
	totalShortCircuits atomic.Uint64

	events []Event
}

func New(cfg Config) *CircuitBreaker {
	if cfg.Name == "" {
		cfg.Name = "default"
	}
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.OpenTimeout <= 0 {
		cfg.OpenTimeout = 30 * time.Second
	}
	return &CircuitBreaker{cfg: cfg, lastStateChangeAt: time.Now()}
}

// Execute prüft den aktuellen State, lässt den Call passieren oder kürzt ab,
// und protokolliert jede Entscheidung im Log — damit man im Workshop nachvollziehen
// kann, warum der Breaker aktuell so handelt wie er handelt.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if !cb.allowRequest() {
		cb.totalShortCircuits.Add(1)
		return ErrCircuitOpen
	}
	cb.totalCalls.Add(1)
	err := fn(ctx)
	cb.recordResult(err)
	return err
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.maybePromoteFromOpenLocked()

	if cb.state == Open {
		remaining := time.Until(cb.openUntil).Round(time.Second)
		log.Printf("CB[%s] state=OPEN — Call SHORT-CIRCUIT (noch %s bis HALF_OPEN)", cb.cfg.Name, remaining)
		return false
	}

	if cb.state == HalfOpen {
		if cb.probeInFlight.CompareAndSwap(false, true) {
			log.Printf("CB[%s] state=HALF_OPEN — Probe-Slot belegt, Call wird durchgelassen (Test ob Backend wieder gesund)", cb.cfg.Name)
			return true
		}
		log.Printf("CB[%s] state=HALF_OPEN — Probe-Call läuft bereits, weiterer Call SHORT-CIRCUIT", cb.cfg.Name)
		return false
	}

	// CLOSED — keine Logzeile, sonst zu laut
	return true
}

// maybePromoteFromOpenLocked übernimmt die lazy Transition OPEN → HALF_OPEN, wenn
// das Timeout abgelaufen ist. Wird sowohl von allowRequest als auch von Snapshot
// aufgerufen, damit das Dashboard den State-Wechsel sieht, auch wenn keine Calls
// reinkommen.
func (cb *CircuitBreaker) maybePromoteFromOpenLocked() {
	if cb.state == Open && time.Now().After(cb.openUntil) {
		log.Printf("CB[%s] OPEN-Timeout abgelaufen → Übergang zu HALF_OPEN", cb.cfg.Name)
		cb.transitionLocked(HalfOpen, "open timeout elapsed")
	}
}

func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.totalFailures.Add(1)
		cb.failures++

		switch cb.state {
		case HalfOpen:
			cb.probeInFlight.Store(false)
			cb.openUntil = time.Now().Add(cb.cfg.OpenTimeout)
			log.Printf("CB[%s] Probe-Call FEHLGESCHLAGEN (err=%v) → zurück nach OPEN für %s",
				cb.cfg.Name, err, cb.cfg.OpenTimeout)
			cb.transitionLocked(Open, "probe failed in half-open")
		case Closed:
			if cb.failures >= cb.cfg.FailureThreshold {
				cb.openUntil = time.Now().Add(cb.cfg.OpenTimeout)
				log.Printf("CB[%s] Fehler %d/%d (err=%v) — Schwellenwert erreicht → OPEN für %s",
					cb.cfg.Name, cb.failures, cb.cfg.FailureThreshold, err, cb.cfg.OpenTimeout)
				cb.transitionLocked(Open, "failure threshold reached")
			} else {
				log.Printf("CB[%s] Fehler %d/%d gezählt (state=CLOSED, err=%v)",
					cb.cfg.Name, cb.failures, cb.cfg.FailureThreshold, err)
			}
		}
		return
	}

	switch cb.state {
	case HalfOpen:
		cb.probeInFlight.Store(false)
		log.Printf("CB[%s] Probe-Call ERFOLGREICH → CLOSED (Backend gilt wieder als gesund, Fehlerzähler zurückgesetzt)",
			cb.cfg.Name)
		cb.failures = 0
		cb.transitionLocked(Closed, "probe succeeded")
	case Closed:
		if cb.failures > 0 {
			log.Printf("CB[%s] Erfolgreicher Call (state=CLOSED) — Fehlerzähler von %d auf 0 zurückgesetzt",
				cb.cfg.Name, cb.failures)
		}
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) transitionLocked(to State, reason string) {
	if cb.state == to {
		return
	}
	from := cb.state
	cb.state = to
	cb.lastStateChangeAt = time.Now()

	log.Printf("CB[%s] STATE-CHANGE %s → %s (reason=%s)", cb.cfg.Name, from, to, reason)

	cb.events = append(cb.events, Event{
		Timestamp: cb.lastStateChangeAt,
		From:      from.String(),
		To:        to.String(),
		Reason:    reason,
	})
	if len(cb.events) > eventCapacity {
		cb.events = cb.events[len(cb.events)-eventCapacity:]
	}
}

func (cb *CircuitBreaker) Snapshot() Snapshot {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.maybePromoteFromOpenLocked()

	s := Snapshot{
		Name:               cb.cfg.Name,
		State:              cb.state.String(),
		FailureCount:       cb.failures,
		FailureThreshold:   cb.cfg.FailureThreshold,
		LastStateChangeAt:  cb.lastStateChangeAt,
		OpenTimeoutSeconds: int(cb.cfg.OpenTimeout / time.Second),
		TotalCalls:         cb.totalCalls.Load(),
		TotalFailures:      cb.totalFailures.Load(),
		TotalShortCircuits: cb.totalShortCircuits.Load(),
	}
	if cb.state == Open {
		u := cb.openUntil
		s.OpenUntil = &u
	}
	return s
}

func (cb *CircuitBreaker) RecentEvents() []Event {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	out := make([]Event, len(cb.events))
	copy(out, cb.events)
	return out
}
