package chaos

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Hang muss den Client-Timeout überdauern und die Anfrage danach trotzdem
// verarbeiten, inklusive Body: Das ist "Timeout, aber gebucht" aus Story 6.
func TestHangAnswersAfterClientTimeoutAndStillHandlesRequest(t *testing.T) {
	c := &Chaos{state: State{Mode: Hang, HangMs: 300}}

	handled := make(chan string, 1)
	srv := httptest.NewServer(c.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			CustomerName string `json:"customerName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handled <- "decode failed: " + err.Error()
			return
		}
		handled <- body.CustomerName
		w.WriteHeader(http.StatusCreated)
	})))
	defer srv.Close()

	client := &http.Client{Timeout: 100 * time.Millisecond}
	payload, _ := json.Marshal(map[string]string{"customerName": "Max"})
	resp, err := client.Post(srv.URL+"/bookings", "application/json", bytes.NewReader(payload))
	if err == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		t.Fatalf("expected client timeout, got status %d", resp.StatusCode)
	}

	select {
	case got := <-handled:
		if got != "Max" {
			t.Fatalf("handler after hang: got %q, want %q", got, "Max")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("handler never ran after hang")
	}
}

func TestHangSkipsHealthAndAdmin(t *testing.T) {
	c := &Chaos{state: State{Mode: Hang, HangMs: 5000}}
	srv := httptest.NewServer(c.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()

	client := &http.Client{Timeout: 500 * time.Millisecond}
	for _, path := range []string{"/health", "/admin/chaos"} {
		resp, err := client.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("%s should not hang: %v", path, err)
		}
		resp.Body.Close()
	}
}

func TestSetModeAcceptsHang(t *testing.T) {
	c := &Chaos{state: State{Mode: Normal}}
	if got := c.SetMode(Hang).Mode; got != Hang {
		t.Fatalf("SetMode(hang) = %q, want %q", got, Hang)
	}
	if got := c.SetMode("bogus").Mode; got != Normal {
		t.Fatalf("SetMode(bogus) = %q, want %q", got, Normal)
	}
}
