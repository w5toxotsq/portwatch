// Package history records and persists a log of port-change events detected
// by portwatch across daemon poll cycles.
//
// Each time the daemon detects that ports have been opened or closed, it calls
// History.Record to append a timestamped Entry. The history can be saved to
// disk with Save and reloaded across restarts with Load.
//
// Typical usage:
//
//	h, err := history.Load(cfg.HistoryPath)
//	if err != nil { ... }
//	h.Record(diff.Opened, diff.Closed)
//	h.Save()
package history
