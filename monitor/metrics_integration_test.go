package monitor

import (
	"strings"
	"testing"
)

// TestMetricsRecordedOnDiff verifies that Metrics.RecordDiff integrates with
// the Diff function output by simulating a monitoring cycle.
func TestMetricsRecordedOnDiff(t *testing.T) {
	prev := makeSnapshot([]int{8080})
	curr := makeSnapshot([]int{8080, 9090})

	m := NewMetrics()
	changes := Diff(prev, curr)
	if len(changes) > 0 {
		m.RecordDiff()
		for range changes {
			m.RecordAlert()
		}
	}

	if m.DiffsTotal != 1 {
		t.Errorf("expected 1 diff, got %d", m.DiffsTotal)
	}
	if m.AlertsTotal != 1 {
		t.Errorf("expected 1 alert for new port, got %d", m.AlertsTotal)
	}
}

// TestMetricsReporterAfterCycle verifies the reporter summary after a full
// simulated monitor cycle with suppression and rate-limiting recorded.
func TestMetricsReporterAfterCycle(t *testing.T) {
	m := NewMetrics()
	m.RecordDiff()
	m.RecordAlert()
	m.RecordSuppressed()
	m.RecordSuppressed()
	m.RecordRateLimited()

	r := NewMetricsReporter(m)
	summary := r.Summary()

	if !strings.Contains(summary, "diffs:        1") {
		t.Errorf("unexpected diffs line in summary: %s", summary)
	}
	if !strings.Contains(summary, "alerts:       1") {
		t.Errorf("unexpected alerts line in summary: %s", summary)
	}
	if !strings.Contains(summary, "suppressed:   2") {
		t.Errorf("unexpected suppressed line in summary: %s", summary)
	}
	if !strings.Contains(summary, "rate_limited: 1") {
		t.Errorf("unexpected rate_limited line in summary: %s", summary)
	}
}
