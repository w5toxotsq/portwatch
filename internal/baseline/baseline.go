package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a known-good set of open ports that the user has
// explicitly approved. Changes are measured against the baseline rather
// than the previous snapshot when one is present.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Ports     []string  `json:"ports"`
}

// New creates a Baseline from the given list of ports, stamped with the
// current UTC time.
func New(ports []string) *Baseline {
	return &Baseline{
		CreatedAt: time.Now().UTC(),
		Ports:     ports,
	}
}

// Save writes the baseline to path as JSON.
func Save(b *Baseline, path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline from path. Returns ErrNotFound when the file does
// not exist so callers can distinguish a missing baseline from other errors.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Compare returns the ports that are in current but not in the baseline
// (unexpected) and the ports that are in the baseline but not in current
// (missing).
func Compare(b *Baseline, current []string) (unexpected, missing []string) {
	baseSet := make(map[string]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		baseSet[p] = struct{}{}
	}
	currentSet := make(map[string]struct{}, len(current))
	for _, p := range current {
		currentSet[p] = struct{}{}
	}
	for _, p := range current {
		if _, ok := baseSet[p]; !ok {
			unexpected = append(unexpected, p)
		}
	}
	for _, p := range b.Ports {
		if _, ok := currentSet[p]; !ok {
			missing = append(missing, p)
		}
	}
	return unexpected, missing
}

// ErrNotFound is returned by Load when no baseline file exists at the given path.
var ErrNotFound = fmt.Errorf("baseline: file not found")
