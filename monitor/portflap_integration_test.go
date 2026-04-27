package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortFlapDetectedFromRepeatedDiff(t *testing.T) {
	flap := NewPortFlap(3, time.Minute)
	reporter := NewPortFlapReporter()

	added := []Message{
		{Port: 8080, Action: "added"},
		{Port: 8080, Action: "removed"},
		{Port: 8080, Action: "added"},
	}

	now := time.Now()
	for i, msg := range added {
		ev := flap.Observe(msg.Port, now.Add(time.Duration(i)*time.Second))
		reporter.Record(ev)
	}

	if reporter.Len() != 1 {
		t.Fatalf("expected 1 flapping port, got %d", reporter.Len())
	}
	s := reporter.Summary()
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", s)
	}
}

func TestPortFlapStablePortNotFlagged(t *testing.T) {
	flap := NewPortFlap(4, time.Minute)
	reporter := NewPortFlapReporter()

	now := time.Now()
	// Only two transitions — below threshold of 4
	flap.Observe(443, now)
	ev := flap.Observe(443, now.Add(time.Second))
	reporter.Record(ev)

	if reporter.Len() != 0 {
		t.Errorf("stable port should not be flagged, got len=%d", reporter.Len())
	}
}

func TestPortFlapSummaryAfterMultiplePorts(t *testing.T) {
	flap := NewPortFlap(2, time.Minute)
	reporter := NewPortFlapReporter()

	now := time.Now()
	for _, port := range []int{80, 443} {
		flap.Observe(port, now)
		ev := flap.Observe(port, now.Add(time.Second))
		reporter.Record(ev)
	}

	s := reporter.Summary()
	if !strings.Contains(s, "80") || !strings.Contains(s, "443") {
		t.Errorf("expected both ports in summary, got: %s", s)
	}
}
