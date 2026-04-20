package monitor

import (
	"testing"
	"time"
)

func newTestAnomaly(window time.Duration, thresh int) *AnomalyDetector {
	ad := NewAnomalyDetector(window, thresh)
	return ad
}

func TestAnomalyNotFlaggedBelowThreshold(t *testing.T) {
	ad := newTestAnomaly(time.Minute, 3)
	if ad.Record(8080) {
		t.Error("expected not anomalous on first event")
	}
	if ad.Record(8080) {
		t.Error("expected not anomalous on second event")
	}
}

func TestAnomalyFlaggedAtThreshold(t *testing.T) {
	ad := newTestAnomaly(time.Minute, 3)
	ad.Record(8080)
	ad.Record(8080)
	if !ad.Record(8080) {
		t.Error("expected anomalous at threshold")
	}
}

func TestAnomalyIsAnomalous(t *testing.T) {
	ad := newTestAnomaly(time.Minute, 2)
	ad.Record(9000)
	if ad.IsAnomalous(9000) {
		t.Error("should not be anomalous yet")
	}
	ad.Record(9000)
	if !ad.IsAnomalous(9000) {
		t.Error("should be anomalous after reaching threshold")
	}
}

func TestAnomalyIndependentPorts(t *testing.T) {
	ad := newTestAnomaly(time.Minute, 2)
	ad.Record(1111)
	ad.Record(1111)
	if ad.IsAnomalous(2222) {
		t.Error("port 2222 should not be affected by port 1111 events")
	}
}

func TestAnomalyEvictsStaleEvents(t *testing.T) {
	base := time.Now()
	ad := newTestAnomaly(time.Second*10, 2)

	// inject old events via fake clock
	old := base.Add(-time.Second * 20)
	ad.events[5000] = []time.Time{old, old, old}

	ad.now = func() time.Time { return base }
	ad.Evict()

	if ad.IsAnomalous(5000) {
		t.Error("stale events should have been evicted")
	}
	if _, ok := ad.events[5000]; ok {
		t.Error("port entry should be removed after full eviction")
	}
}

func TestAnomalyWindowExpiry(t *testing.T) {
	base := time.Now()
	call := 0
	ad := NewAnomalyDetector(time.Second*5, 2)
	ad.now = func() time.Time {
		call++
		if call <= 2 {
			return base
		}
		// advance past window
		return base.Add(time.Second * 10)
	}
	ad.Record(3000)
	ad.Record(3000)
	// now time has advanced, old events should be pruned on next Record
	if ad.Record(3000) {
		t.Error("events outside window should not count toward threshold")
	}
}
