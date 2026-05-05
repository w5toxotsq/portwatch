// Package baseline manages a user-approved set of open ports that serves as
// the reference point for anomaly detection.
//
// A baseline is created by capturing the current port state and saving it to
// disk. Subsequent scans are compared against the baseline to surface
// unexpected new ports (ports present in the scan but absent from the
// baseline) and missing ports (ports present in the baseline but no longer
// open).
//
// Usage:
//
//	b := baseline.New(ports)
//	_ = baseline.Save(b, "/var/lib/portwatch/baseline.json")
//
//	loaded, err := baseline.Load("/var/lib/portwatch/baseline.json")
//	unexpected, missing := baseline.Compare(loaded, currentPorts)
package baseline
