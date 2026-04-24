package notifier

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogNotifier writes port change events to a writer in a structured log format.
type LogNotifier struct {
	out io.Writer
}

// NewLogNotifier creates a LogNotifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogNotifier{out: w}
}

// Notify writes a formatted log line for the given event.
func (l *LogNotifier) Notify(event Event) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := fmt.Fprintf(
		l.out,
		"%s [portwatch] port %s %s/%d\n",
		timestamp,
		string(event.Type),
		event.Protocol,
		event.Port,
	)
	return err
}
