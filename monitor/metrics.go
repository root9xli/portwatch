package monitor

import (
	"sync"
	"time"
)

// Metrics tracks runtime counters for the monitor daemon.
type Metrics struct {
	mu           sync.Mutex
	AlertsTotal  int
	DiffsTotal   int
	Suppressed   int
	RateLimited  int
	StartTime    time.Time
	lastDiff     time.Time
}

// NewMetrics creates a new Metrics instance with StartTime set to now.
func NewMetrics() *Metrics {
	return &Metrics{StartTime: time.Now()}
}

// RecordDiff increments the diff counter and records the time.
func (m *Metrics) RecordDiff() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DiffsTotal++
	m.lastDiff = time.Now()
}

// RecordAlert increments the alert counter.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AlertsTotal++
}

// RecordSuppressed increments the suppressed counter.
func (m *Metrics) RecordSuppressed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Suppressed++
}

// RecordRateLimited increments the rate-limited counter.
func (m *Metrics) RecordRateLimited() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RateLimited++
}

// Snapshot returns a copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Metrics{
		AlertsTotal: m.AlertsTotal,
		DiffsTotal:  m.DiffsTotal,
		Suppressed:  m.Suppressed,
		RateLimited: m.RateLimited,
		StartTime:   m.StartTime,
		lastDiff:    m.lastDiff,
	}
}

// Uptime returns the duration since the monitor started.
func (m *Metrics) Uptime() time.Duration {
	return time.Since(m.StartTime)
}

// LastDiff returns the time of the most recent diff.
func (m *Metrics) LastDiff() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastDiff
}
