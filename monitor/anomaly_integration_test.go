package monitor

import (
	"testing"
	"time"
)

// TestAnomalyDetectedOnRepeatedDiff verifies that repeated add events for the
// same port across multiple Diff calls trigger the anomaly detector.
func TestAnomalyDetectedOnRepeatedDiff(t *testing.T) {
	ad := NewAnomalyDetector(time.Minute, 3)

	before := makeSnapshot([]int{80, 443})
	after1 := makeSnapshot([]int{80, 443, 9999})
	after2 := makeSnapshot([]int{80, 443})
	after3 := makeSnapshot([]int{80, 443, 9999})

	// simulate three churn cycles for port 9999
	for _, diff := range [][]Message{
		Diff(before, after1), // added 9999
		Diff(after1, after2), // removed 9999
		Diff(after2, after3), // added 9999 again
	} {
		for _, msg := range diff {
			ad.Record(msg.Port)
		}
	}

	if !ad.IsAnomalous(9999) {
		t.Error("port 9999 should be flagged as anomalous after repeated churn")
	}
}

// TestAnomalyStablePortNotFlagged ensures that a port that only appears once
// is never considered anomalous.
func TestAnomalyStablePortNotFlagged(t *testing.T) {
	ad := NewAnomalyDetector(time.Minute, 3)

	before := makeSnapshot([]int{80})
	after := makeSnapshot([]int{80, 8080})

	for _, msg := range Diff(before, after) {
		ad.Record(msg.Port)
	}

	if ad.IsAnomalous(8080) {
		t.Error("stable new port should not be anomalous after a single add")
	}
}
