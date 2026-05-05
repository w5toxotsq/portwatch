package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// runBaseline implements the "baseline" sub-command.  It scans the configured
// ports, saves the result as the new baseline, and prints a summary.
func runBaseline(cfg *config.Config, format string) error {
	ports, err := scanner.OpenPorts(cfg.Ports, cfg.Protocol, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	b := baseline.New(ports)
	if err := baseline.Save(b, cfg.BaselinePath); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}

	switch format {
	case "json":
		return printBaselineJSON(b)
	default:
		return printBaselineText(b)
	}
}

func printBaselineText(b *baseline.Baseline) error {
	fmt.Fprintf(os.Stdout, "Baseline saved at %s\n", b.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	if len(b.Ports) == 0 {
		fmt.Fprintln(os.Stdout, "  (no open ports recorded)")
		return nil
	}
	for _, p := range b.Ports {
		fmt.Fprintf(os.Stdout, "  %s\n", p)
	}
	return nil
}

func printBaselineJSON(b *baseline.Baseline) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}
