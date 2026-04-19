package monitor

import (
	"strings"
	"testing"
)

func TestBaselineReporterEmptySummary(t *testing.T) {
	b := NewBaseline("")
	r := NewBaselineReporter(b)
	if !strings.Contains(r.Summary(), "no ports") {
		t.Fatal("expected 'no ports' in empty summary")
	}
}

func TestBaselineReporterSummaryContainsPort(t *testing.T) {
	b := NewBaseline("")
	b.Add(8080)
	r := NewBaselineReporter(b)
	if !strings.Contains(r.Summary(), "8080") {
		t.Fatal("expected port 8080 in summary")
	}
}

func TestBaselineReporterContainsPort(t *testing.T) {
	b := NewBaseline("")
	b.Add(443)
	r := NewBaselineReporter(b)
	out := r.ContainsPort(443)
	if !strings.Contains(out, "in baseline") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestBaselineReporterNotContainsPort(t *testing.T) {
	b := NewBaseline("")
	r := NewBaselineReporter(b)
	out := r.ContainsPort(9999)
	if !strings.Contains(out, "NOT in baseline") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestBaselineReporterCount(t *testing.T) {
	b := NewBaseline("")
	b.Add(22)
	b.Add(80)
	b.Add(443)
	r := NewBaselineReporter(b)
	if !strings.Contains(r.Summary(), "3 port") {
		t.Fatal("expected count 3 in summary")
	}
}
