package scanner

import (
	"net"
	"testing"
	"time"
)

// startTCPServer starts a local TCP listener and returns its port and a stop function.
func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	states, err := Scan("127.0.0.1", "tcp", []int{port}, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 result, got %d", len(states))
	}
	if !states[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	states, err := Scan("127.0.0.1", "tcp", []int{1}, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if states[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_InvalidProtocol(t *testing.T) {
	_, err := Scan("127.0.0.1", "icmp", []int{80}, 200*time.Millisecond)
	if err == nil {
		t.Error("expected error for unsupported protocol, got nil")
	}
}

func TestOpenPorts_Filter(t *testing.T) {
	states := []PortState{
		{Port: 80, Open: true},
		{Port: 81, Open: false},
		{Port: 443, Open: true},
	}
	open := OpenPorts(states)
	if len(open) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(open))
	}
	for _, s := range open {
		if !s.Open {
			t.Errorf("OpenPorts returned a closed port: %d", s.Port)
		}
	}
}
