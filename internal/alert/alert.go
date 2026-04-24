// Package alert provides functionality for notifying users of port changes
// detected by the portwatch daemon.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends alerts to a configured output.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes alerts based on the diff between two snapshots.
func (n *Notifier) Notify(diff snapshot.Diff) []Alert {
	var alerts []Alert

	for _, p := range diff.Opened {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port opened: %s/%d", p.Protocol, p.Port),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	for _, p := range diff.Closed {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port closed: %s/%d", p.Protocol, p.Port),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	if len(alerts) == 0 {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelInfo,
			Message:   "no port changes detected",
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	return alerts
}
