package monitor

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// PortFlapReporter accumulates flap events and produces a human-readable summary.
type PortFlapReporter struct {
	mu     sync.Mutex
	events map[int]*PortFlapEvent // last event per port
}

// NewPortFlapReporter creates a PortFlapReporter.
func NewPortFlapReporter() *PortFlapReporter {
	return &PortFlapReporter{
		events: make(map[int]*PortFlapEvent),
	}
}

// Record stores a flap event.
func (r *PortFlapReporter) Record(ev *PortFlapEvent) {
	if ev == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[ev.Port] = ev
}

// Summary returns a multi-line string describing all recorded flap events.
func (r *PortFlapReporter) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.events) == 0 {
		return "portflap: no flapping ports detected"
	}

	ports := make([]int, 0, len(r.events))
	for p := range r.events {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	var sb strings.Builder
	sb.WriteString("portflap: flapping ports\n")
	for _, p := range ports {
		ev := r.events[p]
		sb.WriteString(fmt.Sprintf("  port=%d flaps=%d window=%s\n",
			ev.Port, ev.Flaps, ev.Window))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Len returns the number of distinct flapping ports recorded.
func (r *PortFlapReporter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}
