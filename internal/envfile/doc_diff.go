// Package envfile provides utilities for parsing, diffing, and managing
// .env files across environments.
//
// # Diff
//
// The diff sub-feature compares two env maps and reports which keys were
// added, removed, or changed.
//
//	// Compute a diff between two environments.
//	entries := envfile.Diff(base, override)
//
//	// Pretty-print the diff.
//	fmt.Println(envfile.FormatDiff(entries))
//
//	// Summarise the diff into counts and key lists.
//	summary := envfile.SummarizeDiff(entries)
//	fmt.Println(envfile.FormatDiffSummary(summary))
//
//	// Export the diff to JSON or unified-patch format.
//	out, err := envfile.ExportDiff(entries, envfile.DiffExportJSON)
//
//	// Write the diff directly to a file (format inferred from extension).
//	err = envfile.ExportDiffToFile(entries, "changes.json")
package envfile
