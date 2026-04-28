// Package watcher implements the port-polling loop for portwatch.
//
// A Watcher periodically scans the configured set of ports, compares the
// results against the most recent snapshot, records any changes to the
// history log, and dispatches alerts through the configured notifier chain.
//
// Typical usage:
//
//	w := watcher.New(cfg, alerter, hist)
//	if err := w.Run(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// Run blocks until the provided context is cancelled, making it easy to
// integrate with signal-based shutdown in the daemon layer.
package watcher
