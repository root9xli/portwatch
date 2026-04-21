package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTrendReporter() (*PortTrendReporter, *PortTrend, *time.Time) {
	pt, now := newTestTrend(time.Minute)
	labels := NewPortLabeler(nil)
	rep := NewPortTrendReporter(pt, labels)
	return rep, pt, now
}

func TestPortTrendReporterEmptySummary(t *testing.T) {
	rep, _, _ := newTrendReporter()
	got := rep.Summary(nil)
	if got != "no trend data" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestPortTrendReporterContainsPort(t *testing.T) {
	rep, pt, _ := newTrendReporter()
	pt.Record(8080, 1)
	got := rep.Summary([]uint16{8080})
	if !strings.Contains(got, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", got)
	}
}

func TestPortTrendReporterContainsDirection(t *testing.T) {
	rep, pt, now := newTrendReporter()
	pt.Record(9000, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 1)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 8)
	*now = now.Add(5 * time.Second)
	pt.Record(9000, 9)
	got := rep.Summary([]uint16{9000})
	if !strings.Contains(got, string(TrendUp)) {
		t.Errorf("expected 'up' in summary, got: %s", got)
	}
}

func TestPortTrendReporterSortedByPort(t *testing.T) {
	rep, pt, _ := newTrendReporter()
	pt.Record(9000, 1)
	pt.Record(80, 1)
	entries := rep.Entries([]uint16{9000, 80})
	if entries[0].Port != 80 || entries[1].Port != 9000 {
		t.Errorf("expected sorted by port, got %d %d", entries[0].Port, entries[1].Port)
	}
}

func TestPortTrendReporterLabelIncluded(t *testing.T) {
	rep, pt, _ := newTrendReporter()
	pt.Record(80, 1)
	got := rep.Summary([]uint16{80})
	if !strings.Contains(got, "http") {
		t.Errorf("expected label 'http' in summary, got: %s", got)
	}
}
