// Package daemon provides the core polling loop for portwatch.
// It periodically scans open ports, compares against the last snapshot,
// and triggers alerts when changes are detected.
package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon holds the runtime state for the polling loop.
type Daemon struct {
	cfg     *config.Config
	alerter *alert.Alerter
}

// New creates a new Daemon with the provided configuration and alerter.
func New(cfg *config.Config, alerter *alert.Alerter) *Daemon {
	return &Daemon{
		cfg:     cfg,
		alerter: alerter,
	}
}

// Run starts the polling loop. It blocks until the provided stop channel is closed.
func (d *Daemon) Run(stop <-chan struct{}) error {
	ticker := time.NewTicker(time.Duration(d.cfg.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("portwatch daemon started (interval: %ds, ports: %v, protocol: %s)",
		d.cfg.IntervalSeconds, d.cfg.Ports, d.cfg.Protocol)

	// Run an immediate scan on startup before waiting for the first tick.
	if err := d.poll(); err != nil {
		log.Printf("warn: initial poll failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := d.poll(); err != nil {
				log.Printf("warn: poll failed: %v", err)
			}
		case <-stop:
			log.Println("portwatch daemon stopped")
			return nil
		}
	}
}

// poll performs a single scan-compare-alert cycle.
func (d *Daemon) poll() error {
	current := scanner.OpenPorts(d.cfg.Ports, d.cfg.Protocol)

	prev, err := snapshot.Load(d.cfg.SnapshotPath)
	if err != nil {
		// No previous snapshot yet — save current state and return.
		next := snapshot.New(current)
		return snapshot.Save(next, d.cfg.SnapshotPath)
	}

	changes := snapshot.Compare(prev, current)
	if err := d.alerter.Notify(changes); err != nil {
		log.Printf("warn: alert notify failed: %v", err)
	}

	next := snapshot.New(current)
	return snapshot.Save(next, d.cfg.SnapshotPath)
}
