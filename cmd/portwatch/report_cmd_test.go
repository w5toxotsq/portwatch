package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
)

func writeTempHistory(t *testing.T, h *history.History) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	if err := h.Save(path); err != nil {
		t.Fatalf("failed to write temp history: %v", err)
	}
	return path
}

func TestRunReport_TextOutput(t *testing.T) {
	h := history.New()
	h.Record(snapshot.Changes{
		Opened: []string{"tcp:9090"},
		Closed: []string{},
	}, time.Now())
	path := writeTempHistory(t, h)

	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runReport([]string{"--file", path, "--format", "text", "--limit", "10"})

	w.Close()
	os.Stdout = old

	var buf [4096]byte
	n, _ := r.Read(buf[:])
	out := string(buf[:n])
	if len(out) == 0 {
		t.Error("expected non-empty text output")
	}
}

func TestRunReport_JSONOutput(t *testing.T) {
	h := history.New()
	h.Record(snapshot.Changes{
		Opened: []string{"tcp:7070"},
		Closed: []string{},
	}, time.Now())
	path := writeTempHistory(t, h)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runReport([]string{"--file", path, "--format", "json"})

	w.Close()
	os.Stdout = old

	var buf [4096]byte
	n, _ := r.Read(buf[:])

	var entries []map[string]interface{}
	if err := json.Unmarshal(buf[:n], &entries); err != nil {
		t.Fatalf("expected valid JSON output: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestRunReport_MissingFile(t *testing.T) {
	// Should not panic; missing history file is treated as empty.
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
	}()

	runReport([]string{"--file", "/nonexistent/path/history.json"})
}
