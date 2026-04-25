package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatNormalizeResult returns a human-readable summary of a NormalizeResult.
func FormatNormalizeResult(r NormalizeResult) string {
	if len(r.Changes) == 0 {
		return "No changes made during normalization.\n"
	}

	// Sort changes by (new) key for deterministic output.
	sorted := make([]NormalizeChange, len(r.Changes))
	copy(sorted, r.Changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Normalization applied %d change(s):\n", len(sorted)))

	for _, c := range sorted {
		if c.OldKey != "" && c.NewVal == "" {
			// removed
			sb.WriteString(fmt.Sprintf("  - REMOVED  %s  (%s)\n", c.OldKey, c.Reason))
			continue
		}
		if c.OldKey != "" {
			// renamed + possibly value changed
			sb.WriteString(fmt.Sprintf("  ~ RENAMED  %s -> %s  (%s)\n", c.OldKey, c.Key, c.Reason))
			continue
		}
		// value-only change
		sb.WriteString(fmt.Sprintf("  ~ CHANGED  %s: %q -> %q  (%s)\n", c.Key, c.OldVal, c.NewVal, c.Reason))
	}

	return sb.String()
}
