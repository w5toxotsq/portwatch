package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
)

func runReport(args []string) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	formatFlag := fs.String("format", "text", "Output format: text or json")
	limitFlag := fs.Int("limit", 50, "Maximum number of entries to display")
	sinceFlag := fs.String("since", "", "Show entries since duration ago (e.g. 24h, 7d)")
	fileFlag := fs.String("file", "portwatch-history.json", "Path to history file")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("report: failed to parse flags: %v", err)
	}

	h, err := history.Load(*fileFlag)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("report: failed to load history: %v", err)
	}
	if h == nil {
		h = history.New()
	}

	opts := report.DefaultOptions()
	opts.Writer = os.Stdout

	switch *formatFlag {
	case "json":
		opts.Format = report.FormatJSON
	case "text":
		opts.Format = report.FormatText
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q, defaulting to text\n", *formatFlag)
	}

	opts.Limit = *limitFlag

	if *sinceFlag != "" {
		d, err := time.ParseDuration(*sinceFlag)
		if err != nil {
			log.Fatalf("report: invalid --since value %q: %v", *sinceFlag, err)
		}
		opts.Since = time.Now().Add(-d)
	}

	if err := report.Generate(h, opts); err != nil {
		log.Fatalf("report: failed to generate report: %v", err)
	}
}
