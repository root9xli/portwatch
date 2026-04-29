package monitor

import (
	"testing"
	"time"
)

func TestPortScoreHistoryRecordedFromScorer(t *testing.T) {
	scorer := newTestScorer()
	history := NewPortScoreHistory(10, 5*time.Minute)

	msgs := []Message{
		{Port: 8080, Action: "added"},
	}

	for _, m := range msgs {
		result := scorer.Score(m)
		history.Record(m.Port, result.Score, result.Level)
	}

	entries := history.Get(8080)
	if len(entries) == 0 {
		t.Fatal("expected at least one score history entry for port 8080")
	}
	if entries[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Port)
	}
}

func TestPortScoreHistoryMultiplePortsFromDiff(t *testing.T) {
	scorer := newTestScorer()
	history := NewPortScoreHistory(10, 5*time.Minute)

	ports := []int{80, 443, 9090}
	for _, p := range ports {
		m := Message{Port: p, Action: "added"}
		result := scorer.Score(m)
		history.Record(p, result.Score, result.Level)
	}

	for _, p := range ports {
		if len(history.Get(p)) == 0 {
			t.Errorf("expected history entry for port %d", p)
		}
	}
}

func TestPortScoreHistorySummaryAfterScoring(t *testing.T) {
	scorer := newTestScorer()
	history := NewPortScoreHistory(10, 5*time.Minute)

	m := Message{Port: 3306, Action: "added"}
	result := scorer.Score(m)
	history.Record(m.Port, result.Score, result.Level)

	summary := history.Summary(3306)
	if !contains(summary, "3306") {
		t.Errorf("expected summary to contain port 3306, got: %s", summary)
	}
}
