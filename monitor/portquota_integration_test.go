package monitor

import (
	"strings"
	"testing"
	"time"
)

// TestPortQuotaWithDiffMessages verifies that PortQuota flags a port that
// appears in repeated diffs more than the allowed quota.
func TestPortQuotaWithDiffMessages(t *testing.T) {
	q := NewPortQuota(2, 5*time.Second)

	snap1 := makeSnapshot(80, 443)
	snap2 := makeSnapshot(80, 443, 9999)

	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	// Observe the added port twice — should not yet exceed quota.
	for _, m := range msgs {
		if exceeded := q.Observe(m.Port); exceeded {
			t.Fatalf("expected quota not exceeded on first diff, port %d", m.Port)
		}
	}

	// Second diff cycle — same port added again (simulate flapping).
	for _, m := range msgs {
		if exceeded := q.Observe(m.Port); exceeded {
			t.Fatalf("expected quota not exceeded on second diff, port %d", m.Port)
		}
	}

	// Third diff cycle — quota should now be exceeded.
	exceeded := false
	for _, m := range msgs {
		if q.Observe(m.Port) {
			exceeded = true
		}
	}
	if !exceeded {
		t.Fatal("expected quota to be exceeded after third observation")
	}
}

// TestPortQuotaSummaryAfterDiff verifies that the summary reflects ports
// observed through diff messages.
func TestPortQuotaSummaryAfterDiff(t *testing.T) {
	q := NewPortQuota(5, 10*time.Second)

	snap1 := makeSnapshot(22)
	snap2 := makeSnapshot(22, 2222)

	msgs := Diff(snap1, snap2)
	for _, m := range msgs {
		q.Observe(m.Port)
	}

	summary := q.Summary()
	if !strings.Contains(summary, "2222") {
		t.Fatalf("expected summary to contain port 2222, got: %s", summary)
	}
}
