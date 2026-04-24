package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNew(t *testing.T) {
	ports := []int{80, 443, 8080}
	s := snapshot.New(ports)
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(s.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(s.Ports))
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New([]int{22, 80, 443})
	if err := snapshot.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: want %d, got %d", len(orig.Ports), len(loaded.Ports))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestCompare_Opened(t *testing.T) {
	prev := snapshot.New([]int{80, 443})
	curr := snapshot.New([]int{80, 443, 8080})

	diff := snapshot.Compare(prev, curr)
	if len(diff.Opened) != 1 || diff.Opened[0] != 8080 {
		t.Errorf("expected opened=[8080], got %v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diff.Closed)
	}
}

func TestCompare_Closed(t *testing.T) {
	prev := snapshot.New([]int{80, 443, 22})
	curr := snapshot.New([]int{80, 443})

	diff := snapshot.Compare(prev, curr)
	if len(diff.Closed) != 1 || diff.Closed[0] != 22 {
		t.Errorf("expected closed=[22], got %v", diff.Closed)
	}
	if len(diff.Opened) != 0 {
		t.Errorf("expected no opened ports, got %v", diff.Opened)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	prev := snapshot.New([]int{80, 443})
	curr := snapshot.New([]int{80, 443})

	diff := snapshot.Compare(prev, curr)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}
