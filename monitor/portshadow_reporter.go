package monitor

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// PortShadowReporter aggregates ShadowEvents and produces structured reports.
type PortShadowReporter struct {
	mu     sync.Mutex
	events []ShadowEvent
}

// NewPortShadowReporter creates a new PortShadowReporter.
func NewPortShadowReporter() *PortShadowReporter {
	return &PortShadowReporter{}
}

// Record adds a ShadowEvent to the reporter.
func (r *PortShadowReporter) Record(ev ShadowEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, ev)
}

// Len returns the number of recorded events.
func (r *PortShadowReporter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}

// TopN returns the N ports with the most shadow observations, sorted descending.
func (r *PortShadowReporter) TopN(n int) []ShadowEvent {
	r.mu.Lock()
	defer r.mu.Unlock()

	// deduplicate by port, keeping highest count
	byPort := make(map[int]ShadowEvent)
	for _, ev := range r.events {
		if existing, ok := byPort[ev.Port]; !ok || ev.Count > existing.Count {
			byPort[ev.Port] = ev
		}
	}

	result := make([]ShadowEvent, 0, len(byPort))
	for _, ev := range byPort {
		result = append(result, ev)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if n > 0 && len(result) > n {
		return result[:n]
	}
	return result
}

// Summary returns a formatted report of all shadow events.
func (r *PortShadowReporter) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.events) == 0 {
		return "portshadow_reporter: no events recorded"
	}

	lines := []string{fmt.Sprintf("portshadow_reporter: %d event(s)", len(r.events))}
	for _, ev := range r.events {
		lines = append(lines, fmt.Sprintf("  port=%d %s->%s count=%d",
			ev.Port, ev.KnownAddr, ev.NewAddr, ev.Count))
	}
	return strings.Join(lines, "\n")
}
