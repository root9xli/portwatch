package monitor

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExporterWritesRecord(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "export-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	ex, err := NewExporter(f.Name())
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}

	rec := ExportRecord{
		Timestamp: time.Now(),
		Hostname:  "testhost",
		Port:      8080,
		Action:    "added",
		Level:     "warn",
	}
	if err := ex.Write(rec); err != nil {
		t.Fatalf("Write: %v", err)
	}
	ex.Close()

	data, _ := os.ReadFile(f.Name())
	var got ExportRecord
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Port != 8080 {
		t.Errorf("port: want 8080, got %d", got.Port)
	}
	if got.Hostname != "testhost" {
		t.Errorf("hostname: want testhost, got %s", got.Hostname)
	}
}

func TestExporterAppendsMultiple(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/out.jsonl"

	for i := 0; i < 3; i++ {
		ex, _ := NewExporter(path)
		ex.Write(ExportRecord{Port: 9000 + i, Action: "added", Level: "info"})
		ex.Close()
	}

	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	_ = bufio.NewScanner(strings.NewReader(string(data)))
}

func TestNewExporterBadPath(t *testing.T) {
	_, err := NewExporter("/no/such/dir/file.jsonl")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}
