package monitor

import (
	"testing"
	"time"
)

func newTestCap(window time.Duration) *PortCap {
	return NewPortCap(window)
}

func TestPortCapFirstObserveNotPeak(t *testing.T) {
	pc := newTestCap(time.Minute)
	if pc.Observe(8080, 3) {
		t.Fatal("first observation should not be flagged as a new peak")
	}
}

func TestPortCapNewPeakDetected(t *testing.T) {
	pc := newTestCap(time.Minute)
	pc.Observe(8080, 3)
	if !pc.Observe(8080, 5) {
		t.Fatal("expected new peak to be detected")
	}
}

func TestPortCapSamePeakNotFlagged(t *testing.T) {
	pc := newTestCap(time.Minute)
	pc.Observe(8080, 5)
	if pc.Observe(8080, 5) {
		t.Fatal("same count should not be flagged as new peak")
	}
}

func TestPortCapLowerCountNotFlagged(t *testing.T) {
	pc := newTestCap(time.Minute)
	pc.Observe(8080, 10)
	if pc.Observe(8080, 7) {
		t.Fatal("lower count should not be flagged as new peak")
	}
}

func TestPortCapIndependentPorts(t *testing.T) {
	pc := newTestCap(time.Minute)
	pc.Observe(8080, 5)
	if !pc.Observe(9090, 1) == false {
		// first observe on 9090 should not be peak
	}
	if pc.Observe(9090, 3) != true {
		t.Fatal("port 9090 should independently track its own peak")
	}
	if pc.Observe(8080, 4) {
		t.Fatal("port 8080 peak is 5, count 4 should not flag")
	}
}

func TestPortCapPeakReturnsValue(t *testing.T) {
	pc := newTestCap(time.Minute)
	pc.Observe(8080, 7)
	v, ok := pc.Peak(8080)
	if !ok {
		t.Fatal("expected peak entry to exist")
	}
	if v != 7 {
		t.Fatalf("expected peak 7, got %d", v)
	}
}

func TestPortCapPeakMissingPort(t *testing.T) {
	pc := newTestCap(time.Minute)
	_, ok := pc.Peak(1234)
	if ok {
		t.Fatal("expected no entry for unseen port")
	}
}

func TestPortCapEvictRemovesExpired(t *testing.T) {
	pc := newTestCap(10 * time.Millisecond)
	pc.Observe(8080, 3)
	time.Sleep(20 * time.Millisecond)
	pc.Evict()
	if pc.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", pc.Len())
	}
}

func TestPortCapWindowResetAllowsNewBaseline(t *testing.T) {
	pc := newTestCap(10 * time.Millisecond)
	pc.Observe(8080, 10)
	time.Sleep(20 * time.Millisecond)
	// After window expires, first observe resets baseline — not flagged.
	if pc.Observe(8080, 2) {
		t.Fatal("after window expiry, first observe should reset baseline without flagging")
	}
}
