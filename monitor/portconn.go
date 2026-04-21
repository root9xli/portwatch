package monitor

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ConnRecord holds connection count and last-seen time for a port.
type ConnRecord struct {
	Port      int
	Count     int
	LastSeen  time.Time
}

// PortConn tracks active TCP connection counts per port.
type PortConn struct {
	mu      sync.Mutex
	records map[int]*ConnRecord
	timeout time.Duration
}

// NewPortConn creates a PortConn with the given dial timeout.
func NewPortConn(timeout time.Duration) *PortConn {
	return &PortConn{
		records: make(map[int]*ConnRecord),
		timeout: timeout,
	}
}

// Probe attempts a TCP connection to the given port on localhost and records the result.
func (pc *PortConn) Probe(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, pc.timeout)

	pc.mu.Lock()
	defer pc.mu.Unlock()

	rec, ok := pc.records[port]
	if !ok {
		rec = &ConnRecord{Port: port}
		pc.records[port] = rec
	}
	if err == nil {
		conn.Close()
		rec.Count++
		rec.LastSeen = time.Now()
	}
	return err
}

// Get returns the ConnRecord for a port, or false if not found.
func (pc *PortConn) Get(port int) (ConnRecord, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	rec, ok := pc.records[port]
	if !ok {
		return ConnRecord{}, false
	}
	return *rec, true
}

// Reset clears the connection record for a port.
func (pc *PortConn) Reset(port int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.records, port)
}

// ProbeAll probes each port in the provided list and returns any errors keyed by port.
func (pc *PortConn) ProbeAll(ports []int) map[int]error {
	errs := make(map[int]error)
	for _, p := range ports {
		if err := pc.Probe(p); err != nil {
			errs[p] = err
		}
	}
	return errs
}
