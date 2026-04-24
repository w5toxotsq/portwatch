package notifier_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
)

// fakeNotifier records events and optionally returns an error.
type fakeNotifier struct {
	events []notifier.Event
	errOn  bool
}

func (f *fakeNotifier) Notify(e notifier.Event) error {
	f.events = append(f.events, e)
	if f.errOn {
		return errors.New("notify error")
	}
	return nil
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	ln := notifier.NewLogNotifier(&buf)

	evt := notifier.Event{Type: notifier.EventOpened, Protocol: "tcp", Port: 8080}
	if err := ln.Notify(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "opened") {
		t.Errorf("expected 'opened' in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected 'tcp/8080' in output, got: %s", out)
	}
}

func TestLogNotifier_DefaultsToStdout(t *testing.T) {
	ln := notifier.NewLogNotifier(nil)
	if ln == nil {
		t.Fatal("expected non-nil LogNotifier")
	}
}

func TestMulti_Notify_DispatchesAll(t *testing.T) {
	a := &fakeNotifier{}
	b := &fakeNotifier{}
	m := notifier.NewMulti(a, b)

	evt := notifier.Event{Type: notifier.EventClosed, Protocol: "udp", Port: 53}
	if err := m.Notify(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(a.events) != 1 || len(b.events) != 1 {
		t.Errorf("expected each notifier to receive 1 event")
	}
}

func TestMulti_Notify_ReturnsFirstError(t *testing.T) {
	a := &fakeNotifier{errOn: true}
	b := &fakeNotifier{}
	m := notifier.NewMulti(a, b)

	evt := notifier.Event{Type: notifier.EventOpened, Protocol: "tcp", Port: 443}
	err := m.Notify(evt)
	if err == nil {
		t.Fatal("expected an error from multi notifier")
	}
	// b should still have received the event
	if len(b.events) != 1 {
		t.Errorf("expected second notifier to still receive event")
	}
}
