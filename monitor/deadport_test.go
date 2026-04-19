package monitor

import (
	"testing"
	"time"
)

func TestDeadPortNotRecycledIfNeverGone(t *testing.T) {
	d := NewDeadPort(5 * time.Second)
	if d.IsRecycled(8080) {
		t.Fatal("expected false for port never marked gone")
	}
}

func TestDeadPortRecycledWithinWindow(t *testing.T) {
	now := time.Now()
	d := NewDeadPort(5 * time.Second)
	d.clock = func() time.Time { return now }
	d.MarkGone(8080)
	d.clock = func() time.Time { return now.Add(2 * time.Second) }
	if !d.IsRecycled(8080) {
		t.Fatal("expected recycled within window")
	}
}

func TestDeadPortNotRecycledAfterWindow(t *testing.T) {
	now := time.Now()
	d := NewDeadPort(5 * time.Second)
	d.clock = func() time.Time { return now }
	d.MarkGone(8080)
	d.clock = func() time.Time { return now.Add(10 * time.Second) }
	if d.IsRecycled(8080) {
		t.Fatal("expected not recycled after window")
	}
}

func TestDeadPortIndependentPorts(t *testing.T) {
	now := time.Now()
	d := NewDeadPort(5 * time.Second)
	d.clock = func() time.Time { return now }
	d.MarkGone(9090)
	if d.IsRecycled(8080) {
		t.Fatal("expected false for different port")
	}
}

func TestDeadPortEvictRemovesExpired(t *testing.T) {
	now := time.Now()
	d := NewDeadPort(5 * time.Second)
	d.clock = func() time.Time { return now }
	d.MarkGone(8080)
	d.clock = func() time.Time { return now.Add(10 * time.Second) }
	d.Evict()
	if _, ok := d.seen[8080]; ok {
		t.Fatal("expected evicted entry to be removed")
	}
}
