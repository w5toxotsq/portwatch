package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func writeTempHistoryFile(t *testing.T, h *history.History) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	if err := h.Save(path); err != nil {
		t.Fatalf("save history: %v", err)
	}
	return path
}

func buildTestHistory(t *testing.T) *history.History {
	t.Helper()
	h := history.New()
	h.Record(time.Now().Add(-2*time.Hour), history.Changes{
		Opened: []string{"tcp:8080"},
		Closed: []string{},
	})
	h.Record(time.Now().Add(-1*time.Hour), history.Changes{
		Opened: []string{},
		Closed: []string{"tcp:8080"},
	})
	return h
}

func TestRunHistory_TextOutput(t *testing.T) {
	h := buildTestHistory(t)
	path := writeTempHistoryFile(t, h)

	cmd := newHistoryCmd()
	cmd.SetArgs([]string{"--file", path, "--format", "text"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tcp:8080") {
		t.Errorf("expected port in output, got:\n%s", out)
	}
}

func TestRunHistory_JSONOutput(t *testing.T) {
	h := buildTestHistory(t)
	path := writeTempHistoryFile(t, h)

	cmd := newHistoryCmd()
	cmd.SetArgs([]string{"--file", path, "--format", "json"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}
}

func TestRunHistory_MissingFile(t *testing.T) {
	cmd := newHistoryCmd()
	cmd.SetArgs([]string{"--file", "/nonexistent/path/history.json"})
	cmd.SetOut(os.Stdout)

	if err := cmd.Execute(); err != nil {
		t.Errorf("expected graceful handling of missing file, got error: %v", err)
	}
}

func TestRunHistory_SinceFilter(t *testing.T) {
	h := buildTestHistory(t)
	path := writeTempHistoryFile(t, h)

	sinceStr := time.Now().Add(-90 * time.Minute).Format(time.RFC3339)

	cmd := newHistoryCmd()
	cmd.SetArgs([]string{"--file", path, "--since", sinceStr, "--format", "text"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
