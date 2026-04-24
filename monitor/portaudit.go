package monitor

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// AuditEntry records a single audit event for a port.
type AuditEntry struct {
	Port      int
	Event     string
	Level     string
	Timestamp time.Time
}

// PortAudit maintains an append-only audit log of port events.
type PortAudit struct {
	mu      sync.Mutex
	entries []AuditEntry
	maxSize int
}

// NewPortAudit creates a PortAudit with the given maximum log size.
func NewPortAudit(maxSize int) *PortAudit {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &PortAudit{maxSize: maxSize}
}

// Record appends an audit event for the given port.
func (a *PortAudit) Record(port int, event, level string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, AuditEntry{
		Port:      port,
		Event:     event,
		Level:     level,
		Timestamp: time.Now(),
	})
}

// EntriesForPort returns all audit entries for a specific port.
func (a *PortAudit) EntriesForPort(port int) []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	var result []AuditEntry
	for _, e := range a.entries {
		if e.Port == port {
			result = append(result, e)
		}
	}
	return result
}

// Len returns the total number of audit entries.
func (a *PortAudit) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.entries)
}

// Summary returns a human-readable audit log summary.
func (a *PortAudit) Summary() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) == 0 {
		return "audit log: no entries"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("audit log: %d entries\n", len(a.entries)))
	for _, e := range a.entries {
		sb.WriteString(fmt.Sprintf("  [%s] port=%d event=%s level=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Port, e.Event, e.Level))
	}
	return strings.TrimRight(sb.String(), "\n")
}
