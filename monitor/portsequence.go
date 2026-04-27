package monitor

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// PortSequenceEvent describes a run of consecutive port numbers observed together.
type PortSequenceEvent struct {
	Start     int
	End       int
	Length    int
	ObservedAt time.Time
}

func (e *PortSequenceEvent) String() string {
	return fmt.Sprintf("sequence %d-%d (len=%d)", e.Start, e.End, e.Length)
}

// PortSequence detects when a run of consecutive ports appears in a diff,
// which may indicate a port scan or misconfigured service binding.
type PortSequence struct {
	mu        sync.Mutex
	minRun    int
	events    []*PortSequenceEvent
	maxEvents int
}

// NewPortSequence creates a PortSequence detector.
// minRun is the minimum consecutive-port run length to flag.
func NewPortSequence(minRun, maxEvents int) *PortSequence {
	if minRun < 2 {
		minRun = 2
	}
	if maxEvents < 1 {
		maxEvents = 64
	}
	return &PortSequence{minRun: minRun, maxEvents: maxEvents}
}

// Observe checks a slice of ports for consecutive runs and records any events.
// Returns the detected events (may be empty).
func (ps *PortSequence) Observe(ports []int) []*PortSequenceEvent {
	if len(ports) == 0 {
		return nil
	}
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	var found []*PortSequenceEvent
	start := sorted[0]
	prev := sorted[0]

	flush := func(end int) {
		run := end - start + 1
		if run >= ps.minRun {
			found = append(found, &PortSequenceEvent{
				Start:      start,
				End:        end,
				Length:     run,
				ObservedAt: time.Now(),
			})
		}
	}

	for _, p := range sorted[1:] {
		if p == prev+1 {
			prev = p
			continue
		}
		flush(prev)
		start = p
		prev = p
	}
	flush(prev)

	if len(found) == 0 {
		return nil
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.events = append(ps.events, found...)
	if len(ps.events) > ps.maxEvents {
		ps.events = ps.events[len(ps.events)-ps.maxEvents:]
	}
	return found
}

// Events returns all stored sequence events.
func (ps *PortSequence) Events() []*PortSequenceEvent {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	out := make([]*PortSequenceEvent, len(ps.events))
	copy(out, ps.events)
	return out
}

// Len returns the number of stored events.
func (ps *PortSequence) Len() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.events)
}
