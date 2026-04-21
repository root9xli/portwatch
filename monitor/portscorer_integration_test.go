package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortScorerWithFullPipeline(t *testing.T) {
	labeler := NewPortLabeler(nil)
	anomaly := NewAnomalyDetector(2, time.Minute)
	dp := NewDeadPort(time.Minute)
	bl := NewBaseline("", 0)

	scorer := NewPortScorer(labeler, anomaly, dp, bl)
	reporter := NewPortScorerReporter(scorer)

	// Simulate anomalous port activity.
	for i := 0; i < 3; i++ {
		anomaly.Record(31337)
	}

	sp := reporter.Record(31337)
	if sp.Score < 40 {
		t.Errorf("expected elevated score for anomalous unknown port, got %d", sp.Score)
	}

	found := false
	for _, r := range sp.Reasons {
		if r == "anomalous frequency" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected anomalous frequency in reasons: %v", sp.Reasons)
	}
}

func TestPortScorerBaselineReducesRisk(t *testing.T) {
	labeler := NewPortLabeler(nil)
	anomaly := NewAnomalyDetector(10, time.Minute)
	dp := NewDeadPort(time.Minute)
	bl := NewBaseline("", 0)
	bl.Add(8080)

	scorer := NewPortScorer(labeler, anomaly, dp, bl)
	reporter := NewPortScorerReporter(scorer)

	sp := reporter.Record(8080)
	for _, r := range sp.Reasons {
		if r == "not in baseline" {
			t.Errorf("baseline port should not have 'not in baseline' reason")
		}
	}

	sum := reporter.Summary()
	if !strings.Contains(sum, "8080") {
		t.Errorf("expected 8080 in summary")
	}
}
