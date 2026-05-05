// main is the entry point for the portwatch CLI.
// It wires together all subcommands and global flags.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfigPath = "/etc/portwatch/config.json"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// newRootCmd builds the root cobra command with all subcommands attached.
func newRootCmd() *cobra.Command {
	var configPath string

	root := &cobra.Command{
		Use:   "portwatch",
		Short: "Monitor open ports and alert on unexpected changes",
		Long: `portwatch is a lightweight CLI daemon that monitors open TCP/UDP ports
and alerts you when ports are opened or closed unexpectedly.

Use 'portwatch watch' to start the daemon, or 'portwatch scan' for a
one-shot port scan.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVarP(
		&configPath, "config", "c", defaultConfigPath,
		"path to the portwatch config file",
	)

	// scan: one-shot port scan
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Perform a one-shot port scan and print results",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			return runScan(configPath, format, os.Stdout)
		},
	}
	scanCmd.Flags().String("format", "text", "output format: text or json")

	// baseline: capture or compare a port baseline
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Capture a port baseline or compare against the current state",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			return runBaseline(configPath, format, os.Stdout)
		},
	}
	baselineCmd.Flags().String("format", "text", "output format: text or json")

	// watch: start the monitoring daemon
	watchCmd := newWatchCmd(&configPath)

	// history: display past change events
	historyCmd := newHistoryCmd(&configPath)

	// report: generate a summary report
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a summary report from recorded history",
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			limit, _ := cmd.Flags().GetInt("limit")
			return runReport(configPath, format, limit, os.Stdout)
		},
	}
	reportCmd.Flags().String("format", "text", "output format: text or json")
	reportCmd.Flags().Int("limit", 0, "maximum number of history entries to include (0 = all)")

	root.AddCommand(scanCmd, baselineCmd, watchCmd, historyCmd, reportCmd)

	return root
}
