package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRunScan_TextOutput(t *testing.T) {
	port, stop := startTestListener(t)
	defer stop()

	_ = port // port is open; scanOptions uses config so we test printScanText directly
	ports := []scanner.Port{
		{Port: port, Protocol: "tcp", Address: "127.0.0.1"},
	}
	out := captureStdout(t, func() {
		if err := printScanText(ports); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, strconv.Itoa(port)) {
		t.Errorf("expected port %d in output, got: %s", port, out)
	}
	if !strings.Contains(out, "tcp") {
		t.Errorf("expected protocol 'tcp' in output, got: %s", out)
	}
}

func TestRunScan_EmptyTextOutput(t *testing.T) {
	out := captureStdout(t, func() {
		if err := printScanText([]scanner.Port{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "No open ports") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestRunScan_JSONOutput(t *testing.T) {
	ports := []scanner.Port{
		{Port: 8080, Protocol: "tcp", Address: "127.0.0.1"},
	}
	out := captureStdout(t, func() {
		if err := printScanJSON(ports); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	var result []scanner.Port
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 1 || result[0].Port != 8080 {
		t.Errorf("unexpected result: %+v", result)
	}
}
