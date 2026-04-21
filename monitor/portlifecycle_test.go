package monitor

import (
	"testing"
	"time"
)

func TestPortLifecycleObserveFirstTime(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(8080)
	entry, ok := pl.Get(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if entry.SeenCount != 1 {
		t.Errorf("expected SeenCount=1, got %d", entry.SeenCount)
	}
	if entry.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", entry.Port)
	}
}

func TestPortLifecycleObserveIncrementsCount(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(9000)
	pl.Observe(9000)
	pl.Observe(9000)
	entry, ok := pl.Get(9000)
	if !ok {
		t.Fatal("expected entry for port 9000")
	}
	if entry.SeenCount != 3 {
		t.Errorf("expected SeenCount=3, got %d", entry.SeenCount)
	}
}

func TestPortLifecycleUptimeGrowsOverTime(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(443)
	time.Sleep(10 * time.Millisecond)
	pl.Observe(443)
	entry, _ := pl.Get(443)
	if entry.Uptime <= 0 {
		t.Errorf("expected uptime > 0, got %v", entry.Uptime)
	}
}

func TestPortLifecycleForgetRemovesEntry(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(22)
	pl.Forget(22)
	_, ok := pl.Get(22)
	if ok {
		t.Error("expected entry to be removed after Forget")
	}
}

func TestPortLifecycleGetMissingReturnsFalse(t *testing.T) {
	pl := NewPortLifecycle()
	_, ok := pl.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestPortLifecycleAllReturnsAllPorts(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(80)
	pl.Observe(443)
	pl.Observe(8080)
	entries := pl.All()
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestPortLifecycleIndependentPorts(t *testing.T) {
	pl := NewPortLifecycle()
	pl.Observe(80)
	pl.Observe(80)
	pl.Observe(443)
	e80, _ := pl.Get(80)
	e443, _ := pl.Get(443)
	if e80.SeenCount != 2 {
		t.Errorf("expected port 80 SeenCount=2, got %d", e80.SeenCount)
	}
	if e443.SeenCount != 1 {
		t.Errorf("expected port 443 SeenCount=1, got %d", e443.SeenCount)
	}
}
