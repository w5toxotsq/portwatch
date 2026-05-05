package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/watcher"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Start the port watcher daemon",
		Long:  "Continuously monitors open ports and alerts on unexpected changes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(configPath)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file (optional)")
	return cmd
}

func runWatch(configPath string) error {
	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.Load(configPath)
	} else {
		cfg = config.Default()
	}
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	hist, err := history.Load(cfg.HistoryPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading history: %w", err)
	}
	if hist == nil {
		hist = history.New(cfg.HistoryPath)
	}

	logNotifier := notifier.NewLogNotifier(os.Stdout)
	multi := notifier.NewMulti(logNotifier)

	w := watcher.New(cfg, hist, multi)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintf(os.Stderr, "portwatch: watching ports every %s\n", cfg.Interval)
	w.Run(ctx)
	fmt.Fprintln(os.Stderr, "portwatch: stopped")
	return nil
}
