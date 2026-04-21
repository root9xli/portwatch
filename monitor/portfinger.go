package monitor

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortFinger attempts to identify the service/banner on a given port
// by reading the first bytes after connecting.
type PortFinger struct {
	mu      sync.Mutex
	results map[int]FingerResult
	timeout time.Duration
}

// FingerResult holds the banner/fingerprint data for a port.
type FingerResult struct {
	Port    int
	Banner  string
	Probed  time.Time
	Reached bool
}

// NewPortFinger creates a PortFinger with the given dial timeout.
func NewPortFinger(timeout time.Duration) *PortFinger {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &PortFinger{
		results: make(map[int]FingerResult),
		timeout: timeout,
	}
}

// Probe connects to localhost:port and captures up to 256 bytes of banner.
func (pf *PortFinger) Probe(port int) FingerResult {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	result := FingerResult{Port: port, Probed: time.Now()}

	conn, err := net.DialTimeout("tcp", addr, pf.timeout)
	if err != nil {
		pf.store(result)
		return result
	}
	defer conn.Close()

	result.Reached = true
	_ = conn.SetReadDeadline(time.Now().Add(pf.timeout))

	buf := make([]byte, 256)
	n, _ := conn.Read(buf)
	if n > 0 {
		result.Banner = sanitize(buf[:n])
	}

	pf.store(result)
	return result
}

// ProbeAll probes each port in the slice and returns all results.
func (pf *PortFinger) ProbeAll(ports []int) []FingerResult {
	out := make([]FingerResult, 0, len(ports))
	for _, p := range ports {
		out = append(out, pf.Probe(p))
	}
	return out
}

// Get returns the last stored result for a port, if any.
func (pf *PortFinger) Get(port int) (FingerResult, bool) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	r, ok := pf.results[port]
	return r, ok
}

func (pf *PortFinger) store(r FingerResult) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	pf.results[r.Port] = r
}

// sanitize replaces non-printable bytes with '.'.
func sanitize(b []byte) string {
	out := make([]byte, len(b))
	for i, c := range b {
		if c >= 32 && c < 127 {
			out[i] = c
		} else {
			out[i] = '.'
		}
	}
	return string(out)
}
