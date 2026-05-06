package metrics_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/metrics"
)

func TestNew_InitialisesUptimeSince(t *testing.T) {
	before := time.Now()
	c := metrics.New()
	after := time.Now()

	s := c.Get()
	if s.UptimeSince.Before(before) || s.UptimeSince.After(after) {
		t.Errorf("UptimeSince %v not in expected range [%v, %v]", s.UptimeSince, before, after)
	}
}

func TestRecordPoll_IncrementsCounters(t *testing.T) {
	c := metrics.New()
	c.RecordPoll(50*time.Millisecond, 3, 1)
	c.RecordPoll(20*time.Millisecond, 0, 2)

	s := c.Get()
	if s.PollCount != 2 {
		t.Errorf("expected PollCount 2, got %d", s.PollCount)
	}
	if s.OpenedTotal != 3 {
		t.Errorf("expected OpenedTotal 3, got %d", s.OpenedTotal)
	}
	if s.ClosedTotal != 3 {
		t.Errorf("expected ClosedTotal 3, got %d", s.ClosedTotal)
	}
	if s.LastPollDur != 20*time.Millisecond {
		t.Errorf("expected LastPollDur 20ms, got %v", s.LastPollDur)
	}
}

func TestRecordAlert_IncrementsAlertsSent(t *testing.T) {
	c := metrics.New()
	c.RecordAlert()
	c.RecordAlert()

	if got := c.Get().AlertsSent; got != 2 {
		t.Errorf("expected AlertsSent 2, got %d", got)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := metrics.New()
	c.RecordPoll(10*time.Millisecond, 1, 0)

	s1 := c.Get()
	c.RecordPoll(10*time.Millisecond, 1, 0)
	s2 := c.Get()

	if s1.PollCount == s2.PollCount {
		t.Error("expected Get to return independent copies")
	}
}
