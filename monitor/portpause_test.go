package monitor

import (
	"testing"
	"time"
)

func newTestPause() *PortPause {
	p := NewPortPause()
	p.now = func() time.Time { return time.Unix(1000, 0) }
	return p
}

func TestPortPauseNotPausedInitially(t *testing.T) {
	p := newTestPause()
	if p.IsPaused(8080) {
		t.Fatal("expected port not to be paused initially")
	}
}

func TestPortPausePausedAfterPause(t *testing.T) {
	p := newTestPause()
	p.Pause(8080, 10*time.Second, "maintenance")
	if !p.IsPaused(8080) {
		t.Fatal("expected port to be paused")
	}
}

func TestPortPauseResumeRemovesPause(t *testing.T) {
	p := newTestPause()
	p.Pause(8080, 10*time.Second, "test")
	p.Resume(8080)
	if p.IsPaused(8080) {
		t.Fatal("expected port to no longer be paused after resume")
	}
}

func TestPortPauseExpiresAfterDuration(t *testing.T) {
	base := time.Unix(1000, 0)
	p := NewPortPause()
	p.now = func() time.Time { return base }
	p.Pause(8080, 5*time.Second, "expiry test")

	p.now = func() time.Time { return base.Add(6 * time.Second) }
	if p.IsPaused(8080) {
		t.Fatal("expected pause to have expired")
	}
}

func TestPortPauseGetReturnsEntry(t *testing.T) {
	p := newTestPause()
	p.Pause(8080, 10*time.Second, "reason-x")
	e := p.Get(8080)
	if e == nil {
		t.Fatal("expected non-nil entry")
	}
	if e.Reason != "reason-x" {
		t.Fatalf("expected reason 'reason-x', got %q", e.Reason)
	}
}

func TestPortPauseIndependentPorts(t *testing.T) {
	p := newTestPause()
	p.Pause(8080, 10*time.Second, "a")
	if p.IsPaused(9090) {
		t.Fatal("port 9090 should not be paused")
	}
}

func TestPortPauseEvictRemovesExpired(t *testing.T) {
	base := time.Unix(1000, 0)
	p := NewPortPause()
	p.now = func() time.Time { return base }
	p.Pause(8080, 2*time.Second, "evict")
	p.Pause(9090, 20*time.Second, "keep")

	p.now = func() time.Time { return base.Add(5 * time.Second) }
	p.Evict()
	if p.Len() != 1 {
		t.Fatalf("expected 1 entry after evict, got %d", p.Len())
	}
}

func TestPortPauseFilterRemovesPausedPorts(t *testing.T) {
	p := newTestPause()
	p.Pause(8080, 10*time.Second, "filtered")
	msgs := []Message{
		{Port: 8080, Action: "added"},
		{Port: 9090, Action: "added"},
	}
	out := p.Filter(msgs)
	if len(out) != 1 {
		t.Fatalf("expected 1 message after filter, got %d", len(out))
	}
	if out[0].Port != 9090 {
		t.Fatalf("expected port 9090, got %d", out[0].Port)
	}
}
