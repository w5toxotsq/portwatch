package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/history"
)

// Format defines the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures report generation.
type Options struct {
	Format  Format
	Limit   int
	Since   time.Time
	Writer  io.Writer
}

// DefaultOptions returns sensible report defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatText,
		Limit:  50,
		Writer: os.Stdout,
	}
}

// Generate writes a formatted report of history entries to the configured writer.
func Generate(h *history.History, opts Options) error {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	entries := h.Entries()
	if !opts.Since.IsZero() {
		filtered := entries[:0]
		for _, e := range entries {
			if !e.Timestamp.Before(opts.Since) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}
	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[len(entries)-opts.Limit:]
	}

	switch opts.Format {
	case FormatJSON:
		return generateJSON(entries, opts.Writer)
	default:
		return generateText(entries, opts.Writer)
	}
}

func generateText(entries []history.Entry, w io.Writer) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No history entries found.")
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tOPENED\tCLOSED")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%d\n",
			e.Timestamp.Format(time.RFC3339),
			len(e.Changes.Opened),
			len(e.Changes.Closed),
		)
	}
	return tw.Flush()
}

func generateJSON(entries []history.Entry, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
