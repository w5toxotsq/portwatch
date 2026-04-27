package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// PrintOptions controls how history entries are rendered.
type PrintOptions struct {
	Writer io.Writer
	Limit  int
	Since  time.Time
}

// DefaultPrintOptions returns PrintOptions with sensible defaults.
func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		Writer: os.Stdout,
		Limit:  0, // 0 means no limit
	}
}

// Print writes formatted history entries to the configured writer.
func (h *History) Print(opts PrintOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	entries := h.Entries

	// Apply since filter
	if !opts.Since.IsZero() {
		filtered := make([]Entry, 0, len(entries))
		for _, e := range entries {
			if e.Timestamp.After(opts.Since) || e.Timestamp.Equal(opts.Since) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No history entries found.")
		return err
	}

	// Apply limit (take last N entries)
	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[len(entries)-opts.Limit:]
	}

	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		if len(e.Opened) > 0 {
			fmt.Fprintf(w, "[%s] OPENED: %s\n", ts, join(e.Opened))
		}
		if len(e.Closed) > 0 {
			fmt.Fprintf(w, "[%s] CLOSED: %s\n", ts, join(e.Closed))
		}
	}

	return nil
}

// join formats a slice of port strings into a human-readable list.
func join(ports []string) string {
	return strings.Join(ports, ", ")
}
