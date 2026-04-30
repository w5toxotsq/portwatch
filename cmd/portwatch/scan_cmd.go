package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

type scanOptions struct {
	configPath string
	outputJSON bool
}

func runScan(opts scanOptions) error {
	cfg, err := config.Load(opts.configPath)
	if err != nil {
		cfg = config.Default()
	}

	ports, err := scanner.OpenPorts(cfg)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if opts.outputJSON {
		return printScanJSON(ports)
	}

	return printScanText(ports)
}

func printScanText(ports []scanner.Port) error {
	if len(ports) == 0 {
		fmt.Println("No open ports found.")
		return nil
	}
	fmt.Printf("%-10s %-8s %s\n", "PORT", "PROTO", "ADDRESS")
	fmt.Println("-----------------------------")
	for _, p := range ports {
		fmt.Printf("%-10d %-8s %s\n", p.Port, p.Protocol, p.Address)
	}
	return nil
}

func printScanJSON(ports []scanner.Port) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(ports)
}
