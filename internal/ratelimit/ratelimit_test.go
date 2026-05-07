package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if !l.Allow("tcp:80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_DuplicateWithinWindowBlocked(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:80")
	if l.Allow("tcp:80") {
		t.Fatal("expected duplicate within window to be blocked")
	}
}

func TestAllow_DuplicateAfterWindowAllowed(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(1 * time.Second)

	// Manually inject a stale timestamp via the exported clock hook.
	// We use Reset + a fake clock instead.
	_ = now

	// Simulate elapsed window by using a tiny window and sleeping.
	l2 := ratelimit.New(10 * time.Millisecond)
	l2.Allow("tcp:443")
	time.Sleep(20 * time.Millisecond)
	if !l2.Allow("tcp:443") {
		t.Fatal("expected key to be allowed after window elapsed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:80")
	if !l.Allow("tcp:443") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:80")
	l.Reset()
	if !l.Allow("tcp:80") {
		t.Fatal("expected key to be allowed after reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:80")
	l.Allow("tcp:443")
	l.Allow("udp:53")
	if got := l.Len(); got != 3 {
		t.Fatalf("expected 3 tracked keys, got %d", got)
	}
}

func TestLen_AfterReset(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:80")
	l.Reset()
	if got := l.Len(); got != 0 {
		t.Fatalf("expected 0 tracked keys after reset, got %d", got)
	}
}
