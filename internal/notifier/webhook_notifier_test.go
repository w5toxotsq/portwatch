package notifier_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/snapshot"
)

func buildDiff(opened, closed []snapshot.PortEntry) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestWebhookNotifier_PostsOnChange(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.NewWebhookNotifier(ts.URL, ts.Client())
	diff := buildDiff(
		[]snapshot.PortEntry{{Proto: "tcp", Port: 8080}},
		[]snapshot.PortEntry{{Proto: "tcp", Port: 22}},
	)

	if err := n.Notify(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opened, _ := received["opened"].([]interface{})
	if len(opened) != 1 || opened[0] != "tcp:8080" {
		t.Errorf("opened = %v, want [tcp:8080]", opened)
	}

	closed, _ := received["closed"].([]interface{})
	if len(closed) != 1 || closed[0] != "tcp:22" {
		t.Errorf("closed = %v, want [tcp:22]", closed)
	}

	if received["timestamp"] == nil {
		t.Error("expected timestamp field")
	}
}

func TestWebhookNotifier_SkipsEmptyDiff(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.NewWebhookNotifier(ts.URL, ts.Client())
	if err := n.Notify(snapshot.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP call for empty diff")
	}
}

func TestWebhookNotifier_ReturnsErrorOnNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notifier.NewWebhookNotifier(ts.URL, ts.Client())
	diff := buildDiff([]snapshot.PortEntry{{Proto: "tcp", Port: 9000}}, nil)

	if err := n.Notify(diff); err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestWebhookNotifier_DefaultClient(t *testing.T) {
	// Passing nil should not panic — a default client is created internally.
	n := notifier.NewWebhookNotifier("http://localhost:0", nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
