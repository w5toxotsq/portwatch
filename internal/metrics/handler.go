package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Handler returns an http.Handler that serves current metrics as JSON.
// It is intended to be mounted on a lightweight debug HTTP server.
func Handler(c *Collector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := c.Get()
		payload := map[string]any{
			"poll_count":           s.PollCount,
			"last_poll_at":         s.LastPollAt.Format(time.RFC3339),
			"last_poll_duration_ms": s.LastPollDur.Milliseconds(),
			"opened_total":         s.OpenedTotal,
			"closed_total":         s.ClosedTotal,
			"alerts_sent":          s.AlertsSent,
			"uptime_seconds":       int64(time.Since(s.UptimeSince).Seconds()),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
		}
	})
}
