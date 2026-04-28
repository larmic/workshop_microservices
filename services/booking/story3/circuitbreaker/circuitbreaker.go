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

// IsFailure entscheidet, ob ein Fehler den Circuit-Breaker-Counter erhöhen soll.
// 4xx-Antworten zählen nicht — die deuten auf Client-Fehler hin, nicht auf einen
// kaputten Backend-Service.
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

	if cb.state == Open {
		if time.Now().After(cb.openUntil) {
			cb.transitionLocked(HalfOpen, "open timeout elapsed")
		} else {
			return false
		}
	}

	if cb.state == HalfOpen {
		// Nur ein Probe-Call gleichzeitig zulassen.
		return cb.probeInFlight.CompareAndSwap(false, true)
	}

	return cb.state == Closed
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
			cb.transitionLocked(Open, "probe failed in half-open")
		case Closed:
			if cb.failures >= cb.cfg.FailureThreshold {
				cb.openUntil = time.Now().Add(cb.cfg.OpenTimeout)
				cb.transitionLocked(Open, "failure threshold reached")
			}
		}
		return
	}

	switch cb.state {
	case HalfOpen:
		cb.probeInFlight.Store(false)
		cb.failures = 0
		cb.transitionLocked(Closed, "probe succeeded")
	case Closed:
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

	log.Printf("CB[%s] %s → %s (reason=%s)", cb.cfg.Name, from, to, reason)

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
