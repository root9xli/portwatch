package monitor

import (
	"testing"
	"time"
)

func newTestRepeat() *PortRepeat {
	r := NewPortRepeat(5 * time.Minute)
	return r
}

func TestPortRepeatFirstObserveReturnsOne(t *testing.T) {
	r := newTestRepeat()
	count := r.Observe(8080)
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
}

func TestPortRepeatIncrementsOnRepeat(t *testing.T) {
	r := newTestRepeat()
	r.Observe(8080)
	r.Observe(8080)
	count := r.Observe(8080)
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestPortRepeatIndependentPorts(t *testing.T) {
	r := newTestRepeat()
	r.Observe(80)
	r.Observe(80)
	r.Observe(443)
	e80, _ := r.Get(80)
	e443, _ := r.Get(443)
	if e80.Count != 2 {
		t.Errorf("port 80: expected 2, got %d", e80.Count)
	}
	if e443.Count != 1 {
		t.Errorf("port 443: expected 1, got %d", e443.Count)
	}
}

func TestPortRepeatGetMissingReturnsFalse(t *testing.T) {
	r := newTestRepeat()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected false for missing port")
	}
}

func TestPortRepeatResetClearsCount(t *testing.T) {
	r := newTestRepeat()
	r.Observe(8080)
	r.Observe(8080)
	r.Reset(8080)
	_, ok := r.Get(8080)
	if ok {
		t.Fatal("expected entry to be removed after Reset")
	}
}

func TestPortRepeatEvictsOldEntries(t *testing.T) {
	r := NewPortRepeat(1 * time.Second)
	now := time.Now()
	r.now = func() time.Time { return now }
	r.Observe(8080)
	// advance time past maxAge
	r.now = func() time.Time { return now.Add(2 * time.Second) }
	r.Observe(9090) // triggers evict
	_, ok := r.Get(8080)
	if ok {
		t.Fatal("expected port 8080 to be evicted")
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Len())
	}
}

func TestPortRepeatLenTracksEntries(t *testing.T) {
	r := newTestRepeat()
	if r.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	r.Observe(80)
	r.Observe(443)
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}
