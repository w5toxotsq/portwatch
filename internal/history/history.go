package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single recorded change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []string  `json:"opened,omitempty"`
	Closed    []string  `json:"closed,omitempty"`
}

// History holds an ordered list of change entries.
type History struct {
	Entries []Entry `json:"entries"`
	path    string
}

// New returns a History backed by the given file path.
func New(path string) *History {
	return &History{path: path}
}

// Record appends a new entry with the current timestamp.
func (h *History) Record(opened, closed []string) {
	if len(opened) == 0 && len(closed) == 0 {
		return
	}
	h.Entries = append(h.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Opened:    opened,
		Closed:    closed,
	})
}

// Save persists the history to disk as JSON.
func (h *History) Save() error {
	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	f, err := os.Create(h.path)
	if err != nil {
		return fmt.Errorf("history: create: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(h); err != nil {
		return fmt.Errorf("history: encode: %w", err)
	}
	return nil
}

// Load reads history from disk. If the file does not exist, an empty History is returned.
func Load(path string) (*History, error) {
	h := New(path)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: open: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(h); err != nil {
		return nil, fmt.Errorf("history: decode: %w", err)
	}
	return h, nil
}
