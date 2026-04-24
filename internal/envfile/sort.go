package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// SortOptions controls how environment variable keys are sorted.
type SortOptions struct {
	// Reverse reverses the sort order.
	Reverse bool
	// GroupByPrefix groups keys sharing a common prefix together.
	GroupByPrefix bool
}

// SortResult holds the output of a Sort operation.
type SortResult struct {
	Sorted map[string]string
	Keys   []string // keys in sorted order
}

// Sort returns a SortResult with keys ordered according to SortOptions.
// The original env map is not mutated.
func Sort(env map[string]string, opts SortOptions) SortResult {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if opts.GroupByPrefix {
		sort.Slice(keys, func(i, j int) bool {
			pi := primaryPrefix(keys[i])
			pj := primaryPrefix(keys[j])
			if pi != pj {
				if opts.Reverse {
					return pi > pj
				}
				return pi < pj
			}
			if opts.Reverse {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			if opts.Reverse {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	}

	sorted := make(map[string]string, len(env))
	for k, v := range env {
		sorted[k] = v
	}

	return SortResult{Sorted: sorted, Keys: keys}
}

// FormatSortResult renders the sorted env as a dotenv-style string.
func FormatSortResult(r SortResult) string {
	var sb strings.Builder
	for _, k := range r.Keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Sorted[k])
	}
	return sb.String()
}

// primaryPrefix returns the portion of a key before the first underscore,
// used for prefix-based grouping.
func primaryPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
