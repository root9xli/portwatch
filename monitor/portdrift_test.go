package monitor

import (
	"testing"
	"time"
)

func newTestDrift(threshold int, window time.Duration) *PortDrift {
	d := NewPortDrift(threshold, window)
	base := time.Now()
	d.now = func() time.Time { return base }
	return d
}

func TestPortDriftNotFlappingBelowThreshold(t *testing.T) {
	d := newTestDrift(3, time.Minute)
	for i := 0; i < 2; i++ {
		if d.Observe(8080) {
			t.Fatal("expected no flap below threshold")
		}
	}
	if d.IsFlapping(8080) {
		t.Fatal("should not be flapping below threshold")
	}
}

func TestPortDriftFlapsAtThreshold(t *testing.T) {
	d := newTestDrift(3, time.Minute)
	var flapped bool
	for i := 0; i < 3; i++ {
		flapped = d.Observe(8080)
	}
	if !flapped {
		t.Fatal("expected flap at threshold")
	}
	if !d.IsFlapping(8080) {
		t.Fatal("IsFlapping should return true after threshold")
	}
}

func TestPortDriftIndependentPorts(t *testing.T) {
	d := newTestDrift(2, time.Minute)
	d.Observe(9090)
	d.Observe(9090)
	if d.IsFlapping(8080) {
		t.Fatal("unrelated port should not be flagged")
	}
}

func TestPortDriftResetsAfterWindow(t *testing.T) {
	base := time.Now()
	d := NewPortDrift(2, time.Second)
	d.now = func() time.Time { return base }
	d.Observe(7070)
	d.Observe(7070)
	// advance past window
	d.now = func() time.Time { return base.Add(2 * time.Second) }
	d.Observe(7070) // should reset
	if d.IsFlapping(7070) {
		t.Fatal("should not be flapping after window reset")
	}
}

func TestPortDriftEvictRemovesExpired(t *testing.T) {
	base := time.Now()
	d := NewPortDrift(2, time.Second)
	d.now = func() time.Time { return base }
	d.Observe(6060)
	d.Observe(6060)
	d.now = func() time.Time { return base.Add(5 * time.Second) }
	d.Evict()
	if d.IsFlapping(6060) {
		t.Fatal("evicted entry should not show as flapping")
	}
}

func TestPortDriftSummaryNoFlapping(t *testing.T) {
	d := newTestDrift(5, time.Minute)
	s := d.Summary()
	if s != "no flapping ports" {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestPortDriftSummaryContainsPort(t *testing.T) {
	d := newTestDrift(2, time.Minute)
	d.Observe(3000)
	d.Observe(3000)
	s := d.Summary()
	if s == "no flapping ports" {
		t.Fatal("expected flapping port in summary")
	}
}
