package baseline_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func TestNew_SetsPortsAndTime(t *testing.T) {
	ports := []string{"tcp:80", "tcp:443"}
	b := baseline.New(ports)
	if len(b.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(b.Ports))
	}
	if b.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	orig := baseline.New([]string{"tcp:22", "tcp:80"})
	if err := baseline.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("expected %d ports, got %d", len(orig.Ports), len(loaded.Ports))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err != baseline.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_UnexpectedAndMissing(t *testing.T) {
	b := baseline.New([]string{"tcp:22", "tcp:80"})
	current := []string{"tcp:22", "tcp:443"}

	unexpected, missing := baseline.Compare(b, current)

	if len(unexpected) != 1 || unexpected[0] != "tcp:443" {
		t.Errorf("unexpected ports: got %v, want [tcp:443]", unexpected)
	}
	if len(missing) != 1 || missing[0] != "tcp:80" {
		t.Errorf("missing ports: got %v, want [tcp:80]", missing)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	b := baseline.New([]string{"tcp:22", "tcp:80"})
	unexpected, missing := baseline.Compare(b, []string{"tcp:22", "tcp:80"})
	if len(unexpected) != 0 || len(missing) != 0 {
		t.Errorf("expected no changes, got unexpected=%v missing=%v", unexpected, missing)
	}
}

func TestSave_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	b := baseline.New([]string{"tcp:8080"})
	if err := baseline.Save(b, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data, _ := os.ReadFile(path)
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
}
