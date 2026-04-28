package watcher_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/watcher"
)

func defaultConfig(t *testing.T) *config.Config {
	t.Helper()
	dir := t.TempDir()
	cfg := config.Default()
	cfg.SnapshotPath = filepath.Join(dir, "snap.json")
	cfg.Interval = 50 * time.Millisecond
	cfg.Ports = []int{}
	return cfg
}

func TestWatcher_CreatesSnapshotOnFirstPoll(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(os.Stdout)
	hist, _ := history.New(filepath.Join(t.TempDir(), "hist.json"))

	w := watcher.New(cfg, alerter, hist)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if _, err := os.Stat(cfg.SnapshotPath); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to be created")
	}
}

func TestWatcher_DetectsChanges(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(os.Stdout)
	hist, _ := history.New(filepath.Join(t.TempDir(), "hist.json"))

	// Pre-seed snapshot with a port that is no longer open.
	initial := snapshot.New([]int{9999})
	if err := snapshot.Save(initial, cfg.SnapshotPath); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	w := watcher.New(cfg, alerter, hist)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	// History should have recorded the closure of port 9999.
	loaded, err := history.Load(hist.Path())
	if err != nil {
		t.Fatalf("load history: %v", err)
	}
	if len(loaded.Entries) == 0 {
		t.Fatal("expected at least one history entry")
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(os.Stdout)
	hist, _ := history.New(filepath.Join(t.TempDir(), "hist.json"))

	w := watcher.New(cfg, alerter, hist)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancellation")
	}
}

// Ensure snapshot written after a poll is valid JSON.
func TestWatcher_SnapshotIsValidJSON(t *testing.T) {
	cfg := defaultConfig(t)
	alerter := alert.New(os.Stdout)
	hist, _ := history.New(filepath.Join(t.TempDir(), "hist.json"))

	w := watcher.New(cfg, alerter, hist)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)

	data, err := os.ReadFile(cfg.SnapshotPath)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("snapshot is not valid JSON: %v", err)
	}
}
