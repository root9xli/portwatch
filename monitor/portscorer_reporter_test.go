package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestScorerReporter(t *testing.T) *PortScorerReporter {
	t.Helper()
	labeler := NewPortLabeler(nil)
	anomaly := NewAnomalyDetector(3, time.Minute)
	dp := NewDeadPort(time.Minute)
	bl := NewBaseline("", 0)
	scorer := NewPortScorer(labeler, anomaly, dp, bl)
	return NewPortScorerReporter(scorer)
}

func TestScorerReporterEmptySummary(t *testing.T) {
	r := newTestScorerReporter(t)
	if !strings.Contains(r.Summary(), "no ports recorded") {
		t.Errorf("expected empty summary message")
	}
}

func TestScorerReporterRecordAppearsInSummary(t *testing.T) {
	r := newTestScorerReporter(t)
	r.Record(19999)
	sum := r.Summary()
	if !strings.Contains(sum, "19999") {
		t.Errorf("expected port 19999 in summary, got: %s", sum)
	}
}

func TestScorerReporterTopLimitN(t *testing.T) {
	r := newTestScorerReporter(t)
	for _, p := range []int{1111, 2222, 3333, 4444, 5555, 6666} {
		r.Record(p)
	}
	top := r.Top(3)
	if len(top) != 3 {
		t.Errorf("expected 3 results, got %d", len(top))
	}
}

func TestScorerReporterTopSortedDescending(t *testing.T) {
	r := newTestScorerReporter(t)
	// 80 is well-known (http-alt), 19999 is unknown
	r.Record(80)
	r.Record(19999)
	top := r.Top(2)
	if top[0].Score < top[1].Score {
		t.Errorf("expected descending order: %+v", top)
	}
}

func TestScorerReporterSummaryContainsLevel(t *testing.T) {
	r := newTestScorerReporter(t)
	r.Record(19999)
	sum := r.Summary()
	if !strings.Contains(sum, "level=") {
		t.Errorf("expected level= in summary, got: %s", sum)
	}
}
