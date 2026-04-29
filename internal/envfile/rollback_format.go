package envfile

import (
	"fmt"
	"strings"
)

// FormatRollbackResult returns a human-readable summary of a rollback.
func FormatRollbackResult(r RollbackResult) string {
	var sb strings.Builder

	mode := "applied"
	if r.DryRun {
		mode = "dry-run"
	}

	sb.WriteString(fmt.Sprintf("Rollback [%s] to archive %q (%s)\n",
		mode, r.Label, r.Timestamp.Format("2006-01-02 15:04:05")))

	added, removed, changed, unchanged := 0, 0, 0, 0
	for _, e := range r.Diff {
		switch e.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusChanged:
			changed++
		case StatusUnchanged:
			unchanged++
		}
	}

	sb.WriteString(fmt.Sprintf("  added: %d  removed: %d  changed: %d  unchanged: %d\n",
		added, removed, changed, unchanged))

	if len(r.Diff) > 0 {
		sb.WriteString("\nChanges:\n")
		sb.WriteString(FormatDiff(r.Diff))
	}

	if r.DryRun {
		sb.WriteString("\n(no changes written — dry-run mode)\n")
	}

	return sb.String()
}
