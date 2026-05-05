package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/config"
)

func startBaselineListener(t *testing.T) (port int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestRunBaseline_TextOutput(t *testing.T) {
	port := startBaselineListener(t)
	dir := t.TempDir()
	cfg := &config.Config{
		Ports:        []int{port},
		Protocol:     "tcp",
		Timeout:      200 * time.Millisecond,
		BaselinePath: filepath.Join(dir, "baseline.json"),
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runBaseline(cfg, "text")
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("runBaseline: %v", err)
	}
	if !strings.Contains(buf.String(), "Baseline saved") {
		t.Errorf("expected 'Baseline saved' in output, got: %s", buf.String())
	}
}

func TestRunBaseline_JSONOutput(t *testing.T) {
	port := startBaselineListener(t)
	dir := t.TempDir()
	cfg := &config.Config{
		Ports:        []int{port},
		Protocol:     "tcp",
		Timeout:      200 * time.Millisecond,
		BaselinePath: filepath.Join(dir, "baseline.json"),
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runBaseline(cfg, "json")
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("runBaseline: %v", err)
	}
	var b baseline.Baseline
	if err := json.Unmarshal(buf.Bytes(), &b); err != nil {
		t.Fatalf("output is not valid JSON: %v — output: %s", err, buf.String())
	}
}
