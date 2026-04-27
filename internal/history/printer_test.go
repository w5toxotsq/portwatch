package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_ShowsEntries(t *testing.T) {
	h := &History{
		Entries: []Entry{
			{
				Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Opened:    []string{"tcp/8080"},
				Closed:    []string{"tcp/22"},
			},
		},
	}

	var buf bytes.Buffer
	opts := DefaultPrintOptions()
	opts.Writer = &buf

	if err := h.Print(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "OPENED: tcp/8080") {
		t.Errorf("expected OPENED line, got: %s", out)
	}
	if !strings.Contains(out, "CLOSED: tcp/22") {
		t.Errorf("expected CLOSED line, got: %s", out)
	}
}

func TestPrint_EmptyHistory(t *testing.T) {
	h := &History{}

	var buf bytes.Buffer
	opts := DefaultPrintOptions()
	opts.Writer = &buf

	if err := h.Print(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No history entries found.") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrint_LimitEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	h := &History{}
	for i := 0; i < 5; i++ {
		h.Entries = append(h.Entries, Entry{
			Timestamp: base.Add(time.Duration(i) * time.Hour),
			Opened:    []string{"tcp/808" + string(rune('0'+i))},
		})
	}

	var buf bytes.Buffer
	opts := DefaultPrintOptions()
	opts.Writer = &buf
	opts.Limit = 2

	if err := h.Print(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %s", len(lines), buf.String())
	}
}

func TestPrint_SinceFilter(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	h := &History{
		Entries: []Entry{
			{Timestamp: base, Opened: []string{"tcp/80"}},
			{Timestamp: base.Add(2 * time.Hour), Opened: []string{"tcp/443"}},
			{Timestamp: base.Add(4 * time.Hour), Opened: []string{"tcp/8080"}},
		},
	}

	var buf bytes.Buffer
	opts := DefaultPrintOptions()
	opts.Writer = &buf
	opts.Since = base.Add(1 * time.Hour)

	if err := h.Print(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "tcp/80") && !strings.Contains(out, "tcp/443") {
		t.Errorf("since filter not applied correctly: %s", out)
	}
	if !strings.Contains(out, "tcp/443") || !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected newer entries in output: %s", out)
	}
}
