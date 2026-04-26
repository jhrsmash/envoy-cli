package envfile

import (
	"fmt"
	"strings"
)

// FormatScopeResult returns a human-readable summary of a ScopeResult.
func FormatScopeResult(r ScopeResult) string {
	var sb strings.Builder

	if len(r.Included) == 0 && len(r.Excluded) == 0 {
		sb.WriteString("scope: no keys in environment\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("scope: %d included, %d excluded\n",
		len(r.Included), len(r.Excluded)))

	if len(r.Included) > 0 {
		sb.WriteString("\nIncluded keys:\n")
		for _, k := range r.Included {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	if len(r.Excluded) > 0 {
		sb.WriteString("\nExcluded keys:\n")
		for _, k := range r.Excluded {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return sb.String()
}
