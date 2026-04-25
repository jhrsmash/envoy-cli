package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// DiffFilterOptions controls which diff entries are included in the result.
type DiffFilterOptions struct {
	// Statuses limits results to entries with the given statuses.
	// Valid values: "added", "removed", "changed", "unchanged".
	// If empty, all statuses are included.
	Statuses []string

	// KeyPrefix, if set, only includes entries whose key starts with the prefix.
	KeyPrefix string

	// ExcludeUnchanged is a convenience flag that drops unchanged entries.
	ExcludeUnchanged bool
}

// FilterDiff returns a subset of diff entries matching the given options.
func FilterDiff(entries []DiffEntry, opts DiffFilterOptions) ([]DiffEntry, error) {
	statusSet := make(map[string]bool)
	for _, s := range opts.Statuses {
		norm := strings.ToLower(strings.TrimSpace(s))
		switch norm {
		case "added", "removed", "changed", "unchanged":
			statusSet[norm] = true
		default:
			return nil, fmt.Errorf("invalid diff status filter %q: must be added, removed, changed, or unchanged", s)
		}
	}

	var result []DiffEntry
	for _, e := range entries {
		if opts.ExcludeUnchanged && e.Status == StatusUnchanged {
			continue
		}
		if len(statusSet) > 0 {
			if !statusSet[strings.ToLower(string(e.Status))] {
				continue
			}
		}
		if opts.KeyPrefix != "" && !strings.HasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result, nil
}

// FormatFilteredDiff formats a filtered slice of DiffEntry values for display.
func FormatFilteredDiff(entries []DiffEntry) string {
	if len(entries) == 0 {
		return "(no matching diff entries)\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.NewValue)
		case StatusRemoved:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.OldValue)
		case StatusChanged:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", e.Key, e.OldValue, e.NewValue)
		case StatusUnchanged:
			fmt.Fprintf(&sb, "  %s=%s\n", e.Key, e.NewValue)
		}
	}
	return sb.String()
}
