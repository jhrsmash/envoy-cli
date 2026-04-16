package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatMergeResult returns a human-readable summary of a MergeResult.
func FormatMergeResult(r MergeResult) string {
	var sb strings.Builder

	if len(r.Added) == 0 && len(r.Conflicts) == 0 {
		sb.WriteString("No changes during merge.\n")
		return sb.String()
	}

	if len(r.Added) > 0 {
		sort.Strings(r.Added)
		sb.WriteString(fmt.Sprintf("Added (%d):\n", len(r.Added)))
		for _, k := range r.Added {
			sb.WriteString(fmt.Sprintf("  + %s=%s\n", k, r.Merged[k]))
		}
	}

	if len(r.Conflicts) > 0 {
		sort.Strings(r.Conflicts)
		sb.WriteString(fmt.Sprintf("Conflicts (%d):\n", len(r.Conflicts)))
		for _, k := range r.Conflicts {
			sb.WriteString(fmt.Sprintf("  ~ %s => %s\n", k, r.Merged[k]))
		}
	}

	return sb.String()
}
