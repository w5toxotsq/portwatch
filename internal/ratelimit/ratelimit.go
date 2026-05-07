package ratelimit

import (
	"sync"
	"time"
)

// Limiter throttles alert dispatch so that repeated identical change sets
// do not flood downstream notifiers.
type Limiter struct {
	mu       sync.Mutex
	window   time.Duration
	last     map[string]time.Time
	clock    func() time.Time
}

// New returns a Limiter that suppresses duplicate keys within window.
func New(window time.Duration) *Limiter {
	return &Limiter{
		window: window,
		last:   make(map[string]time.Time),
		clock:  time.Now,
	}
}

// Allow reports whether the given key should be allowed through.
// A key is allowed if it has never been seen, or if the window has
// elapsed since it was last allowed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.window {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears all recorded keys, allowing every key through on the
// next call to Allow.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// Len returns the number of keys currently tracked.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.last)
}
