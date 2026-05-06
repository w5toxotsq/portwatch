package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	c := metrics.New()
	c.RecordPoll(30*time.Millisecond, 2, 1)
	c.RecordAlert()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler(c).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	assertField := func(key string, want any) {
		t.Helper()
		got, ok := body[key]
		if !ok {
			t.Errorf("missing field %q", key)
			return
		}
		// JSON numbers decode as float64
		if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
			t.Errorf("field %q: got %v, want %v", key, got, want)
		}
	}

	assertField("poll_count", float64(1))
	assertField("opened_total", float64(2))
	assertField("closed_total", float64(1))
	assertField("alerts_sent", float64(1))
	assertField("last_poll_duration_ms", float64(30))

	if _, ok := body["uptime_seconds"]; !ok {
		t.Error("missing field uptime_seconds")
	}
	if _, ok := body["last_poll_at"]; !ok {
		t.Error("missing field last_poll_at")
	}
}

func TestHandler_ContentTypeJSON(t *testing.T) {
	c := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metrics.Handler(c).ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
