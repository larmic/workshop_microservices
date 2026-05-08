package saga

import (
	"encoding/json"
	"sort"
	"sync"
	"time"
)

type Status string

const (
	StatusPending      Status = "PENDING"
	StatusCompleted    Status = "COMPLETED"
	StatusCompensating Status = "COMPENSATING"
	StatusFailed       Status = "FAILED"
)

type StepStatus string

const (
	StepBooked             StepStatus = "BOOKED"
	StepCompensated        StepStatus = "COMPENSATED"
	StepCompensationFailed StepStatus = "COMPENSATION_FAILED"
)

type Step struct {
	Service   string          `json:"service"`
	BookingID string          `json:"bookingId,omitempty"`
	Status    StepStatus      `json:"status"`
	Detail    json.RawMessage `json:"detail,omitempty"`
	Reason    string          `json:"reason,omitempty"`
}

type Saga struct {
	SagaID       string    `json:"sagaId"`
	CustomerName string    `json:"customerName"`
	Status       Status    `json:"status"`
	Steps        []Step    `json:"steps"`
	FailedAt     string    `json:"failedAt,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Store hält Sagas in-memory. Nach Crash sind sie weg — eine echte
// Recovery braucht eine persistente Ablage (siehe docs/questions/story5.md,
// Frage 5 zu Monitoring und Story 5 AK „Saga-Status persistieren").
// Für den Workshop reicht in-memory.
type Store struct {
	mu    sync.RWMutex
	sagas map[string]*Saga
}

func NewStore() *Store {
	return &Store{sagas: make(map[string]*Saga)}
}

func (s *Store) Save(saga Saga) {
	s.mu.Lock()
	defer s.mu.Unlock()
	saga.UpdatedAt = time.Now().UTC()
	stored := saga
	if len(saga.Steps) > 0 {
		stored.Steps = make([]Step, len(saga.Steps))
		copy(stored.Steps, saga.Steps)
	}
	s.sagas[saga.SagaID] = &stored
}

func (s *Store) Get(id string) (Saga, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stored, ok := s.sagas[id]
	if !ok {
		return Saga{}, false
	}
	out := *stored
	if len(stored.Steps) > 0 {
		out.Steps = make([]Step, len(stored.Steps))
		copy(out.Steps, stored.Steps)
	}
	return out, true
}

// List liefert alle Sagas, neueste zuerst. Snapshot — Mutationen
// am Ergebnis wirken sich nicht auf den Store aus.
func (s *Store) List() []Saga {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Saga, 0, len(s.sagas))
	for _, stored := range s.sagas {
		copy := *stored
		if len(stored.Steps) > 0 {
			copy.Steps = append([]Step(nil), stored.Steps...)
		}
		out = append(out, copy)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// Reset leert den Store. Workshop-Helfer für saubere Demo-Runden.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sagas = make(map[string]*Saga)
}
