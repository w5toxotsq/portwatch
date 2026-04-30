// Package report provides functionality for generating formatted reports
// of portwatch history data.
//
// Reports can be rendered as human-readable text tables or as structured
// JSON output, making them suitable for both interactive CLI use and
// downstream tooling or automation.
//
// Usage:
//
//	opts := report.DefaultOptions()
//	opts.Format = report.FormatJSON
//	if err := report.Generate(h, opts); err != nil {
//	    log.Fatal(err)
//	}
package report
