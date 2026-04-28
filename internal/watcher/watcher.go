package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher polls open ports at a regular interval and triggers alerts on changes.
type Watcher struct {
	cfg     *config.Config
	alerter *alert.Alert
	hist    *history.History
}

// New creates a new Watcher with the given configuration, alerter, and history.
func New(cfg *config.Config, alerter *alert.Alert, hist *history.History) *Watcher {
	return &Watcher{
		cfg:     cfg,
		alerter: alerter,
		hist:    hist,
	}
}

// Run starts the polling loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := w.poll(); err != nil {
				return fmt.Errorf("watcher poll: %w", err)
			}
		}
	}
}

// poll performs a single scan cycle: scan ports, compare to last snapshot,
// record history, and fire alerts if anything changed.
func (w *Watcher) poll() error {
	ports, err := scanner.OpenPorts(w.cfg.Ports, w.cfg.Protocol)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	current := snapshot.New(ports)

	prev, err := snapshot.Load(w.cfg.SnapshotPath)
	if err != nil {
		// No previous snapshot yet; save current and return.
		return snapshot.Save(current, w.cfg.SnapshotPath)
	}

	changes := snapshot.Compare(prev, current)

	if err := w.hist.Record(changes); err != nil {
		return fmt.Errorf("history record: %w", err)
	}

	if err := w.alerter.Notify(changes); err != nil {
		return fmt.Errorf("alert notify: %w", err)
	}

	return snapshot.Save(current, w.cfg.SnapshotPath)
}
