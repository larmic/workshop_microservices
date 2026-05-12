package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseValid(t *testing.T) {
	header := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
	tc, ok := Parse(header)
	if !ok {
		t.Fatalf("expected valid parse, got ok=false")
	}
	if tc.TraceID != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("trace id mismatch: %q", tc.TraceID)
	}
	if tc.SpanID != "b7ad6b7169203331" {
		t.Errorf("span id mismatch: %q", tc.SpanID)
	}
	if tc.Flags != "01" {
		t.Errorf("flags mismatch: %q", tc.Flags)
	}
}

func TestParseInvalid(t *testing.T) {
	cases := map[string]string{
		"empty":              "",
		"wrong version":      "01-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01",
		"too short trace":    "00-0af7-b7ad6b7169203331-01",
		"too short span":     "00-0af7651916cd43dd8448eb211c80319c-b7ad-01",
		"non-hex trace":      "00-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz-b7ad6b7169203331-01",
		"all-zero trace":     "00-00000000000000000000000000000000-b7ad6b7169203331-01",
		"all-zero span":      "00-0af7651916cd43dd8448eb211c80319c-0000000000000000-01",
		"missing flags":      "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331",
		"extra field":        "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01-extra",
		"non-hex flags":      "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-zz",
		"uppercase rejected": "00-0AF7651916CD43DD8448EB211C80319C-b7ad6b7169203331-01",
	}
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			if _, ok := Parse(header); ok {
				t.Errorf("expected invalid, but parse succeeded")
			}
		})
	}
}

func TestNewRoundTrip(t *testing.T) {
	tc := New()
	parsed, ok := Parse(tc.Header())
	if !ok {
		t.Fatalf("round-trip failed: header %q rejected by Parse", tc.Header())
	}
	if parsed != tc {
		t.Errorf("round-trip mismatch: in=%+v out=%+v", tc, parsed)
	}
}

func TestNewUniqueness(t *testing.T) {
	a := New()
	b := New()
	if a.TraceID == b.TraceID {
		t.Errorf("two consecutive New() calls produced same trace id %q", a.TraceID)
	}
	if a.SpanID == b.SpanID {
		t.Errorf("two consecutive New() calls produced same span id %q", a.SpanID)
	}
}

func TestMiddlewareGeneratesWhenMissing(t *testing.T) {
	var captured TraceContext
	var hadCtx bool
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, hadCtx = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !hadCtx {
		t.Fatal("handler did not see trace context in request")
	}
	if captured.TraceID == "" {
		t.Error("expected generated trace id, got empty")
	}
	if got := rec.Header().Get(HeaderName); got != captured.Header() {
		t.Errorf("response header %q does not match context %q", got, captured.Header())
	}
}

func TestMiddlewareAdoptsIncoming(t *testing.T) {
	incoming := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
	var captured TraceContext
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderName, incoming)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if captured.TraceID != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("trace id not adopted from header: got %q", captured.TraceID)
	}
}

func TestMiddlewareRegeneratesOnInvalidHeader(t *testing.T) {
	var captured TraceContext
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderName, "garbage")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if captured.TraceID == "" {
		t.Error("expected fresh trace context after invalid header, got empty")
	}
	if strings.Contains(captured.TraceID, "garbage") {
		t.Error("captured context contains garbage from invalid header")
	}
}

func TestInjectKeepsTraceIDChangesSpanID(t *testing.T) {
	original := TraceContext{
		TraceID: "0af7651916cd43dd8448eb211c80319c",
		SpanID:  "b7ad6b7169203331",
		Flags:   "01",
	}
	ctx := WithContext(context.Background(), original)
	req := httptest.NewRequest(http.MethodGet, "http://example/", nil)
	Inject(ctx, req)

	hop, ok := Parse(req.Header.Get(HeaderName))
	if !ok {
		t.Fatal("inject produced invalid traceparent header")
	}
	if hop.TraceID != original.TraceID {
		t.Errorf("trace id changed across hop: %q -> %q", original.TraceID, hop.TraceID)
	}
	if hop.SpanID == original.SpanID {
		t.Errorf("span id should change per hop, but stayed %q", hop.SpanID)
	}
}

func TestInjectNoopWithoutContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example/", nil)
	Inject(context.Background(), req)
	if req.Header.Get(HeaderName) != "" {
		t.Error("inject set header even though no trace context was in ctx")
	}
}
