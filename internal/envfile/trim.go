package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// TrimResult holds the outcome of a trim operation.
type TrimResult struct {
	Output   map[string]string
	Removed  []string // keys removed because their values were blank/whitespace
	Trimmed  []string // keys whose values had leading/trailing whitespace stripped
}

// TrimOptions controls the behaviour of Trim.
type TrimOptions struct {
	// RemoveEmpty removes keys whose value is empty (or all whitespace) after
	// trimming.
	RemoveEmpty bool

	// TrimValues strips leading and trailing whitespace from every value.
	TrimValues bool
}

// Trim cleans up an env map according to the supplied options.
// It never mutates the input map.
func Trim(env map[string]string, opts TrimOptions) TrimResult {
	output := make(map[string]string, len(env))
	var removed, trimmed []string

	for k, v := range env {
		newVal := v

		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
			if newVal != v {
				trimmed = append(trimmed, k)
			}
		}

		if opts.RemoveEmpty && strings.TrimSpace(newVal) == "" {
			removed = append(removed, k)
			continue
		}

		output[k] = newVal
	}

	sort.Strings(removed)
	sort.Strings(trimmed)

	return TrimResult{
		Output:  output,
		Removed: removed,
		Trimmed: trimmed,
	}
}

// FormatTrimResult returns a human-readable summary of a TrimResult.
func FormatTrimResult(r TrimResult) string {
	var sb strings.Builder

	if len(r.Trimmed) == 0 && len(r.Removed) == 0 {
		sb.WriteString("trim: nothing to do — all values are already clean\n")
		return sb.String()
	}

	if len(r.Trimmed) > 0 {
		sb.WriteString(fmt.Sprintf("trimmed whitespace from %d key(s):\n", len(r.Trimmed)))
		for _, k := range r.Trimmed {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", k))
		}
	}

	if len(r.Removed) > 0 {
		sb.WriteString(fmt.Sprintf("removed %d empty key(s):\n", len(r.Removed)))
		for _, k := range r.Removed {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}

	return sb.String()
}
