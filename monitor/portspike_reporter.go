package monitor

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// PortSpikeReporter records PortSpikeEvents and produces summaries.
type PortSpikeReporter struct {
	mu     sync.Mutex
	events []*PortSpikeEvent
}

// NewPortSpikeReporter creates an empty PortSpikeReporter.
func NewPortSpikeReporter() *PortSpikeReporter {
	return &PortSpikeReporter{}
}

// Record stores a spike event. Nil events are ignored.
func (r *PortSpikeReporter) Record(ev *PortSpikeEvent) {
	if ev == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, ev)
}

// Len returns the number of recorded spike events.
func (r *PortSpikeReporter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}

// TopN returns up to n ports with the highest spike counts, sorted descending.
func (r *PortSpikeReporter) TopN(n int) []*PortSpikeEvent {
	r.mu.Lock()
	defer r.mu.Unlock()

	sorted := make([]*PortSpikeEvent, len(r.events))
	copy(sorted, r.events)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})
	if n < len(sorted) {
		return sorted[:n]
	}
	return sorted
}

// Summary returns a human-readable report of recorded spike events.
func (r *PortSpikeReporter) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.events) == 0 {
		return "portspike: no spikes detected"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("portspike: %d spike(s) detected\n", len(r.events)))
	for _, ev := range r.events {
		sb.WriteString(fmt.Sprintf("  port=%d count=%d threshold=%d at=%s\n",
			ev.Port, ev.Count, ev.Threshold, ev.DetectedAt.Format("15:04:05")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
