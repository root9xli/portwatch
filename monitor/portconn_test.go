package monitor

import (
	"net"
	"testing"
	"time"
)

func startTCPListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestPortConnProbeReachable(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	pc := NewPortConn(200 * time.Millisecond)
	if err := pc.Probe(port); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	rec, ok := pc.Get(port)
	if !ok {
		t.Fatal("expected record to exist")
	}
	if rec.Count != 1 {
		t.Errorf("expected count 1, got %d", rec.Count)
	}
}

func TestPortConnProbeUnreachable(t *testing.T) {
	pc := NewPortConn(100 * time.Millisecond)
	err := pc.Probe(1)
	if err == nil {
		t.Fatal("expected error for unreachable port")
	}
	_, ok := pc.Get(1)
	if ok {
		t.Error("expected no record for failed probe")
	}
}

func TestPortConnReset(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	pc := NewPortConn(200 * time.Millisecond)
	_ = pc.Probe(port)
	pc.Reset(port)
	_, ok := pc.Get(port)
	if ok {
		t.Error("expected record to be cleared after reset")
	}
}

func TestPortConnProbeAllReturnsErrors(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	pc := NewPortConn(100 * time.Millisecond)
	errs := pc.ProbeAll([]int{port, 2})
	if _, bad := errs[port]; bad {
		t.Errorf("reachable port %d should not have an error", port)
	}
	if _, bad := errs[2]; !bad {
		t.Error("unreachable port 2 should have an error")
	}
}

func TestPortConnCountIncrements(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	pc := NewPortConn(200 * time.Millisecond)
	for i := 0; i < 3; i++ {
		_ = pc.Probe(port)
	}
	rec, ok := pc.Get(port)
	if !ok {
		t.Fatal("expected record")
	}
	if rec.Count != 3 {
		t.Errorf("expected count 3, got %d", rec.Count)
	}
}
