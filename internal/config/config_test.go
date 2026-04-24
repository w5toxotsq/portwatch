package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeConfigFile(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeConfigFile(t, map[string]any{
		"ports":         []int{80, 443, 8080},
		"protocol":      "tcp",
		"interval":      int(10 * time.Second),
		"snapshot_path": "/tmp/snap.json",
		"alert_on_start": true,
	})

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Protocol != "tcp" {
		t.Errorf("protocol: got %q, want \"tcp\"", cfg.Protocol)
	}
	if len(cfg.Ports) != 3 {
		t.Errorf("ports: got %d, want 3", len(cfg.Ports))
	}
	if !cfg.AlertOnStart {
		t.Error("alert_on_start: got false, want true")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("{invalid json")
	f.Close()

	_, err := config.Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoad_InvalidProtocol(t *testing.T) {
	path := writeConfigFile(t, map[string]any{
		"protocol": "ftp",
		"interval": int(5 * time.Second),
	})

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid protocol, got nil")
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Protocol != "tcp" {
		t.Errorf("default protocol: got %q, want \"tcp\"", cfg.Protocol)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("default interval: got %v, want 30s", cfg.Interval)
	}
	if cfg.SnapshotPath != config.DefaultSnapshotPath {
		t.Errorf("default snapshot path: got %q, want %q", cfg.SnapshotPath, config.DefaultSnapshotPath)
	}
}
