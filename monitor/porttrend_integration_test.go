package monitor_test

import (
	"testing"
	"time"
)

func TestPortTrendRecordedFromDiff(t *testing.T) {
	trend := newTestTrend(t, 5, 30*time.Minute)

	snap1 := makeSnapshot(8080, 9090)
	snap2 := makeSnapshot(8080, 9090)

	diffs1 := Diff(nil, snap1)
	for _, d := range diffs1 {
		trend.Record(d.Port, 1)
	}

	diffs2 := Diff(snap1, snap2)
	_ = diffs2 // no changes expected

	trend.Record(8080, 1)
	trend.Record(8080, 1)

	dir := trend.Direction(8080)
	if dir == "" {
		t.Fatal("expected a direction for port 8080")
	}
}

func TestPortTrendStablePortNotFlagged(t *testing.T) {
	trend := newTestTrend(t, 5, 30*time.Minute)

	// Record the same count repeatedly — should be stable
	for i := 0; i < 5; i++ {
		trend.Record(443, 3)
	}

	dir := trend.Direction(443)
	if dir != "stable" {
		t.Errorf("expected stable direction for steady port 443, got %q", dir)
	}
}

func TestPortTrendReporterWithMultiplePorts(t *testing.T) {
	trend := newTestTrend(t, 5, 30*time.Minute)
	reporter := NewPortTrendReporter(trend)

	ports := []uint16{80, 443, 8080}
	for _, p := range ports {
		for i := 0; i < 3; i++ {
			trend.Record(p, uint32(i+1))
		}
	}

	summary := reporter.Summary()
	for _, p := range ports {
		pStr := itoa(int(p))
		if !containsSubstring(summary, pStr) {
			t.Errorf("expected summary to contain port %s, got: %s", pStr, summary)
		}
	}
}
