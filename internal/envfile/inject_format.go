package envfile

import (
	"fmt"
	"strings"
)

// FormatInjectResult returns a human-readable summary of an InjectResult.
func FormatInjectResult(r InjectResult) string {
	var sb strings.Builder

	if len(r.Injected) == 0 && len(r.Overwritten) == 0 && len(r.Skipped) == 0 {
		sb.WriteString("No variables injected.\n")
		return sb.String()
	}

	if len(r.Injected) > 0 {
		sb.WriteString(fmt.Sprintf("Injected (%d):\n", len(r.Injected)))
		for _, k := range r.Injected {
			sb.WriteString(fmt.Sprintf("  + %s\n", k))
		}
	}

	if len(r.Overwritten) > 0 {
		sb.WriteString(fmt.Sprintf("Overwritten (%d):\n", len(r.Overwritten)))
		for _, k := range r.Overwritten {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}

	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped (%d, already set):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return sb.String()
}
