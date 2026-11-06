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

	"github.com/team-neusta-skills/workshop_microservices/shared/tracing"
)

const (
	defaultLatencyMs = 2000
	// Hängen muss länger dauern als jeder Client-Timeout im Stack
	// (3 s ab Story 4, 5 s in Story 1 und 3), sonst ist es nur langsam.
	defaultHangMs = 10000
)

type Mode string

const (
	Normal Mode = "normal"
	Slow   Mode = "slow"
	Fail   Mode = "fail"
	// Hang antwortet erst nach HangMs, also nach dem Client-Timeout. Die
	// Anfrage wird danach trotzdem verarbeitet: Der Aufrufer hat aufgegeben,
	// der Service bucht (Story 6: "Timeout, aber gebucht").
	Hang Mode = "hang"
)

type State struct {
	Mode      Mode `json:"mode"`
	LatencyMs int  `json:"latencyMs"`
	HangMs    int  `json:"hangMs"`
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
	if !mode.valid() {
		mode = Normal
	}
	return &Chaos{state: State{
		Mode:      mode,
		LatencyMs: msFromEnv("CHAOS_LATENCY_MS", defaultLatencyMs),
		HangMs:    msFromEnv("CHAOS_HANG_MS", defaultHangMs),
	}}
}

func (m Mode) valid() bool {
	return m == Normal || m == Slow || m == Fail || m == Hang
}

// msFromEnv liest eine Wartezeit in Millisekunden (CHAOS_LATENCY_MS für
// Slow, CHAOS_HANG_MS für Hang, z. B. via docker-compose). Fix für die
// Lebensdauer der Instanz; nicht zur Laufzeit änderbar.
func msFromEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
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
	if !mode.valid() {
		mode = Normal
	}
	c.mu.Lock()
	c.state.Mode = mode
	s := c.state
	c.mu.Unlock()

	hostname, _ := os.Hostname()
	log.Printf("[%s] Chaos mode set to %s (latency=%dms, hang=%dms fix)", hostname, s.Mode, s.LatencyMs, s.HangMs)
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
		case Hang:
			// Bewusst ohne Blick auf r.Context(): Auch wenn der Aufrufer
			// längst abgebrochen hat, läuft der Handler danach normal weiter.
			time.Sleep(time.Duration(s.HangMs) * time.Millisecond)
			tracing.Logger(r.Context()).Warn("request answered after hang, caller may have given up", "method", r.Method, "path", r.URL.Path, "hangMs", s.HangMs)
		case Fail:
			// Loggen, bevor wir abbrechen: Sonst fehlt genau der Fehlerfall
			// im Trace (Story 8).
			tracing.Logger(r.Context()).Warn("request failed by chaos", "method", r.Method, "path", r.URL.Path)
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
