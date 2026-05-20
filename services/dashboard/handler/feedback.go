package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	feedbackMaxBodyBytes   = 4 << 10 // 4 KiB
	feedbackMaxTextLen     = 2000
	feedbackMaxNameLen     = 80
	feedbackRateLimitMax   = 10
	feedbackRateLimitWin   = 5 * time.Minute
	discordContentMaxChars = 1900 // Discord-Limit ist 2000, mit Buffer für Header
)

type feedbackRequest struct {
	Thumb string `json:"thumb"`
	Text  string `json:"text"`
	Name  string `json:"name"`
}

type discordPayload struct {
	Content         string         `json:"content"`
	AllowedMentions map[string]any `json:"allowed_mentions"`
}

// FeedbackHandler nimmt Workshop-Feedback vom Dashboard entgegen und leitet
// es als formatierte Nachricht an einen Discord-Webhook weiter.
//
// Bei leerer webhookURL ist das Feature deaktiviert (503). So funktioniert
// das Dashboard auch ohne Webhook-Config — Teilnehmer:innen, die das Repo
// lokal klonen, sollen keinen kaputten Button sehen.
func FeedbackHandler(webhookURL string, now func() time.Time) http.HandlerFunc {
	if now == nil {
		now = time.Now
	}
	client := &http.Client{Timeout: 5 * time.Second}
	limiter := newFeedbackLimiter(feedbackRateLimitMax, feedbackRateLimitWin, now)

	return func(w http.ResponseWriter, r *http.Request) {
		if webhookURL == "" {
			http.Error(w, "feedback disabled", http.StatusServiceUnavailable)
			return
		}

		if !limiter.allow(clientIP(r)) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, feedbackMaxBodyBytes)
		var req feedbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		req.Text = strings.TrimSpace(req.Text)
		req.Name = strings.TrimSpace(req.Name)
		req.Thumb = strings.TrimSpace(req.Thumb)

		if req.Thumb != "up" && req.Thumb != "down" {
			http.Error(w, "thumb must be 'up' or 'down'", http.StatusBadRequest)
			return
		}
		if req.Text == "" {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}
		if len(req.Text) > feedbackMaxTextLen {
			http.Error(w, "text too long", http.StatusRequestEntityTooLarge)
			return
		}
		if len(req.Name) > feedbackMaxNameLen {
			http.Error(w, "name too long", http.StatusBadRequest)
			return
		}

		content := buildDiscordContent(req, now().UTC())
		body, _ := json.Marshal(discordPayload{
			Content:         content,
			AllowedMentions: map[string]any{"parse": []string{}},
		})

		resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
		if err != nil {
			http.Error(w, "discord unreachable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			io.Copy(io.Discard, resp.Body)
			http.Error(w, "discord rejected payload", http.StatusBadGateway)
			return
		}
		io.Copy(io.Discard, resp.Body)
		w.WriteHeader(http.StatusNoContent)
	}
}

// FeedbackStatusHandler signalisiert dem Frontend, ob ein Webhook konfiguriert
// ist. Frontend kann so auf einen Fallback (GitHub Discussions) ausweichen,
// wenn das Feature lokal deaktiviert ist.
func FeedbackStatusHandler(webhookURL string) http.HandlerFunc {
	enabled := webhookURL != ""
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"enabled": enabled})
	}
}

func buildDiscordContent(req feedbackRequest, t time.Time) string {
	emoji := "👍"
	if req.Thumb == "down" {
		emoji = "👎"
	}
	ts := t.Format("2006-01-02 15:04 MST")

	var header strings.Builder
	header.WriteString(emoji)
	header.WriteString("  *")
	header.WriteString(ts)
	header.WriteString("*")
	if req.Name != "" {
		header.WriteString("  —  *")
		header.WriteString(sanitizeForDiscord(req.Name))
		header.WriteString("*")
	}
	header.WriteString("\n")

	text := sanitizeForDiscord(req.Text)
	// Discord-Blockquote-Stil: jede Zeile mit "> " präfixen.
	var quoted strings.Builder
	for _, line := range strings.Split(text, "\n") {
		quoted.WriteString("> ")
		quoted.WriteString(line)
		quoted.WriteString("\n")
	}

	out := header.String() + quoted.String()
	if len(out) > discordContentMaxChars {
		out = out[:discordContentMaxChars] + "…"
	}
	return out
}

// sanitizeForDiscord verhindert, dass User-Input Discord-Mentions oder
// Code-Block-Ausbrüche triggert.
func sanitizeForDiscord(s string) string {
	s = strings.ReplaceAll(s, "@", "@​")
	s = strings.ReplaceAll(s, "```", "`​``")
	return s
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i >= 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type feedbackLimiter struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	now     func() time.Time
	buckets map[string][]time.Time
}

func newFeedbackLimiter(max int, window time.Duration, now func() time.Time) *feedbackLimiter {
	return &feedbackLimiter{
		max:     max,
		window:  window,
		now:     now,
		buckets: make(map[string][]time.Time),
	}
}

func (l *feedbackLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := l.now().Add(-l.window)
	hits := l.buckets[key]
	pruned := hits[:0]
	for _, t := range hits {
		if t.After(cutoff) {
			pruned = append(pruned, t)
		}
	}
	if len(pruned) >= l.max {
		l.buckets[key] = pruned
		return false
	}
	pruned = append(pruned, l.now())
	l.buckets[key] = pruned
	return true
}
