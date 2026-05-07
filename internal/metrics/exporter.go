package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// ExportFormat controls the output format of exported metrics.
type ExportFormat string

const (
	ExportText ExportFormat = "text"
	ExportJSON ExportFormat = "json"
)

// ExportOptions configures how metrics are exported.
type ExportOptions struct {
	Format ExportFormat
	Writer io.Writer
}

// Snapshot is a point-in-time copy of collected metrics.
type Snapshot struct {
	UptimeSince  time.Time `json:"uptime_since"`
	UptimeHuman  string    `json:"uptime_human"`
	TotalPolls   int64     `json:"total_polls"`
	FailedPolls  int64     `json:"failed_polls"`
	AlertsSent   int64     `json:"alerts_sent"`
	CollectedAt  time.Time `json:"collected_at"`
}

// Export writes the current metrics to w using the specified format.
func Export(m *Metrics, opts ExportOptions) error {
	snap := buildSnapshot(m)
	switch opts.Format {
	case ExportJSON:
		return exportJSON(snap, opts.Writer)
	default:
		return exportText(snap, opts.Writer)
	}
}

func buildSnapshot(m *Metrics) Snapshot {
	g := m.Get()
	uptime := time.Since(g.UptimeSince).Round(time.Second)
	return Snapshot{
		UptimeSince: g.UptimeSince,
		UptimeHuman: formatDuration(uptime),
		TotalPolls:  g.TotalPolls,
		FailedPolls: g.FailedPolls,
		AlertsSent:  g.AlertsSent,
		CollectedAt: time.Now().UTC(),
	}
}

func exportJSON(snap Snapshot, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

func exportText(snap Snapshot, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Metric\tValue\n")
	fmt.Fprintf(tw, "------\t-----\n")
	fmt.Fprintf(tw, "Uptime\t%s\n", snap.UptimeHuman)
	fmt.Fprintf(tw, "Total Polls\t%d\n", snap.TotalPolls)
	fmt.Fprintf(tw, "Failed Polls\t%d\n", snap.FailedPolls)
	fmt.Fprintf(tw, "Alerts Sent\t%d\n", snap.AlertsSent)
	fmt.Fprintf(tw, "Collected At\t%s\n", snap.CollectedAt.Format(time.RFC3339))
	return tw.Flush()
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
