package monitor

import (
	"strings"
	"testing"
)

func TestPortGroupAddAndGet(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("web", []int{80, 443, 8080})

	ports, ok := pg.Get("web")
	if !ok {
		t.Fatal("expected group 'web' to exist")
	}
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}
}

func TestPortGroupGetMissing(t *testing.T) {
	pg := NewPortGroup()
	_, ok := pg.Get("nonexistent")
	if ok {
		t.Fatal("expected missing group to return false")
	}
}

func TestPortGroupDeduplicates(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("db", []int{5432, 5432, 3306, 3306})

	ports, _ := pg.Get("db")
	if len(ports) != 2 {
		t.Fatalf("expected 2 unique ports, got %d", len(ports))
	}
}

func TestPortGroupSorted(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("mixed", []int{9000, 80, 443})

	ports, _ := pg.Get("mixed")
	if ports[0] != 80 || ports[1] != 443 || ports[2] != 9000 {
		t.Fatalf("expected sorted ports, got %v", ports)
	}
}

func TestPortGroupContains(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("web", []int{80, 443})

	if !pg.Contains("web", 80) {
		t.Error("expected port 80 to be in 'web'")
	}
	if pg.Contains("web", 8080) {
		t.Error("expected port 8080 NOT to be in 'web'")
	}
	if pg.Contains("missing", 80) {
		t.Error("expected false for missing group")
	}
}

func TestPortGroupNames(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("web", []int{80})
	pg.Add("db", []int{5432})
	pg.Add("cache", []int{6379})

	names := pg.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	if names[0] != "cache" || names[1] != "db" || names[2] != "web" {
		t.Fatalf("expected sorted names, got %v", names)
	}
}

func TestPortGroupSummaryEmpty(t *testing.T) {
	pg := NewPortGroup()
	got := pg.Summary()
	if got != "no port groups defined" {
		t.Fatalf("unexpected summary: %q", got)
	}
}

func TestPortGroupSummaryContainsName(t *testing.T) {
	pg := NewPortGroup()
	pg.Add("web", []int{80, 443})

	got := pg.Summary()
	if !strings.Contains(got, "web") {
		t.Errorf("expected summary to contain group name, got: %s", got)
	}
	if !strings.Contains(got, "80") {
		t.Errorf("expected summary to contain port 80, got: %s", got)
	}
}
