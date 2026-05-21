package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatAnnotateResult returns a human-readable summary of an AnnotateResult.
func FormatAnnotateResult(r AnnotateResult) string {
	var sb strings.Builder

	if len(r.Annotations) == 0 && len(r.Skipped) == 0 && len(r.Missing) == 0 {
		sb.WriteString("No annotations applied.\n")
		return sb.String()
	}

	if len(r.Annotations) > 0 {
		keys := make([]string, 0, len(r.Annotations))
		for k := range r.Annotations {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("Annotated (%d):\n", len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s  %s\n", k, r.Annotations[k]))
		}
	}

	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped (%d, already annotated):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("Missing (%d, not in env):\n", len(r.Missing)))
		for _, k := range r.Missing {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	return sb.String()
}
