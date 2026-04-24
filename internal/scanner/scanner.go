package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	Address  string
}

// Scan checks which ports in the given range are open on the given host.
// protocol must be "tcp" or "udp".
func Scan(host, protocol string, ports []int, timeout time.Duration) ([]PortState, error) {
	protocol = strings.ToLower(protocol)
	if protocol != "tcp" && protocol != "udp" {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}

	results := make([]PortState, 0, len(ports))
	for _, port := range ports {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		open := isOpen(addr, protocol, timeout)
		results = append(results, PortState{
			Port:     port,
			Protocol: protocol,
			Open:     open,
			Address:  addr,
		})
	}
	return results, nil
}

// isOpen attempts a connection to determine if a port is open.
func isOpen(addr, protocol string, timeout time.Duration) bool {
	conn, err := net.DialTimeout(protocol, addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// OpenPorts filters a slice of PortState and returns only the open ones.
func OpenPorts(states []PortState) []PortState {
	open := make([]PortState, 0)
	for _, s := range states {
		if s.Open {
			open = append(open, s)
		}
	}
	return open
}
