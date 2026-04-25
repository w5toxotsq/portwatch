package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// PrintOptions controls output formatting.
type PrintOptions struct {
	Out    io.Writer
	Limit  int
	Since  time.Time
}

// DefaultPrintOptions returns sensible defaults writing to stdout.
func DefaultPrintOptions() PrintOptions {
	return PrintOptions{
		Out:   os.Stdout,
		Limit: 50,
	}
}

// Print writes a human-readable table of history entries to opts.Out.
func (h *History) Print(opts PrintOptions) {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	entries := h.Entries
	if !opts.Since.IsZero() {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Timestamp.After(opts.Since) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}
	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[len(entries)-opts.Limit:]
	}

	if len(entries) == 0 {
		fmt.Fprintln(opts.Out, "no history entries found")
		return
	}

	w := tabwriter.NewWriter(opts.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tOPENED\tCLOSED")
	for _, e := range entries {
		ts := e.Timestamp.Local().Format(time.RFC3339)
		opened := join(e.Opened)
		closed := join(e.Closed)
		fmt.Fprintf(w, "%s\t%s\t%s\n", ts, opened, closed)
	}
	w.Flush()
}

func join(ports []string) string {
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}
