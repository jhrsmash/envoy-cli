package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// SummaryOptions controls what is included in the environment summary.
type SummaryOptions struct {
	// ShowValues includes values in the summary output (default: false for safety).
	ShowValues bool
	// RedactSensitive replaces sensitive values with [REDACTED].
	RedactSensitive bool
	// PrefixGroups groups keys by their prefix delimiter.
	PrefixGroups bool
	// Delimiter is the separator used to detect prefix groups (default: "_").
	Delimiter string
}

// SummaryEntry represents a single key entry in the summary.
type SummaryEntry struct {
	Key       string
	Value     string
	Sensitive bool
	Group     string
}

// SummaryResult holds the output of a Summarize call.
type SummaryResult struct {
	Entries    []SummaryEntry
	TotalKeys  int
	Sensitive  int
	EmptyCount int
	Groups     []string
}

// Summarize produces a structured summary of the provided env map.
func Summarize(env map[string]string, opts SummaryOptions) SummaryResult {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	groupSet := map[string]struct{}{}
	var entries []SummaryEntry
	emptyCount := 0
	sensitiveCount := 0

	for _, k := range keys {
		v := env[k]
		sens := isSensitive(k)

		if v == "" {
			emptyCount++
		}
		if sens {
			sensitiveCount++
		}

		display := v
		if !opts.ShowValues {
			display = "***"
		} else if opts.RedactSensitive && sens {
			display = "[REDACTED]"
		}

		group := ""
		if opts.PrefixGroups {
			if idx := strings.Index(k, opts.Delimiter); idx > 0 {
				group = k[:idx]
				groupSet[group] = struct{}{}
			}
		}

		entries = append(entries, SummaryEntry{
			Key:       k,
			Value:     display,
			Sensitive: sens,
			Group:     group,
		})
	}

	groups := make([]string, 0, len(groupSet))
	for g := range groupSet {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	return SummaryResult{
		Entries:    entries,
		TotalKeys:  len(keys),
		Sensitive:  sensitiveCount,
		EmptyCount: emptyCount,
		Groups:     groups,
	}
}

// FormatSummaryResult returns a human-readable representation of a SummaryResult.
func FormatSummaryResult(r SummaryResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total keys : %d\n", r.TotalKeys))
	sb.WriteString(fmt.Sprintf("Sensitive  : %d\n", r.Sensitive))
	sb.WriteString(fmt.Sprintf("Empty      : %d\n", r.EmptyCount))
	if len(r.Groups) > 0 {
		sb.WriteString(fmt.Sprintf("Groups     : %s\n", strings.Join(r.Groups, ", ")))
	}
	sb.WriteString("\n")
	for _, e := range r.Entries {
		line := fmt.Sprintf("  %-30s = %s", e.Key, e.Value)
		if e.Sensitive {
			line += "  [sensitive]"
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
