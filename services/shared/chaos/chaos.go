package chaos

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultLatencyMs = 2000

type Mode string

const (
	Normal Mode = "normal"
	Slow   Mode = "slow"
	Fail   Mode = "fail"
)

type State struct {
	Mode      Mode `json:"mode"`
	LatencyMs int  `json:"latencyMs"`
}

// Snapshot ist das, was über GET /admin/chaos nach außen geht. Enthält
// zusätzlich zum konfigurierten State den Zeitstempel des letzten von
// einem Aufrufer beobachteten Requests (für die Backend-Card im Dashboard,
// die kurz pulst, wenn eine Replica gerade gerufen wurde).
type Snapshot struct {
	State
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
}

type Chaos struct {
	mu         sync.RWMutex
	state      State
	lastSeenAt time.Time
}

func New() *Chaos {
	mode := Mode(strings.ToLower(os.Getenv("CHAOS_DEFAULT_MODE")))
	if mode != Slow && mode != Fail {
		mode = Normal
	}
	return &Chaos{state: State{Mode: mode, LatencyMs: latencyFromEnv()}}
}

// latencyFromEnv liest die Slow-Latenz aus CHAOS_LATENCY_MS (z. B. via docker-compose).
// Fix für die Lebensdauer der Instanz; nicht zur Laufzeit änderbar.
func latencyFromEnv() int {
	raw := os.Getenv("CHAOS_LATENCY_MS")
	if raw == "" {
		return defaultLatencyMs
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return defaultLatencyMs
	}
	return n
}

func (c *Chaos) Snapshot() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// snapshotWithLastSeen erstellt die nach außen sichtbare Sicht für GET
// /admin/chaos: aktueller State + Zeitpunkt des letzten "regulären"
// Requests (nicht-whitelisted) gegen die Instanz.
func (c *Chaos) snapshotWithLastSeen() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snap := Snapshot{State: c.state}
	if !c.lastSeenAt.IsZero() {
		t := c.lastSeenAt
		snap.LastSeenAt = &t
	}
	return snap
}

// markSeen merkt sich, dass die Instanz gerade einen regulären Request
// beantwortet. Wird vom Middleware-Pfad atomar gesetzt.
func (c *Chaos) markSeen(now time.Time) {
	c.mu.Lock()
	c.lastSeenAt = now
	c.mu.Unlock()
}

func (c *Chaos) SetMode(mode Mode) State {
	if mode != Normal && mode != Slow && mode != Fail {
		mode = Normal
	}
	c.mu.Lock()
	c.state.Mode = mode
	s := c.state
	c.mu.Unlock()

	hostname, _ := os.Hostname()
	log.Printf("[%s] Chaos mode set to %s (latency=%dms fix)", hostname, s.Mode, s.LatencyMs)
	return s
}

// Health-Check und Admin-Pfade müssen unbeeinflusst bleiben — sonst deregistriert
// Consul die Instanz im Fail-Mode, bevor der Circuit Breaker triggern kann.
func isWhitelisted(path string) bool {
	if path == "/health" || path == "/info" || path == "/openapi" {
		return true
	}
	return strings.HasPrefix(path, "/admin/")
}

func (c *Chaos) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWhitelisted(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		c.markSeen(time.Now().UTC())
		s := c.Snapshot()
		switch s.Mode {
		case Slow:
			time.Sleep(time.Duration(s.LatencyMs) * time.Millisecond)
		case Fail:
			http.Error(w, "service intentionally failing (chaos)", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *Chaos) GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(c.snapshotWithLastSeen())
}

func (c *Chaos) SetHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Mode Mode `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	updated := c.SetMode(req.Mode)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}
