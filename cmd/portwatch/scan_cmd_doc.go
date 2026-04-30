// Package main provides the portwatch CLI.
//
// The scan subcommand performs a one-shot scan of the ports defined in the
// active configuration and prints the results to stdout. It does not persist
// a snapshot or compare against a previous baseline — use the daemon for
// continuous monitoring.
//
// Usage:
//
//	portwatch scan [--config <path>] [--json]
//
// Flags:
//
//	--config  Path to the configuration file (default: portwatch.json)
//	--json    Emit results as a JSON array instead of a human-readable table
package main
