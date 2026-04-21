package monitor

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Whitelist holds a set of ports that are always allowed.
type Whitelist struct {
	ports map[int]struct{}
}

// NewWhitelist creates an empty Whitelist.
func NewWhitelist() *Whitelist {
	return &Whitelist{ports: make(map[int]struct{})}
}

// Add adds a port to the whitelist.
func (w *Whitelist) Add(port int) {
	w.ports[port] = struct{}{}
}

// Remove removes a port from the whitelist. It is a no-op if the port is not present.
func (w *Whitelist) Remove(port int) {
	delete(w.ports, port)
}

// Contains reports whether port is whitelisted.
func (w *Whitelist) Contains(port int) bool {
	_, ok := w.ports[port]
	return ok
}

// Len returns the number of whitelisted ports.
func (w *Whitelist) Len() int {
	return len(w.ports)
}

// Ports returns a slice of all whitelisted port numbers in no particular order.
func (w *Whitelist) Ports() []int {
	ports := make([]int, 0, len(w.ports))
	for p := range w.ports {
		ports = append(ports, p)
	}
	return ports
}

// LoadWhitelistFile reads a file of newline-separated port numbers and
// returns a populated Whitelist. Lines starting with '#' are ignored.
func LoadWhitelistFile(path string) (*Whitelist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := NewWhitelist()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		port, err := strconv.Atoi(line)
		if err != nil {
			continue // skip malformed lines
		}
		w.Add(port)
	}
	return w, scanner.Err()
}
