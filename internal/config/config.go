// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DefaultSnapshotPath is used when no snapshot path is specified.
const DefaultSnapshotPath = "/tmp/portwatch_snapshot.json"

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Ports is the list of ports to monitor. If empty, all open ports are monitored.
	Ports []int `json:"ports"`

	// Protocol specifies the protocol to scan: "tcp" or "udp".
	Protocol string `json:"protocol"`

	// Interval is how often portwatch scans for changes.
	Interval time.Duration `json:"interval"`

	// SnapshotPath is the file path used to persist port snapshots.
	SnapshotPath string `json:"snapshot_path"`

	// AlertOnStart controls whether an alert fires on the initial scan.
	AlertOnStart bool `json:"alert_on_start"`
}

// Load reads a JSON config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Protocol:     "tcp",
		Interval:     30 * time.Second,
		SnapshotPath: DefaultSnapshotPath,
		AlertOnStart: false,
	}
}

// validate checks that required fields contain acceptable values.
func (c *Config) validate() error {
	if c.Protocol == "" {
		c.Protocol = "tcp"
	}
	if c.Protocol != "tcp" && c.Protocol != "udp" {
		return fmt.Errorf("config: invalid protocol %q: must be \"tcp\" or \"udp\"", c.Protocol)
	}
	if c.Interval <= 0 {
		c.Interval = 30 * time.Second
	}
	if c.SnapshotPath == "" {
		c.SnapshotPath = DefaultSnapshotPath
	}
	return nil
}
