package metrics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestExport_TextFormat(t *testing.T) {
	m := New()
	m.RecordPoll(false)
	m.RecordPoll(true)
	m.RecordAlert()

	var buf bytes.Buffer
	err := Export(m, ExportOptions{Format: ExportText, Writer: &buf})
	if err != nil {
		t.Fatalf("Export returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Uptime", "Total Polls", "Failed Polls", "Alerts Sent", "Collected At"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestExport_JSONFormat(t *testing.T) {
	m := New()
	m.RecordPoll(false)
	m.RecordPoll(false)
	m.RecordPoll(true)
	m.RecordAlert()
	m.RecordAlert()

	var buf bytes.Buffer
	err := Export(m, ExportOptions{Format: ExportJSON, Writer: &buf})
	if err != nil {
		t.Fatalf("Export returned error: %v", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(buf.Bytes(), &snap); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if snap.TotalPolls != 3 {
		t.Errorf("TotalPolls = %d, want 3", snap.TotalPolls)
	}
	if snap.FailedPolls != 1 {
		t.Errorf("FailedPolls = %d, want 1", snap.FailedPolls)
	}
	if snap.AlertsSent != 2 {
		t.Errorf("AlertsSent = %d, want 2", snap.AlertsSent)
	}
	if snap.CollectedAt.IsZero() {
		t.Error("CollectedAt should not be zero")
	}
}

func TestExport_DefaultsToText(t *testing.T) {
	m := New()
	var buf bytes.Buffer
	err := Export(m, ExportOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("Export returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "Total Polls") {
		t.Error("default format should produce text output")
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		input time.Duration
		want  string
	}{
		{45 * time.Second, "45s"},
		{2*time.Minute + 30*time.Second, "2m 30s"},
		{3*time.Hour + 5*time.Minute + 10*time.Second, "3h 5m 10s"},
	}
	for _, tc := range cases {
		got := formatDuration(tc.input)
		if got != tc.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
