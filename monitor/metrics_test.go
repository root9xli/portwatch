package monitor

import (
	"testing"
	"time"
)

func TestMetricsInitialState(t *testing.T) {
	m := NewMetrics()
	if m.AlertsTotal != 0 || m.DiffsTotal != 0 || m.Suppressed != 0 || m.RateLimited != 0 {
		t.Error("expected all counters to be zero")
	}
	if m.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestMetricsRecordDiff(t *testing.T) {
	m := NewMetrics()
	m.RecordDiff()
	m.RecordDiff()
	if m.DiffsTotal != 2 {
		t.Errorf("expected 2 diffs, got %d", m.DiffsTotal)
	}
	if m.LastDiff().IsZero() {
		t.Error("expected LastDiff to be set after RecordDiff")
	}
}

func TestMetricsRecordAlert(t *testing.T) {
	m := NewMetrics()
	m.RecordAlert()
	if m.AlertsTotal != 1 {
		t.Errorf("expected 1 alert, got %d", m.AlertsTotal)
	}
}

func TestMetricsRecordSuppressedAndRateLimited(t *testing.T) {
	m := NewMetrics()
	m.RecordSuppressed()
	m.RecordRateLimited()
	m.RecordRateLimited()
	if m.Suppressed != 1 {
		t.Errorf("expected 1 suppressed, got %d", m.Suppressed)
	}
	if m.RateLimited != 2 {
		t.Errorf("expected 2 rate-limited, got %d", m.RateLimited)
	}
}

func TestMetricsSnapshot(t *testing.T) {
	m := NewMetrics()
	m.RecordAlert()
	m.RecordDiff()
	snap := m.Snapshot()
	if snap.AlertsTotal != 1 || snap.DiffsTotal != 1 {
		t.Error("snapshot does not reflect recorded values")
	}
}

func TestMetricsUptime(t *testing.T) {
	m := NewMetrics()
	time.Sleep(2 * time.Millisecond)
	if m.Uptime() < time.Millisecond {
		t.Error("expected uptime > 0")
	}
}
