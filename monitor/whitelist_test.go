package monitor

import (
	"os"
	"testing"
)

func TestWhitelistAddAndContains(t *testing.T) {
	w := NewWhitelist()
	w.Add(80)
	w.Add(443)

	if !w.Contains(80) {
		t.Error("expected 80 to be whitelisted")
	}
	if !w.Contains(443) {
		t.Error("expected 443 to be whitelisted")
	}
	if w.Contains(8080) {
		t.Error("expected 8080 not to be whitelisted")
	}
}

func TestWhitelistLen(t *testing.T) {
	w := NewWhitelist()
	if w.Len() != 0 {
		t.Fatalf("expected 0, got %d", w.Len())
	}
	w.Add(22)
	w.Add(22) // duplicate
	if w.Len() != 1 {
		t.Fatalf("expected 1 after duplicate add, got %d", w.Len())
	}
}

func TestLoadWhitelistFile(t *testing.T) {
	content := "# comment\n22\n80\n\n443\nbad\n"
	f, err := os.CreateTemp("", "whitelist-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	w, err := LoadWhitelistFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Len() != 3 {
		t.Fatalf("expected 3 ports, got %d", w.Len())
	}
	for _, p := range []int{22, 80, 443} {
		if !w.Contains(p) {
			t.Errorf("expected port %d to be whitelisted", p)
		}
	}
}

func TestLoadWhitelistFileMissing(t *testing.T) {
	_, err := LoadWhitelistFile("/nonexistent/path/whitelist.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
