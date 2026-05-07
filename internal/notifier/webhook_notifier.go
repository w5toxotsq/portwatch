package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// WebhookNotifier sends port change alerts as JSON POST requests to a URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Timestamp string   `json:"timestamp"`
	Opened    []string `json:"opened"`
	Closed    []string `json:"closed"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
// An optional *http.Client may be provided; if nil, a default client with a
// 10-second timeout is used.
func NewWebhookNotifier(url string, client *http.Client) *WebhookNotifier {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &WebhookNotifier{url: url, client: client}
}

// Notify serialises the diff as JSON and POSTs it to the configured URL.
// It returns an error if marshalling, the HTTP request, or a non-2xx
// response status is encountered.
func (w *WebhookNotifier) Notify(diff snapshot.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}

	payload := webhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Opened:    formatPorts(diff.Opened),
		Closed:    formatPorts(diff.Closed),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// formatPorts converts a slice of scanner.Port to "proto:port" strings.
func formatPorts(ports []snapshot.PortEntry) []string {
	out := make([]string, len(ports))
	for i, p := range ports {
		out[i] = fmt.Sprintf("%s:%d", p.Proto, p.Port)
	}
	return out
}
