// Package ratelimit provides a simple key-based rate limiter used to
// suppress repeated alert notifications within a configurable time window.
//
// The Limiter is safe for concurrent use. A typical key is a canonical
// string representation of a change set (e.g. "opened:tcp:8080") so that
// the same alert is not dispatched more than once per window even if the
// daemon polls frequently.
//
// Example:
//
//	l := ratelimit.New(5 * time.Minute)
//	if l.Allow(key) {
//		notifier.Notify(changes)
//	}
package ratelimit
