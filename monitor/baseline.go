package monitor

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// BaselineEntry records when a port was first seen.
type BaselineEntry struct {
	Port      int       `json:"port"`
	FirstSeen time.Time `json:"first_seen"`
}

// Baseline tracks the set of ports considered "known good".
type Baseline struct {
	mu      sync.RWMutex
	entries map[int]BaselineEntry
	path    string
}

// NewBaseline creates an empty Baseline backed by the given file path.
func NewBaseline(path string) *Baseline {
	return &Baseline{
		entries: make(map[int]BaselineEntry),
		path:    path,
	}
}

// Contains returns true if the port is in the baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[port]
	return ok
}

// Add adds a port to the baseline with the current timestamp.
func (b *Baseline) Add(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.entries[port]; !ok {
		b.entries[port] = BaselineEntry{Port: port, FirstSeen: time.Now()}
	}
}

// Len returns the number of baselined ports.
func (b *Baseline) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.entries)
}

// Save persists the baseline to disk.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	f, err := os.Create(b.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b.entries)
}

// Load reads a previously saved baseline from disk.
func (b *Baseline) Load() error {
	f, err := os.Open(b.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	b.mu.Lock()
	defer b.mu.Unlock()
	return json.NewDecoder(f).Decode(&b.entries)
}
