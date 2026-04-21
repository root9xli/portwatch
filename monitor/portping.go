package monitor

import (
	"net"
	"strconv"
	"time"
)

// PortPing attempts a TCP dial to a port and reports whether it is reachable.
type PortPing struct {
	timeout time.Duration
}

// PingResult holds the outcome of a single port probe.
type PingResult struct {
	Port      int
	Reachable bool
	Latency   time.Duration
	Err       error
}

// NewPortPing creates a PortPing with the given dial timeout.
func NewPortPing(timeout time.Duration) *PortPing {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &PortPing{timeout: timeout}
}

// Probe dials localhost on the given port and returns a PingResult.
func (p *PortPing) Probe(port int) PingResult {
	addr := "127.0.0.1:" + strconv.Itoa(port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return PingResult{Port: port, Reachable: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return PingResult{Port: port, Reachable: true, Latency: latency}
}

// ProbeAll probes each port in the slice and returns all results.
func (p *PortPing) ProbeAll(ports []int) []PingResult {
	results := make([]PingResult, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Probe(port))
	}
	return results
}
