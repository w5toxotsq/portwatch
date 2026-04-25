package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/history"
)

func TestRecord_AddsEntry(t *testing.T) {
	h := history.New("/tmp/unused.json")
	h.Record([]string{"tcp:8080"}, nil)
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	if h.Entries[0].Opened[0] != "tcp:8080" {
		t.Errorf("unexpected opened port: %s", h.Entries[0].Opened[0])
	}
}

func TestRecord_SkipsEmptyChanges(t *testing.T) {
	h := history.New("/tmp/unused.json")
	h.Record(nil, nil)
	if len(h.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(h.Entries))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := history.New(path)
	h.Record([]string{"tcp:9090"}, []string{"tcp:22"})

	if err := h.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := history.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Closed[0] != "tcp:22" {
		t.Errorf("unexpected closed port: %s", loaded.Entries[0].Closed[0])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	h, err := history.Load("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := history.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "history.json")
	h := history.New(path)
	h.Record([]string{"tcp:443"}, nil)
	if err := h.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON written: %v", err)
	}
}
