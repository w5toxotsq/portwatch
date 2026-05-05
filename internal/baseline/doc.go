// Package baseline provides functionality for capturing and comparing
// a known-good set of open ports against the current state.
//
// A baseline represents an intentional snapshot of expected open ports.
// It can be saved to disk and later loaded to compare against a live scan,
// highlighting ports that have appeared or disappeared unexpectedly.
//
// Usage:
//
//	b := baseline.New(ports)
//	if err := b.Save(path); err != nil { ... }
//
//	loaded, err := baseline.Load(path)
//	if err != nil { ... }
//
//	unexpected, missing := loaded.Compare(currentPorts)
package baseline
