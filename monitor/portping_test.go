package monitor

import (
	"net"
	"testing"
	"time"
)

func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return port, func() { _ = ln.Close() }
}

func TestPortPingReachable(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	pp := NewPortPing(time.Second)
	res := pp.Probe(port)
	if !res.Reachable {
		t.Fatalf("expected reachable, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestPortPingUnreachable(t *testing.T) {
	pp := NewPortPing(100 * time.Millisecond)
	res := pp.Probe(1) // port 1 should be unreachable in test env
	if res.Reachable {
		t.Skip("port 1 unexpectedly reachable; skipping")
	}
	if res.Err == nil {
		t.Error("expected non-nil error for unreachable port")
	}
}

func TestPortPingProbeAll(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	pp := NewPortPing(time.Second)
	results := pp.ProbeAll([]int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Reachable {
		t.Errorf("first port should be reachable")
	}
}

func TestNewPortPingDefaultTimeout(t *testing.T) {
	pp := NewPortPing(0)
	if pp.timeout != 2*time.Second {
		t.Errorf("expected default 2s timeout, got %v", pp.timeout)
	}
}
