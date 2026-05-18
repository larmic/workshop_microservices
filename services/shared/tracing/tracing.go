// Package tracing implementiert eine minimale Variante des
// W3C-Trace-Context-Standards (https://www.w3.org/TR/trace-context/).
// Pflicht-Pfad für Story 7: Trace-ID aus `traceparent`-Header lesen oder
// neu generieren, durch alle ausgehenden Aufrufe propagieren und in jede
// strukturierte Logzeile schreiben.
//
// Vollwertige Span-Bäume (Parent/Child) und Backend-Export (Jaeger,
// Tempo, OTLP) sind bewusst out-of-scope — dafür existieren Libraries
// wie OpenTelemetry.
package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
)

const (
	HeaderName     = "traceparent"
	version        = "00"
	flagsSampled   = "01"
	traceIDLen     = 32
	spanIDLen      = 16
	flagsLen       = 2
	headerPartsLen = 4
)

type TraceContext struct {
	TraceID string
	SpanID  string
	Flags   string
}

func (tc TraceContext) Header() string {
	return version + "-" + tc.TraceID + "-" + tc.SpanID + "-" + tc.Flags
}

// New erzeugt einen frischen Trace-Kontext mit zufälliger Trace-ID und
// Span-ID. Sampling-Flag wird auf 01 (sampled) gesetzt.
func New() TraceContext {
	return TraceContext{
		TraceID: randomHex(16),
		SpanID:  randomHex(8),
		Flags:   flagsSampled,
	}
}

// Parse zerlegt einen `traceparent`-Header. Bei strikter
// Format-Validierung (Version 00, korrekte Längen, valide Hex-Zeichen,
// nicht-nullige IDs) wird der Kontext zurückgegeben — sonst (TraceContext{}, false).
func Parse(header string) (TraceContext, bool) {
	if header == "" {
		return TraceContext{}, false
	}
	parts := strings.Split(header, "-")
	if len(parts) != headerPartsLen {
		return TraceContext{}, false
	}
	if parts[0] != version {
		return TraceContext{}, false
	}
	traceID, spanID, flags := parts[1], parts[2], parts[3]
	if len(traceID) != traceIDLen || !isHex(traceID) || isAllZero(traceID) {
		return TraceContext{}, false
	}
	if len(spanID) != spanIDLen || !isHex(spanID) || isAllZero(spanID) {
		return TraceContext{}, false
	}
	if len(flags) != flagsLen || !isHex(flags) {
		return TraceContext{}, false
	}
	return TraceContext{TraceID: traceID, SpanID: spanID, Flags: flags}, true
}

type ctxKey struct{}

func WithContext(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, tc)
}

func FromContext(ctx context.Context) (TraceContext, bool) {
	tc, ok := ctx.Value(ctxKey{}).(TraceContext)
	return tc, ok
}

// Middleware liest einen eingehenden `traceparent`-Header oder erzeugt
// einen neuen Trace-Kontext und legt ihn in `r.Context()` ab. Setzt den
// Header zusätzlich auf der Response, damit Clients sehen, mit welcher
// Trace-ID ihr Request bearbeitet wurde.
//
// Geeignet für Entry-Point-Services (z.B. den Booking-Service), die als
// erste Station eines Vorgangs auch dann einen Trace starten sollen,
// wenn der Aufrufer keinen mitschickt. Downstream-Services nutzen
// stattdessen Propagate.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tc, ok := Parse(r.Header.Get(HeaderName))
		if !ok {
			tc = New()
		}
		w.Header().Set(HeaderName, tc.Header())
		next.ServeHTTP(w, r.WithContext(WithContext(r.Context(), tc)))
	})
}

// Propagate liest einen eingehenden `traceparent`-Header und legt ihn in
// `r.Context()` ab. Im Unterschied zu Middleware wird **kein** Trace
// erzeugt, wenn kein gültiger Header vorhanden ist — der Service bleibt
// dann ohne Trace-Kontext, Logger fällt auf den Default-Logger zurück
// und `trace_id` erscheint nicht in den Logs.
//
// Geeignet für Downstream-Services (Flight, Hotel, Car), die niemals
// selbst einen Trace initiieren sollen. Erst wenn der Entry-Point den
// Trace explizit propagiert, taucht die Trace-ID hier in den Logs auf.
func Propagate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tc, ok := Parse(r.Header.Get(HeaderName)); ok {
			w.Header().Set(HeaderName, tc.Header())
			r = r.WithContext(WithContext(r.Context(), tc))
		}
		next.ServeHTTP(w, r)
	})
}

// Inject setzt den `traceparent`-Header auf einem ausgehenden Request.
// Die Trace-ID bleibt erhalten, die Span-ID wird pro Hop neu erzeugt.
// Wenn im Context kein Trace-Kontext liegt, passiert nichts.
func Inject(ctx context.Context, req *http.Request) {
	tc, ok := FromContext(ctx)
	if !ok {
		return
	}
	hop := TraceContext{TraceID: tc.TraceID, SpanID: randomHex(8), Flags: tc.Flags}
	req.Header.Set(HeaderName, hop.Header())
}

// Logger liefert einen *slog.Logger, der trace_id und span_id als Felder
// vorgepinnt hat. Wenn im Context kein Trace-Kontext liegt, wird der
// Default-Logger ohne Trace-Felder zurückgegeben.
func Logger(ctx context.Context) *slog.Logger {
	tc, ok := FromContext(ctx)
	if !ok {
		return slog.Default()
	}
	return slog.Default().With(
		slog.String("trace_id", tc.TraceID),
		slog.String("span_id", tc.SpanID),
	)
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		// crypto/rand kann in der Praxis nicht fehlschlagen; falls doch,
		// liefern wir einen festen Wert, der nie als gültige ID akzeptiert
		// wird (Parse() lehnt all-zero ab). Damit fällt der Fehler auf.
		return strings.Repeat("0", n*2)
	}
	return hex.EncodeToString(buf)
}

func isHex(s string) bool {
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		default:
			return false
		}
	}
	return true
}

func isAllZero(s string) bool {
	for _, r := range s {
		if r != '0' {
			return false
		}
	}
	return true
}
