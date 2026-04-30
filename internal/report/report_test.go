package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
	"github.com/user/portwatch/internal/snapshot"
)

func buildHistory(t *testing.T) *history.History {
	t.Helper()
	h := history.New()
	now := time.Now()
	h.Record(snapshot.Changes{
		Opened: []string{"tcp:8080"},
		Closed: []string{},
	}, now.Add(-2*time.Hour))
	h.Record(snapshot.Changes{
		Opened: []string{},
		Closed: []string{"tcp:8080"},
	}, now.Add(-1*time.Hour))
	return h
}

func TestGenerate_TextFormat(t *testing.T) {
	h := buildHistory(t)
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.Writer = &buf

	if err := report.Generate(h, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("expected header row in text output")
	}
	if !strings.Contains(out, "1") {
		t.Error("expected port counts in output")
	}
}

func TestGenerate_JSONFormat(t *testing.T) {
	h := buildHistory(t)
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.Format = report.FormatJSON
	opts.Writer = &buf

	if err := report.Generate(h, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestGenerate_EmptyHistory(t *testing.T) {
	h := history.New()
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.Writer = &buf

	if err := report.Generate(h, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Error("expected empty message")
	}
}

func TestGenerate_LimitEntries(t *testing.T) {
	h := buildHistory(t)
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.Limit = 1
	opts.Writer = &buf

	if err := report.Generate(h, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 1 data line
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header+1), got %d", len(lines))
	}
}

func TestGenerate_SinceFilter(t *testing.T) {
	h := buildHistory(t)
	var buf bytes.Buffer
	opts := report.DefaultOptions()
	opts.Since = time.Now().Add(-90 * time.Minute)
	opts.Writer = &buf

	if err := report.Generate(h, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header+1 recent), got %d", len(lines))
	}
}
