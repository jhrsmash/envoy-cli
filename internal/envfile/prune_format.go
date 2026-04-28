package envfile

import (
	"fmt"
	"strings"
)

// FormatPruneResult returns a human-readable summary of a PruneResult.
func FormatPruneResult(r PruneResult) string {
	var sb strings.Builder

	if len(r.Removed) == 0 {
		sb.WriteString("No keys pruned.\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Pruned %d key(s):\n", len(r.Removed)))
	for _, k := range r.Removed {
		sb.WriteString(fmt.Sprintf("  - %s\n", k))
	}
	sb.WriteString(fmt.Sprintf("Remaining keys: %d\n", len(r.Output)))
	return sb.String()
}
