package monitor

import (
	"strings"
	"testing"
)

func TestReporterSummaryEmpty(t *testing.T) {
	h := NewHistory(10)
	r := NewReporter(h)
	out := r.Summary(5)
	if out != "no recent events" {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestReporterSummaryContainsPort(t *testing.T) {
	h := NewHistory(10)
	h.Record(8080, "added", "warn")
	r := NewReporter(h)
	out := r.Summary(5)
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", out)
	}
}

func TestReporterSummaryContainsLevel(t *testing.T) {
	h := NewHistory(10)
	h.Record(443, "added", "critical")
	r := NewReporter(h)
	out := r.Summary(1)
	if !strings.Contains(out, "critical") {
		t.Errorf("expected level in summary, got: %s", out)
	}
}

func TestReporterCountByLevel(t *testing.T) {
	h := NewHistory(10)
	h.Record(80, "added", "info")
	h.Record(443, "added", "warn")
	h.Record(8080, "added", "warn")
	h.Record(9090, "added", "critical")
	r := NewReporter(h)
	counts := r.CountByLevel()
	if counts["warn"] != 2 {
		t.Errorf("expected 2 warn, got %d", counts["warn"])
	}
	if counts["info"] != 1 {
		t.Errorf("expected 1 info, got %d", counts["info"])
	}
	if counts["critical"] != 1 {
		t.Errorf("expected 1 critical, got %d", counts["critical"])
	}
}

func TestReporterHasCritical(t *testing.T) {
	h := NewHistory(10)
	r := NewReporter(h)
	if r.HasCritical() {
		t.Error("expected no critical on empty history")
	}
	h.Record(22, "added", "critical")
	if !r.HasCritical() {
		t.Error("expected HasCritical to be true")
	}
}
