package monitor

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// PortBurstReporter summarises burst events observed during monitoring.
type PortBurstReporter struct {
	burst   *PortBurst
	events  []BurstEvent
}

// NewPortBurstReporter creates a reporter backed by the given PortBurst.
func NewPortBurstReporter(pb *PortBurst) *PortBurstReporter {
	return &PortBurstReporter{burst: pb}
}

// Record feeds a diff message through the burst detector and stores any
// resulting burst events.
func (r *PortBurstReporter) Record(msgs []Message) {
	now := time.Now()
	for _, m := range msgs {
		if ev := r.burst.Observe(m.Port, now); ev != nil {
			r.events = append(r.events, *ev)
		}
	}
}

// Events returns a copy of all recorded burst events.
func (r *PortBurstReporter) Events() []BurstEvent {
	out := make([]BurstEvent, len(r.events))
	copy(out, r.events)
	return out
}

// Summary returns a human-readable report of burst events grouped by port.
func (r *PortBurstReporter) Summary() string {
	if len(r.events) == 0 {
		return "portburst: no burst events recorded"
	}

	// aggregate max count per port
	peak := make(map[int]int)
	for _, ev := range r.events {
		if ev.Count > peak[ev.Port] {
			peak[ev.Port] = ev.Count
		}
	}

	ports := make([]int, 0, len(peak))
	for p := range peak {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	var sb strings.Builder
	sb.WriteString("portburst: burst events detected\n")
	for _, p := range ports {
		sb.WriteString(fmt.Sprintf("  port %d: peak burst count %d\n", p, peak[p]))
	}
	return strings.TrimRight(sb.String(), "\n")
}
