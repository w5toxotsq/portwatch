// Package snapshot manages point-in-time records of open ports and provides
// utilities to persist, load, and compare them.
//
// A typical workflow:
//
//  1. Scan open ports using the scanner package.
//  2. Create a new Snapshot with [New].
//  3. Persist it to disk with [Save].
//  4. On the next scan cycle, load the previous snapshot with [Load].
//  5. Compare old and new snapshots with [Compare] to obtain a [Diff].
//  6. Use [Diff.HasChanges] to decide whether to emit an alert.
//
// Example:
//
//	prev, _ := snapshot.Load("last.json")
//	curr := snapshot.New(currentPorts)
//	diff := snapshot.Compare(prev, curr)
//	if diff.HasChanges() {
//		// alert: diff.Opened, diff.Closed
//	}
//	snapshot.Save(curr, "last.json")
package snapshot
