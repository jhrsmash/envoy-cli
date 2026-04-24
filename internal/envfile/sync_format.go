package envfile

import (
	"fmt"
	"strings"
)

// FormatSyncResult returns a human-readable summary of a SyncResult.
func FormatSyncResult(r SyncResult) string {
	if len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Updated) == 0 {
		return "sync: no changes"
	}

	var sb strings.Builder

	if len(r.Added) > 0 {
		sb.WriteString(fmt.Sprintf("added (%d):\n", len(r.Added)))
		for _, k := range r.Added {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	if len(r.Updated) > 0 {
		sb.WriteString(fmt.Sprintf("updated (%d):\n", len(r.Updated)))
		for _, k := range r.Updated {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}

	if len(r.Removed) > 0 {
		sb.WriteString(fmt.Sprintf("removed (%d):\n", len(r.Removed)))
		for _, k := range r.Removed {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}
