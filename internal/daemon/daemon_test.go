package daemon_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func defaultConfig(t *testing.T) *config.Config {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), "snapshot-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()
	os.Remove(tmp.Name()) // let snapshot create it fresh

	cfg := config.Default()
	cfg.SnapshotPath = tmp.Name()
	cfg.IntervalSeconds = 1
	cfg.Ports = []int{} // no real ports to avoid flakiness
	return cfg
}

func TestDaemon_StopsOnSignal(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(nil)
	d := daemon.New(cfg, alerter)

	stop := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		done <- d.Run(stop)
	}()

	// Allow at least one poll cycle.
	time.Sleep(150 * time.Millisecond)
	close(stop)

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not stop within timeout")
	}
}

func TestDaemon_CreatesSnapshotOnFirstPoll(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(nil)
	d := daemon.New(cfg, alerter)

	stop := make(chan struct{})
	go func() { d.Run(stop) }() //nolint:errcheck

	time.Sleep(150 * time.Millisecond)
	close(stop)
	time.Sleep(50 * time.Millisecond)

	if _, err := os.Stat(cfg.SnapshotPath); os.IsNotExist(err) {
		t.Error("expected snapshot file to be created after first poll")
	}
}

func TestNew_ReturnsDaemon(t *testing.T) {
	cfg := config.Default()
	alerter := alert.New(nil)
	d := daemon.New(cfg, alerter)
	if d == nil {
		t.Fatal("expected non-nil Daemon")
	}
}
