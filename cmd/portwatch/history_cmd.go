package main

import (
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/spf13/cobra"
)

func newHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show port change history",
		Long:  "Display a log of detected port changes recorded during daemon operation.",
		RunE:  runHistory,
	}

	cmd.Flags().StringP("file", "f", "portwatch-history.json", "Path to history file")
	cmd.Flags().IntP("limit", "n", 20, "Maximum number of entries to show")
	cmd.Flags().StringP("since", "s", "", "Show entries since timestamp (RFC3339, e.g. 2024-01-01T00:00:00Z)")
	cmd.Flags().StringP("format", "o", "text", "Output format: text or json")

	return cmd
}

func runHistory(cmd *cobra.Command, args []string) error {
	filePath, _ := cmd.Flags().GetString("file")
	limit, _ := cmd.Flags().GetInt("limit")
	sinceStr, _ := cmd.Flags().GetString("since")
	format, _ := cmd.Flags().GetString("format")

	h, err := history.Load(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "no history file found:", filePath)
			return nil
		}
		return fmt.Errorf("loading history: %w", err)
	}

	opts := history.DefaultPrintOptions()
	opts.Limit = limit
	opts.Format = format

	if sinceStr != "" {
		t, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return fmt.Errorf("invalid --since value %q: expected RFC3339 format", sinceStr)
		}
		opts.Since = t
	}

	return history.Print(os.Stdout, h, opts)
}
