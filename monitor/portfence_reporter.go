package monitor

import (
	"fmt"
	"sort"
	"time"
)

// PortFenceReporter produces human-readable summaries of the current fence state.
type PortFenceReporter struct {
	fence *PortFence
}

// NewPortFenceReporter creates a reporter backed by the given PortFence.
func NewPortFenceReporter(fence *PortFence) *PortFenceReporter {
	return &PortFenceReporter{fence: fence}
}

// FenceRecord is a snapshot of a single fence entry for reporting.
type FenceRecord struct {
	Port      int
	Reason    string
	FencedAt  time.Time
	ExpiresAt time.Time
	Remaining time.Duration
}

// Records returns a sorted slice of active fence records.
func (r *PortFenceReporter) Records() []FenceRecord {
	r.fence.mu.Lock()
	defer r.fence.mu.Unlock()
	now := r.fence.clock()
	var out []FenceRecord
	for port, e := range r.fence.entries {
		if now.After(e.ExpiresAt) {
			delete(r.fence.entries, port)
			continue
		}
		out = append(out, FenceRecord{
			Port:      e.Port,
			Reason:    e.Reason,
			FencedAt:  e.FencedAt,
			ExpiresAt: e.ExpiresAt,
			Remaining: e.ExpiresAt.Sub(now).Round(time.Second),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })
	return out
}

// Summary returns a formatted multi-line report of all active fences.
func (r *PortFenceReporter) Summary() string {
	recs := r.Records()
	if len(recs) == 0 {
		return "portfence: no active fences"
	}
	s := fmt.Sprintf("portfence: %d active fence(s)\n", len(recs))
	for _, rec := range recs {
		s += fmt.Sprintf("  port=%-6d reason=%-20s remaining=%s\n",
			rec.Port, rec.Reason, rec.Remaining)
	}
	return s
}

// Count returns the number of currently active fences.
func (r *PortFenceReporter) Count() int {
	return len(r.Records())
}
