package monitor

import (
	"testing"
	"time"
)

func newTestScorer(t *testing.T) *PortScorer {
	t.Helper()
	labeler := NewPortLabeler(nil)
	anomaly := NewAnomalyDetector(3, time.Minute)
	dp := NewDeadPort(time.Minute)
	bl := NewBaseline("", 0)
	return NewPortScorer(labeler, anomaly, dp, bl)
}

func TestPortScorerUnknownPortAddsScore(t *testing.T) {
	ps := newTestScorer(t)
	// port 19999 is not a well-known port
	result := ps.Score(19999)
	if result.Score < 10 {
		t.Errorf("expected score >= 10 for unknown port, got %d", result.Score)
	}
}

func TestPortScorerBaselineKnownPortLowerScore(t *testing.T) {
	ps := newTestScorer(t)
	ps.baseline.Add(8080)
	result := ps.Score(8080)
	// baseline known: no "not in baseline" penalty
	for _, r := range result.Reasons {
		if r == "not in baseline" {
			t.Errorf("expected no baseline penalty for known port")
		}
	}
}

func TestPortScorerAnomalyRaisesScore(t *testing.T) {
	ps := newTestScorer(t)
	// trigger anomaly threshold
	for i := 0; i < 4; i++ {
		ps.anomaly.Record(9999)
	}
	result := ps.Score(9999)
	found := false
	for _, r := range result.Reasons {
		if r == "anomalous frequency" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected anomalous frequency reason, got %v", result.Reasons)
	}
}

func TestPortScorerLevelCritical(t *testing.T) {
	ps := newTestScorer(t)
	if ps.Level(70) != "critical" {
		t.Errorf("expected critical for score 70")
	}
}

func TestPortScorerLevelWarning(t *testing.T) {
	ps := newTestScorer(t)
	if ps.Level(50) != "warning" {
		t.Errorf("expected warning for score 50")
	}
}

func TestPortScorerLevelInfo(t *testing.T) {
	ps := newTestScorer(t)
	if ps.Level(5) != "info" {
		t.Errorf("expected info for score 5")
	}
}
