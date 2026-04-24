// Package notifier provides interfaces and implementations for dispatching
// port change events to one or more notification backends.
//
// The core Notifier interface allows different backends (logging, webhooks,
// desktop alerts, etc.) to be used interchangeably. The Multi type enables
// fan-out delivery to several backends simultaneously.
//
// Example usage:
//
//	ln := notifier.NewLogNotifier(os.Stdout)
//	multi := notifier.NewMulti(ln)
//	multi.Notify(notifier.Event{
//		Type:     notifier.EventOpened,
//		Protocol: "tcp",
//		Port:     8080,
//	})
package notifier
