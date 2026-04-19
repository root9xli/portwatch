package monitor

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHealthCheckServesOverHTTP(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	h := NewHealthCheck(addr)
	h.Start()
	defer h.Close()

	// wait for server
	var resp *http.Response
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		resp, err = http.Get(fmt.Sprintf("http://%s/healthz", addr))
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("could not reach healthz: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHealthCheckReadyzAfterCycle(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	h := NewHealthCheck(addr)
	h.Start()
	defer h.Close()

	time.Sleep(50 * time.Millisecond)
	h.RecordCycle()

	resp, err := http.Get(fmt.Sprintf("http://%s/readyz", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 after cycle, got %d", resp.StatusCode)
	}
	_ = strings.Contains 
}
