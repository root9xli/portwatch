package monitor

import (
	"sync"
	"time"
)

// TrendDirection indicates whether a port's activity is increasing, stable, or decreasing.
type TrendDirection string

const (
	TrendUp     TrendDirection = "up"
	TrendStable TrendDirection = "stable"
	TrendDown   TrendDirection = "down"
)

// portSample records a count observation at a point in time.
type portSample struct {
	at    time.Time
	count int
}

// PortTrend tracks how frequently each port appears across diff cycles
// and reports whether activity is trending up, stable, or down.
type PortTrend struct {
	mu      sync.Mutex
	window  time.Duration
	samples map[uint16][]portSample
	now     func() time.Time
}

// NewPortTrend creates a PortTrend that retains samples within the given window.
func NewPortTrend(window time.Duration) *PortTrend {
	return &PortTrend{
		window:  window,
		samples: make(map[uint16][]portSample),
		now:     time.Now,
	}
}

// Record adds an observation of count for port at the current time.
func (pt *PortTrend) Record(port uint16, count int) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	now := pt.now()
	pt.samples[port] = append(pt.samples[port], portSample{at: now, count: count})
	pt.evict(port, now)
}

// Trend returns the TrendDirection for a port based on retained samples.
// It compares the average of the first half of samples to the second half.
func (pt *PortTrend) Trend(port uint16) TrendDirection {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.evict(port, pt.now())
	samples := pt.samples[port]
	if len(samples) < 2 {
		return TrendStable
	}
	mid := len(samples) / 2
	early := avg(samples[:mid])
	late := avg(samples[mid:])
	switch {
	case late > early:
		return TrendUp
	case late < early:
		return TrendDown
	default:
		return TrendStable
	}
}

// evict removes samples outside the retention window. Caller must hold mu.
func (pt *PortTrend) evict(port uint16, now time.Time) {
	cutoff := now.Add(-pt.window)
	ss := pt.samples[port]
	i := 0
	for i < len(ss) && ss[i].at.Before(cutoff) {
		i++
	}
	pt.samples[port] = ss[i:]
}

func avg(ss []portSample) float64 {
	if len(ss) == 0 {
		return 0
	}
	sum := 0
	for _, s := range ss {
		sum += s.count
	}
	return float64(sum) / float64(len(ss))
}
