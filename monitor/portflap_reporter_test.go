package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestFlapReporter() (*PortFlap, *PortFlapReporter) {
	return NewPortFlap(2, time.Minute), NewPortFlapReporter()
}

func TestPortFlapReporterEmptySummary(t *testing.T) {
	_, r := newTestFlapReporter()
	s := r.Summary()
	if !strings.Contains(s, "no flapping") {
		t.Errorf("expected empty summary, got: %s", s)
	}
}

func TestPortFlapReporterRecordAppearsInSummary(t *testing.T) {
	f, r := newTestFlapReporter()
	now := time.Now()
	f.Observe(8080, now)
	ev := f.Observe(8080, now.Add(time.Second))
	if ev == nil {
		t.Fatal("expected flap event")
	}
	r.Record(ev)
	s := r.Summary()
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", s)
	}
}

func TestPortFlapReporterLen(t *testing.T) {
	f, r := newTestFlapReporter()
	now := time.Now()
	for _, port := range []int{80, 443, 9000} {
		f.Observe(port, now)
		ev := f.Observe(port, now.Add(time.Second))
		r.Record(ev)
	}
	if r.Len() != 3 {
		t.Errorf("expected len 3, got %d", r.Len())
	}
}

func TestPortFlapReporterRecordNilNoOp(t *testing.T) {
	_, r := newTestFlapReporter()
	r.Record(nil)
	if r.Len() != 0 {
		t.Errorf("expected len 0 after nil record, got %d", r.Len())
	}
}

func TestPortFlapReporterSummaryContainsFlaps(t *testing.T) {
	f, r := newTestFlapReporter()
	now := time.Now()
	f.Observe(3000, now)
	ev := f.Observe(3000, now.Add(time.Second))
	r.Record(ev)
	s := r.Summary()
	if !strings.Contains(s, "flaps=") {
		t.Errorf("expected flaps= in summary, got: %s", s)
	}
}
