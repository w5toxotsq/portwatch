// Package snapshot provides functionality for capturing and comparing
// port scan results over time to detect unexpected changes.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Diff describes the changes between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// New creates a new Snapshot with the current timestamp.
func New(ports []int) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now(),
		Ports:     ports,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func Save(s *Snapshot, path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Compare returns a Diff between a previous and current snapshot.
// Opened contains ports present in current but not in previous.
// Closed contains ports present in previous but not in current.
func Compare(previous, current *Snapshot) Diff {
	prev := toSet(previous.Ports)
	curr := toSet(current.Ports)

	var diff Diff
	for p := range curr {
		if !prev[p] {
			diff.Opened = append(diff.Opened, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			diff.Closed = append(diff.Closed, p)
		}
	}
	return diff
}

// HasChanges reports whether a Diff contains any port changes.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
