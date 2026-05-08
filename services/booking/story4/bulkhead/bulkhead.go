package bulkhead

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
)

var ErrBulkheadFull = errors.New("bulkhead is full")

type Config struct {
	Name          string
	MaxConcurrent int
}

type Snapshot struct {
	Name          string `json:"name"`
	MaxConcurrent int    `json:"maxConcurrent"`
	InFlight      int    `json:"inFlight"`
	TotalCalls    uint64 `json:"totalCalls"`
	TotalRejected uint64 `json:"totalRejected"`
}

type Bulkhead struct {
	cfg Config
	sem chan struct{}

	inFlight      atomic.Int64
	totalCalls    atomic.Uint64
	totalRejected atomic.Uint64
}

func New(cfg Config) *Bulkhead {
	if cfg.Name == "" {
		cfg.Name = "default"
	}
	if cfg.MaxConcurrent <= 0 {
		cfg.MaxConcurrent = 10
	}
	return &Bulkhead{
		cfg: cfg,
		sem: make(chan struct{}, cfg.MaxConcurrent),
	}
}

// Execute belegt nicht-blockierend einen Slot im Semaphore. Ist der Pool voll,
// wird der Aufruf SOFORT abgewiesen (kein Queueing) — genau das ist der Sinn
// des Bulkheads: ein langsames Backend darf nicht beliebig viele Aufrufe
// stauen, weil sonst der gesamte Booking-Service blockiert.
func (b *Bulkhead) Execute(ctx context.Context, fn func(context.Context) error) error {
	select {
	case b.sem <- struct{}{}:
	default:
		b.totalRejected.Add(1)
		log.Printf("BH[%s] FULL — Call REJECTED (inFlight=%d/%d)",
			b.cfg.Name, b.inFlight.Load(), b.cfg.MaxConcurrent)
		return ErrBulkheadFull
	}

	b.totalCalls.Add(1)
	b.inFlight.Add(1)
	defer func() {
		b.inFlight.Add(-1)
		<-b.sem
	}()

	return fn(ctx)
}

func (b *Bulkhead) Snapshot() Snapshot {
	return Snapshot{
		Name:          b.cfg.Name,
		MaxConcurrent: b.cfg.MaxConcurrent,
		InFlight:      int(b.inFlight.Load()),
		TotalCalls:    b.totalCalls.Load(),
		TotalRejected: b.totalRejected.Load(),
	}
}
