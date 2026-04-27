package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortSpikeDetectedFromRepeatedDiff(t *testing.T) {
	ps := NewPortSpike(500*time.Millisecond, 3)
	rep := NewPortSpikeReporter()

	added := []Message{
		{Port: 4444, Action: "added", Level: "warn"},
	}

	for i := 0; i < 4; i++ {
		for _, msg := range added {
			ev := ps.Observe(msg.Port)
			rep.Record(ev)
		}
	}

	if rep.Len() == 0 {
		t.Fatal("expected at least one spike event recorded")
	}

	ev := rep.TopN(1)
	if len(ev) == 0 || ev[0].Port != 4444 {
		t.Errorf("expected top spike for port 4444")
	}
}

func TestPortSpikeStablePortNotFlagged(t *testing.T) {
	ps := NewPortSpike(500*time.Millisecond, 10)
	rep := NewPortSpikeReporter()

	for i := 0; i < 5; i++ {
		rep.Record(ps.Observe(3333))
	}

	if rep.Len() != 0 {
		t.Errorf("expected no spikes below threshold, got %d", rep.Len())
	}
}

func TestPortSpikeSummaryAfterDiff(t *testing.T) {
	ps := NewPortSpike(500*time.Millisecond, 2)
	rep := NewPortSpikeReporter()

	msgs := []Message{
		{Port: 8888, Action: "added", Level: "warn"},
		{Port: 8888, Action: "added", Level: "warn"},
		{Port: 8888, Action: "added", Level: "warn"},
	}
	for _, m := range msgs {
		rep.Record(ps.Observe(m.Port))
	}

	summary := rep.Summary()
	if !strings.Contains(summary, "8888") {
		t.Errorf("expected port 8888 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "spike") {
		t.Errorf("expected 'spike' in summary, got: %s", summary)
	}
}
