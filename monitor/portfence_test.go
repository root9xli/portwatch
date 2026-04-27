package monitor

import (
	"testing"
	"time"
)

func newTestFence() (*PortFence, *time.Time) {
	now := time.Now()
	f := NewPortFence()
	f.clock = func() time.Time { return now }
	return f, &now
}

func TestPortFenceNotFencedInitially(t *testing.T) {
	f, _ := newTestFence()
	if f.IsFenced(8080) {
		t.Fatal("expected port 8080 not to be fenced")
	}
}

func TestPortFenceFencedAfterFence(t *testing.T) {
	f, _ := newTestFence()
	f.Fence(8080, 5*time.Minute, "test")
	if !f.IsFenced(8080) {
		t.Fatal("expected port 8080 to be fenced")
	}
}

func TestPortFenceLiftRemovesFence(t *testing.T) {
	f, _ := newTestFence()
	f.Fence(8080, 5*time.Minute, "test")
	f.Lift(8080)
	if f.IsFenced(8080) {
		t.Fatal("expected port 8080 not to be fenced after lift")
	}
}

func TestPortFenceExpiresAfterDuration(t *testing.T) {
	f, now := newTestFence()
	f.Fence(8080, 1*time.Second, "expire-test")
	*now = now.Add(2 * time.Second)
	if f.IsFenced(8080) {
		t.Fatal("expected port 8080 fence to have expired")
	}
}

func TestPortFenceGetReturnsEntry(t *testing.T) {
	f, _ := newTestFence()
	f.Fence(9000, 10*time.Minute, "reason-check")
	e := f.Get(9000)
	if e == nil {
		t.Fatal("expected non-nil entry")
	}
	if e.Reason != "reason-check" {
		t.Fatalf("expected reason 'reason-check', got %q", e.Reason)
	}
	if e.Port != 9000 {
		t.Fatalf("expected port 9000, got %d", e.Port)
	}
}

func TestPortFenceGetMissingReturnsNil(t *testing.T) {
	f, _ := newTestFence()
	if f.Get(1234) != nil {
		t.Fatal("expected nil for unfenced port")
	}
}

func TestPortFenceLenCountsActive(t *testing.T) {
	f, now := newTestFence()
	f.Fence(80, 5*time.Minute, "a")
	f.Fence(443, 5*time.Minute, "b")
	f.Fence(8080, 1*time.Second, "c")
	*now = now.Add(2 * time.Second)
	if got := f.Len(); got != 2 {
		t.Fatalf("expected 2 active fences, got %d", got)
	}
}

func TestPortFenceFilterRemovesFenced(t *testing.T) {
	f, _ := newTestFence()
	f.Fence(8080, 5*time.Minute, "block")
	msgs := []Message{
		{Port: 8080, Action: "added"},
		{Port: 443, Action: "added"},
	}
	out := f.Filter(msgs)
	if len(out) != 1 {
		t.Fatalf("expected 1 message after filter, got %d", len(out))
	}
	if out[0].Port != 443 {
		t.Fatalf("expected port 443 in output, got %d", out[0].Port)
	}
}

func TestPortFenceIndependentPorts(t *testing.T) {
	f, _ := newTestFence()
	f.Fence(80, 5*time.Minute, "web")
	if f.IsFenced(443) {
		t.Fatal("port 443 should not be fenced")
	}
}
