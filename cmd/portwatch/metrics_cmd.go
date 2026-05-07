package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portwatch/internal/metrics"
)

func newMetricsCmd() *cobra.Command {
	var (
		format  string
		serveAt string
	)

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Display or serve runtime metrics",
		Long: `Display collected runtime metrics such as poll counts, failures,
and alerts sent. Use --serve to expose a live JSON endpoint instead.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serveAt != "" {
				return runMetricsServe(serveAt)
			}
			return runMetricsPrint(format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	cmd.Flags().StringVar(&serveAt, "serve", "", "Address to serve metrics HTTP endpoint (e.g. :9090)")
	return cmd
}

func runMetricsPrint(format string) error {
	// In a real daemon the Metrics instance would be shared; here we
	// construct a zero-value instance so the command is always safe to run.
	m := metrics.New()

	fmt := metrics.ExportText
	if format == "json" {
		fmt = metrics.ExportJSON
	}

	return metrics.Export(m, metrics.ExportOptions{
		Format: fmt,
		Writer: os.Stdout,
	})
}

func runMetricsServe(addr string) error {
	m := metrics.New()
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler(m))

	fmt.Fprintf(os.Stderr, "serving metrics on %s/metrics\n", addr)
	return http.ListenAndServe(addr, mux)
}
