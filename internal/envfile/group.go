package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds the outcome of grouping env keys by a common prefix delimiter.
type GroupResult struct {
	// Groups maps prefix -> (key -> value) for keys that matched a prefix.
	Groups map[string]map[string]string
	// Ungrouped contains keys that did not match any prefix group.
	Ungrouped map[string]string
}

// GroupOptions controls how Group behaves.
type GroupOptions struct {
	// Delimiter separates the prefix from the rest of the key (default "_").
	Delimiter string
	// MinGroupSize is the minimum number of keys required to form a group.
	// Groups smaller than this are placed into Ungrouped. Default is 1.
	MinGroupSize int
}

// Group partitions env into named groups based on key prefix before the delimiter.
// Keys with no delimiter, or belonging to groups smaller than MinGroupSize,
// are placed in Ungrouped.
func Group(env map[string]string, opts GroupOptions) GroupResult {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}
	if opts.MinGroupSize < 1 {
		opts.MinGroupSize = 1
	}

	candidates := make(map[string]map[string]string)
	ungrouped := make(map[string]string)

	for k, v := range env {
		idx := strings.Index(k, opts.Delimiter)
		if idx <= 0 {
			ungrouped[k] = v
			continue
		}
		prefix := k[:idx]
		if candidates[prefix] == nil {
			candidates[prefix] = make(map[string]string)
		}
		candidates[prefix][k] = v
	}

	groups := make(map[string]map[string]string)
	for prefix, members := range candidates {
		if len(members) < opts.MinGroupSize {
			for k, v := range members {
				ungrouped[k] = v
			}
		} else {
			groups[prefix] = members
		}
	}

	return GroupResult{Groups: groups, Ungrouped: ungrouped}
}

// FormatGroupResult returns a human-readable summary of a GroupResult.
func FormatGroupResult(r GroupResult) string {
	var sb strings.Builder

	prefixes := make([]string, 0, len(r.Groups))
	for p := range r.Groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	for _, p := range prefixes {
		members := r.Groups[p]
		keys := make([]string, 0, len(members))
		for k := range members {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("[%s] (%d keys)\n", p, len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", k, members[k]))
		}
	}

	if len(r.Ungrouped) > 0 {
		keys := make([]string, 0, len(r.Ungrouped))
		for k := range r.Ungrouped {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("[ungrouped] (%d keys)\n", len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", k, r.Ungrouped[k]))
		}
	}

	if sb.Len() == 0 {
		return "(no keys)\n"
	}
	return sb.String()
}
