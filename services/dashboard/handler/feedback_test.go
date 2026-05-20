package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func fixedClock() func() time.Time {
	t := time.Date(2026, 5, 20, 14, 32, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func TestFeedbackHandler_DisabledWhenWebhookEmpty(t *testing.T) {
	h := FeedbackHandler("", fixedClock())
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"up","text":"hi"}`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", w.Code)
	}
}

func TestFeedbackHandler_InvalidJSON(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{not json`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestFeedbackHandler_RejectsInvalidThumb(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"maybe","text":"hi"}`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestFeedbackHandler_RequiresText(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"up","text":"  "}`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestFeedbackHandler_TextTooLong(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	// 4 KiB body limit kommt vor der Text-Längen-Prüfung — wir testen den
	// Body-Limit-Path explizit über payload > 4 KiB.
	big := strings.Repeat("a", 5000)
	h := FeedbackHandler(upstream.URL, fixedClock())
	payload := `{"thumb":"up","text":"` + big + `"}`
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(payload))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("want 413, got %d", w.Code)
	}
}

func TestFeedbackHandler_HappyPath_PostsFormattedContent(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	body := `{"thumb":"down","text":"Story 3 zu schnell.","name":"Max"}`
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(body))
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d, body=%s", w.Code, w.Body.String())
	}
	if upstream.calls != 1 {
		t.Fatalf("want 1 upstream call, got %d", upstream.calls)
	}
	if upstream.lastContentType != "application/json" {
		t.Fatalf("want application/json content-type, got %q", upstream.lastContentType)
	}

	var payload discordPayload
	if err := json.Unmarshal([]byte(upstream.lastBody), &payload); err != nil {
		t.Fatalf("upstream body not JSON: %v", err)
	}
	if !strings.Contains(payload.Content, "👎") {
		t.Errorf("content missing thumb emoji: %q", payload.Content)
	}
	if !strings.Contains(payload.Content, "2026-05-20 14:32 UTC") {
		t.Errorf("content missing timestamp: %q", payload.Content)
	}
	if !strings.Contains(payload.Content, "Max") {
		t.Errorf("content missing name: %q", payload.Content)
	}
	if !strings.Contains(payload.Content, "> Story 3 zu schnell.") {
		t.Errorf("content missing quoted text: %q", payload.Content)
	}
	mentions, _ := payload.AllowedMentions["parse"].([]any)
	if len(mentions) != 0 {
		t.Errorf("want empty allowed_mentions parse, got %v", mentions)
	}
}

func TestFeedbackHandler_DiscordError_ReturnsBadGateway(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusBadRequest)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"up","text":"hi"}`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusBadGateway {
		t.Fatalf("want 502, got %d", w.Code)
	}
}

func TestFeedbackHandler_RateLimit(t *testing.T) {
	upstream := newRecordingUpstream(http.StatusNoContent)
	defer upstream.Close()

	h := FeedbackHandler(upstream.URL, fixedClock())
	for i := 0; i < feedbackRateLimitMax; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"up","text":"ok"}`))
		w := httptest.NewRecorder()
		h(w, r)
		if w.Code != http.StatusNoContent {
			t.Fatalf("req %d want 204, got %d", i, w.Code)
		}
	}
	// (max+1)-ter Request muss geblockt werden
	r := httptest.NewRequest(http.MethodPost, "/api/feedback", strings.NewReader(`{"thumb":"up","text":"ok"}`))
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", w.Code)
	}
}

func TestFeedbackStatusHandler(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"", false},
		{"https://discord.com/api/webhooks/x/y", true},
	}
	for _, c := range cases {
		h := FeedbackStatusHandler(c.url)
		r := httptest.NewRequest(http.MethodGet, "/api/feedback/status", nil)
		w := httptest.NewRecorder()
		h(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("url=%q want 200, got %d", c.url, w.Code)
		}
		var resp map[string]bool
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("body not JSON: %v", err)
		}
		if resp["enabled"] != c.want {
			t.Errorf("url=%q want enabled=%v, got %v", c.url, c.want, resp["enabled"])
		}
	}
}

func TestSanitizeForDiscord(t *testing.T) {
	if got := sanitizeForDiscord("@everyone hi"); !strings.HasPrefix(got, "@") || strings.HasPrefix(got, "@everyone") {
		t.Errorf("at-mention should be neutralized, got %q", got)
	}
	if got := sanitizeForDiscord("```evil```"); strings.Contains(got, "```") {
		t.Errorf("triple-backticks should be neutralized, got %q", got)
	}
}

type recordingUpstream struct {
	*httptest.Server
	calls           int
	lastBody        string
	lastContentType string
}

func newRecordingUpstream(status int) *recordingUpstream {
	u := &recordingUpstream{}
	u.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.calls++
		u.lastContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		u.lastBody = string(b)
		w.WriteHeader(status)
	}))
	return u
}
