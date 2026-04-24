package monitor

import (
	"testing"
)

func TestPortChainObserveSinglePortIgnored(t *testing.T) {
	pc := NewPortChain(16)
	pc.Observe([]int{8080})
	if pc.Len() != 0 {
		t.Fatalf("expected 0 entries for single-port observation, got %d", pc.Len())
	}
}

func TestPortChainObserveTwoPorts(t *testing.T) {
	pc := NewPortChain(16)
	pc.Observe([]int{80, 443})
	if pc.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", pc.Len())
	}
	e, ok := pc.Get([]int{80, 443})
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
}

func TestPortChainIncrementOnRepeat(t *testing.T) {
	pc := NewPortChain(16)
	pc.Observe([]int{80, 443})
	pc.Observe([]int{443, 80}) // order should not matter
	e, ok := pc.Get([]int{80, 443})
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Count != 2 {
		t.Fatalf("expected count 2, got %d", e.Count)
	}
}

func TestPortChainGetMissingReturnsFalse(t *testing.T) {
	pc := NewPortChain(16)
	_, ok := pc.Get([]int{22, 8080})
	if ok {
		t.Fatal("expected false for missing entry")
	}
}

func TestPortChainEvictsWhenFull(t *testing.T) {
	pc := NewPortChain(2)
	pc.Observe([]int{1, 2})
	pc.Observe([]int{3, 4})
	// Both slots full; observing a third group should evict the lowest-count one.
	pc.Observe([]int{5, 6})
	if pc.Len() != 2 {
		t.Fatalf("expected 2 entries after eviction, got %d", pc.Len())
	}
}

func TestPortChainIndependentGroups(t *testing.T) {
	pc := NewPortChain(16)
	pc.Observe([]int{80, 443})
	pc.Observe([]int{22, 8080})
	if pc.Len() != 2 {
		t.Fatalf("expected 2 independent entries, got %d", pc.Len())
	}
	e1, _ := pc.Get([]int{80, 443})
	e2, _ := pc.Get([]int{22, 8080})
	if e1.Count != 1 || e2.Count != 1 {
		t.Fatalf("counts should both be 1, got %d and %d", e1.Count, e2.Count)
	}
}

func TestPortChainSortedCopy(t *testing.T) {
	out := sortedCopy([]int{9, 3, 7, 1})
	expected := []int{1, 3, 7, 9}
	for i, v := range expected {
		if out[i] != v {
			t.Fatalf("sortedCopy mismatch at %d: want %d got %d", i, v, out[i])
		}
	}
}
