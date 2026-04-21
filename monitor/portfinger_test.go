package monitor

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func startBannerListener(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_, _ = conn.Write([]byte(banner))
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestPortFingerReachable(t *testing.T) {
	port := startBannerListener(t, "SSH-2.0-OpenSSH")
	pf := NewPortFinger(time.Second)
	r := pf.Probe(port)
	if !r.Reached {
		t.Fatal("expected port to be reached")
	}
	if r.Banner != "SSH-2.0-OpenSSH" {
		t.Fatalf("unexpected banner: %q", r.Banner)
	}
}

func TestPortFingerUnreachable(t *testing.T) {
	pf := NewPortFinger(100 * time.Millisecond)
	r := pf.Probe(1) // port 1 should be unreachable in test env
	if r.Reached {
		t.Skip("port 1 unexpectedly reachable")
	}
	if r.Banner != "" {
		t.Fatalf("expected empty banner, got %q", r.Banner)
	}
}

func TestPortFingerStoresResult(t *testing.T) {
	port := startBannerListener(t, "HELLO")
	pf := NewPortFinger(time.Second)
	pf.Probe(port)
	r, ok := pf.Get(port)
	if !ok {
		t.Fatal("expected stored result")
	}
	if r.Banner != "HELLO" {
		t.Fatalf("unexpected banner: %q", r.Banner)
	}
}

func TestPortFingerProbeAll(t *testing.T) {
	p1 := startBannerListener(t, "SVC-A")
	p2 := startBannerListener(t, "SVC-B")
	pf := NewPortFinger(time.Second)
	results := pf.ProbeAll([]int{p1, p2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Reached {
			t.Errorf("port %d not reached", r.Port)
		}
	}
}

func TestNewPortFingerDefaultTimeout(t *testing.T) {
	pf := NewPortFinger(0)
	if pf.timeout != 2*time.Second {
		t.Fatalf("expected default 2s timeout, got %v", pf.timeout)
	}
	_ = fmt.Sprintf("%v", pf.timeout) // use fmt
}
