package envfile

import (
	"fmt"
	"strings"
)

// SnapshotDiffResult holds the result of comparing two snapshots.
type SnapshotDiffResult struct {
	From   string
	To     string
	Changes []DiffEntry
}

// DiffSnapshots compares two snapshots and returns a SnapshotDiffResult.
func DiffSnapshots(from, to Snapshot) SnapshotDiffResult {
	changes := Diff(from.Env, to.Env)
	return SnapshotDiffResult{
		From:    from.Label,
		To:      to.Label,
		Changes: changes,
	}
}

// FormatSnapshotDiff returns a human-readable string of a SnapshotDiffResult.
func FormatSnapshotDiff(r SnapshotDiffResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Snapshot diff: %s → %s\n", r.From, r.To)
	if len(r.Changes) == 0 {
		sb.WriteString("  (no changes)\n")
		return sb.String()
	}
	for _, e := range r.Changes {
		switch e.Type {
		case DiffAdded:
			fmt.Fprintf(&sb, "  + %s=%s\n", e.Key, e.NewValue)
		case DiffRemoved:
			fmt.Fprintf(&sb, "  - %s=%s\n", e.Key, e.OldValue)
		case DiffChanged:
			fmt.Fprintf(&sb, "  ~ %s: %s → %s\n", e.Key, e.OldValue, e.NewValue)
		}
	}
	return sb.String()
}
