package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time summary of watcher activity.
type Snapshot struct {
	PollCount      int64         `json:"poll_count"`
	LastPollAt     time.Time     `json:"last_poll_at"`
	LastPollDur    time.Duration `json:"last_poll_duration_ns"`
	OpenedTotal    int64         `json:"opened_total"`
	ClosedTotal    int64         `json:"closed_total"`
	AlertsSent     int64         `json:"alerts_sent"`
	UptimeSince    time.Time     `json:"uptime_since"`
}

// Collector accumulates runtime metrics for the daemon.
type Collector struct {
	mu          sync.RWMutex
	snapshot    Snapshot
}

// New returns a new Collector with UptimeSince set to now.
func New() *Collector {
	return &Collector{
		snapshot: Snapshot{
			UptimeSince: time.Now().UTC(),
		},
	}
}

// RecordPoll records a completed poll cycle.
func (c *Collector) RecordPoll(dur time.Duration, opened, closed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snapshot.PollCount++
	c.snapshot.LastPollAt = time.Now().UTC()
	c.snapshot.LastPollDur = dur
	c.snapshot.OpenedTotal += int64(opened)
	c.snapshot.ClosedTotal += int64(closed)
}

// RecordAlert increments the alerts-sent counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snapshot.AlertsSent++
}

// Get returns a copy of the current metrics snapshot.
func (c *Collector) Get() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snapshot
}
