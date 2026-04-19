package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportRecord is a JSON-serialisable snapshot of a single alert event.
type ExportRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	Port      int       `json:"port"`
	Action    string    `json:"action"`
	Level     string    `json:"level"`
}

// Exporter writes alert events to a newline-delimited JSON file.
type Exporter struct {
	path string
	f    *os.File
	enc  *json.Encoder
}

// NewExporter opens (or creates) the file at path for appending.
func NewExporter(path string) (*Exporter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("exporter: open %s: %w", path, err)
	}
	return &Exporter{path: path, f: f, enc: json.NewEncoder(f)}, nil
}

// Write appends a single record to the file.
func (e *Exporter) Write(r ExportRecord) error {
	if err := e.enc.Encode(r); err != nil {
		return fmt.Errorf("exporter: encode: %w", err)
	}
	return nil
}

// Close closes the underlying file.
func (e *Exporter) Close() error {
	return e.f.Close()
}
