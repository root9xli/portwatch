package monitor

import (
	"fmt"
	"strings"
	"sync"
)

// DeadPortReporter accumulates recycled port events and produces a summary.
type DeadPortReporter struct {
	mu       sync.Mutex
	events   []deadPortEvent
	deadPort *DeadPort
}

type deadPortEvent struct {
	Port    uint16
	Message string
}

func NewDeadPortReporter(dp *DeadPort) *DeadPortReporter {
	return &DeadPortReporter{deadPort: dp}
}

// Observe checks added ports against the DeadPort tracker and records recycled ones.
func (r *DeadPortReporter) Observe(msgs []Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range msgs {
		if m.Action == "removed" {
			r.deadPort.MarkGone(m.Port)
		}
		if m.Action == "added" && r.deadPort.IsRecycled(m.Port) {
			r.events = append(r.events, deadPortEvent{
				Port:    m.Port,
				Message: fmt.Sprintf("port %d recycled within window", m.Port),
			})
		}
	}
}

// Summary returns a human-readable report of recycled port events.
func (r *DeadPortReporter) Summary() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.events) == 0 {
		return "no recycled ports detected"
	}
	lines := make([]string, 0, len(r.events))
	for _, e := range r.events {
		lines = append(lines, e.Message)
	}
	return strings.Join(lines, "\n")
}

// Count returns the number of recycled port events recorded.
func (r *DeadPortReporter) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}
