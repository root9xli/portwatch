package monitor

import (
	"testing"
	"time"
)

func newTestTrend(window time.Duration) (*PortTrend, *time.Time) {
	pt := NewPortTrend(window)
	now := time.Now()
	pt.now = func() time.Time { return now }
	return pt, &now
}

func TestPortTrendStableWithOnePoint(t *testing.T) {
	pt, _ := newTestTrend(time.Minute)
	pt.Record(8080, 3)
	if got := pt.Trend(8080); got != TrendStable {
		t.Errorf("expected stable, got %s", got)
	}
}

func TestPortTrendUp(t *testing.T) {
	pt, now := newTestTrend(time.Minute)
	pt.Record(9000, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 5)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 6)
	if got := pt.Trend(9000); got != TrendUp {
		t.Errorf("expected up, got %s", got)
	}
}

func TestPortTrendDown(t *testing.T) {
	pt, now := newTestTrend(time.Minute)
	pt.Record(443, 10)
	*now = now.Add(5 * time.Second)
	pt.Record(443, 9)
	*now = now.Add(5 * time.Second)
	pt.Record(443, 2)
	*now = now.Add(5 * time.Second)
	pt.Record(443, 1)
	if got := pt.Trend(443); got != TrendDown {
		t.Errorf("expected down, got %s", got)
	}
}

func TestPortTrendEvictsOldSamples(t *testing.T) {
	pt, now := newTestTrend(10 * time.Second)
	pt.Record(22, 10)
	*now = now.Add(20 * time.Second) // beyond window
	pt.Record(22, 1)
	// Only one sample survives eviction, so trend is stable
	if got := pt.Trend(22); got != TrendStable {
		t.Errorf("expected stable after eviction, got %s", got)
	}
}

func TestPortTrendIndependentPorts(t *testing.T) {
	pt, now := newTestTrend(time.Minute)
	pt.Record(80, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(80, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(80, 5)
	*now = now.Add(5 * time.Second)
	pt.Record(80, 6)

	pt.Record(443, 5)
	if got := pt.Trend(443); got != TrendStable {
		t.Errorf("port 443 should be stable, got %s", got)
	}
	if got := pt.Trend(80); got != TrendUp {
		t.Errorf("port 80 should be up, got %s", got)
	}
}
