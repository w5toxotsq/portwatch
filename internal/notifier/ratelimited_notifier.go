package notifier

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/snapshot"
)

// RateLimitedNotifier wraps a Notifier and suppresses duplicate
// notifications that occur within the limiter's window.
type RateLimitedNotifier struct {
	inner   Notifier
	limiter *ratelimit.Limiter
}

// NewRateLimited returns a RateLimitedNotifier that delegates to inner
// but drops notifications whose change key was already dispatched within
// the limiter's window.
func NewRateLimited(inner Notifier, limiter *ratelimit.Limiter) *RateLimitedNotifier {
	return &RateLimitedNotifier{inner: inner, limiter: limiter}
}

// Notify forwards the call to the underlying notifier only when the
// derived change key is not suppressed by the rate limiter.
func (r *RateLimitedNotifier) Notify(changes snapshot.Changes) error {
	key := changeKey(changes)
	if !r.limiter.Allow(key) {
		return nil
	}
	return r.inner.Notify(changes)
}

// changeKey builds a stable string key from a Changes value so that
// identical change sets produce the same key regardless of map iteration
// order.
func changeKey(c snapshot.Changes) string {
	opened := sortedPorts(c.Opened)
	closed := sortedPorts(c.Closed)
	return fmt.Sprintf("opened:%s|closed:%s", strings.Join(opened, ","), strings.Join(closed, ","))
}

func sortedPorts(ports []string) []string {
	copy_ := append([]string(nil), ports...)
	sort.Strings(copy_)
	return copy_
}
