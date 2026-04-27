package monitor

import (
	"testing"
)

func newTestSequence() *PortSequence {
	return NewPortSequence(3, 32)
}

func TestPortSequenceNoRunBelowMinRun(t *testing.T) {
	ps := newTestSequence() // minRun=3
	events := ps.Observe([]int{80, 81}) // run of 2, below threshold
	if len(events) != 0 {
		t.Fatalf("expected no events for run < minRun, got %d", len(events))
	}
}

func TestPortSequenceDetectsRunAtMinRun(t *testing.T) {
	ps := newTestSequence()
	events := ps.Observe([]int{100, 101, 102})
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e := events[0]
	if e.Start != 100 || e.End != 102 || e.Length != 3 {
		t.Errorf("unexpected event: %+v", e)
	}
}

func TestPortSequenceDetectsLongerRun(t *testing.T) {
	ps := newTestSequence()
	events := ps.Observe([]int{8080, 8081, 8082, 8083, 8084})
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Length != 5 {
		t.Errorf("expected length 5, got %d", events[0].Length)
	}
}

func TestPortSequenceMultipleRuns(t *testing.T) {
	ps := newTestSequence()
	// two separate runs: 10-12 and 20-22
	events := ps.Observe([]int{10, 11, 12, 20, 21, 22})
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestPortSequenceNonConsecutiveNoEvent(t *testing.T) {
	ps := newTestSequence()
	events := ps.Observe([]int{80, 443, 8080})
	if len(events) != 0 {
		t.Fatalf("expected no events for non-consecutive ports, got %d", len(events))
	}
}

func TestPortSequenceEventsStoredAndReturned(t *testing.T) {
	ps := newTestSequence()
	ps.Observe([]int{1, 2, 3})
	ps.Observe([]int{5, 6, 7})
	if ps.Len() != 2 {
		t.Errorf("expected Len 2, got %d", ps.Len())
	}
	events := ps.Events()
	if len(events) != 2 {
		t.Errorf("expected 2 events from Events(), got %d", len(events))
	}
}

func TestPortSequenceEvictsOldestWhenFull(t *testing.T) {
	ps := NewPortSequence(2, 3) // maxEvents=3
	ps.Observe([]int{1, 2})
	ps.Observe([]int{10, 11})
	ps.Observe([]int{20, 21})
	ps.Observe([]int{30, 31}) // should evict first
	if ps.Len() != 3 {
		t.Errorf("expected Len 3 after eviction, got %d", ps.Len())
	}
	events := ps.Events()
	if events[0].Start != 10 {
		t.Errorf("expected oldest evicted, first event start=10, got %d", events[0].Start)
	}
}

func TestPortSequenceUnsortedInputHandled(t *testing.T) {
	ps := newTestSequence()
	events := ps.Observe([]int{303, 301, 302})
	if len(events) != 1 {
		t.Fatalf("expected 1 event for unsorted consecutive input, got %d", len(events))
	}
	if events[0].Start != 301 || events[0].End != 303 {
		t.Errorf("unexpected event range: %+v", events[0])
	}
}

func TestPortSequenceEmptyInputReturnsNil(t *testing.T) {
	ps := newTestSequence()
	events := ps.Observe([]int{})
	if events != nil {
		t.Errorf("expected nil for empty input, got %v", events)
	}
}
