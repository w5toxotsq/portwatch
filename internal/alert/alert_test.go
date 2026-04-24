package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/snapshot"
)

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := snapshot.Diff{
		Opened: []snapshot.PortEntry{{Protocol: "tcp", Port: 8080}},
	}

	alerts := n.Notify(diff)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelAlert {
		t.Errorf("expected level ALERT, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "port opened: tcp/8080") {
		t.Errorf("expected output to contain port info, got: %s", buf.String())
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := snapshot.Diff{
		Closed: []snapshot.PortEntry{{Protocol: "udp", Port: 53}},
	}

	alerts := n.Notify(diff)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelWarn {
		t.Errorf("expected level WARN, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "port closed: udp/53") {
		t.Errorf("expected output to contain port info, got: %s", buf.String())
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	alerts := n.Notify(snapshot.Diff{})

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelInfo {
		t.Errorf("expected level INFO, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "no port changes detected") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Ensure New(nil) does not panic
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
