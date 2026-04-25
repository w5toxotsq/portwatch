package history_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestPrint_ShowsEntries(t *testing.T) {
	h := history.New("/tmp/unused.json")
	h.Record([]string{"tcp:8080"}, []string{"tcp:22"})

	var buf bytes.Buffer
	opts := history.DefaultPrintOptions()
	opts.Out = &buf
	h.Print(opts)

	out := buf.String()
	if !strings.Contains(out, "tcp:8080") {
		t.Errorf("expected opened port in output, got:\n%s", out)
	}
	if !strings.Contains(out, "tcp:22") {
		t.Errorf("expected closed port in output, got:\n%s", out)
	}
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestPrint_EmptyHistory(t *testing.T) {
	h := history.New("/tmp/unused.json")
	var buf bytes.Buffer
	opts := history.DefaultPrintOptions()
	opts.Out = &buf
	h.Print(opts)

	if !strings.Contains(buf.String(), "no history entries found") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrint_LimitEntries(t *testing.T) {
	h := history.New("/tmp/unused.json")
	for i := 0; i < 10; i++ {
		h.Record([]string{"tcp:8080"}, nil)
	}

	var buf bytes.Buffer
	opts := history.DefaultPrintOptions()
	opts.Out = &buf
	opts.Limit = 3
	h.Print(opts)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 3 data lines
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header+3), got %d:\n%s", len(lines), buf.String())
	}
}

func TestPrint_SinceFilter(t *testing.T) {
	h := history.New("/tmp/unused.json")
	h.Record([]string{"tcp:22"}, nil)  // old entry via direct append below
	h.Entries[0].Timestamp = time.Now().Add(-2 * time.Hour)
	h.Record([]string{"tcp:443"}, nil) // recent

	var buf bytes.Buffer
	opts := history.DefaultPrintOptions()
	opts.Out = &buf
	opts.Since = time.Now().Add(-30 * time.Minute)
	h.Print(opts)

	out := buf.String()
	if strings.Contains(out, "tcp:22") {
		t.Errorf("old entry should be filtered out")
	}
	if !strings.Contains(out, "tcp:443") {
		t.Errorf("recent entry should appear")
	}
}
