package monitor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single listening port detected on the system.
type PortEntry struct {
	Port     int
	Protocol string
	PID  string
}

// Snapshot holds point-in-time collection of listening ports.
type Snapshot struct {
	Ports []PortEntry
}

// TakeSnapshot reads /proc/net/tcp and /proc/net/tcp6 to collect listening ports.
func TakeSnapshot() (*Snapshot, error) {
	var entries []PortEntry
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		result, err := parseNetTCP(path)
		if err != nil {
			continue // file may not exist on all systems
		}
		entries = append(entries, result...)
	}
	return &Snapshot{Ports: entries}, nil
}

// parseNetTCP parses a /proc/net/tcp style file and returns listening entries.
func parseNetTCP(path string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A == TCP_LISTEN
		if fields[3] != "0A" {
			continue
		}
		port, err := parseHexPort(fields[1])
		if err != nil {
			continue
		}
		entries = append(entries, PortEntry{Port: port, Protocol: "tcp"})
	}
	return entries, scanner.Err()
}

// parseHexPort extracts the port from a hex address field like "0F:1F90".
func parseHexPort(addr string) (int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid addr: %s", addr)
	}
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, err
	}
	return int(port), nil
}
