package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// DedupResult holds the outcome of a deduplication operation.
type DedupResult struct {
	// Output is the deduplicated env map.
	Output map[string]string
	// Duplicates maps each key that appeared more than once to the list of
	// values that were seen (in encounter order). The value kept in Output
	// is always the last one (override semantics, like most shell loaders).
	Duplicates map[string][]string
}

// Dedup scans a raw slice of key=value lines (as would be produced by reading
// an .env file line-by-line before deduplication) and returns a DedupResult.
// Keys that appear only once are placed in Output without a Duplicates entry.
// Comment lines (prefix '#') and blank lines are silently skipped.
func Dedup(lines []string) DedupResult {
	seen := make(map[string][]string) // key -> all values in order
	order := []string{}               // insertion order for deterministic output

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		// Strip surrounding quotes (single or double).
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if _, exists := seen[key]; !exists {
			order = append(order, key)
		}
		seen[key] = append(seen[key], val)
	}

	out := make(map[string]string, len(order))
	duplicates := make(map[string][]string)

	for _, key := range order {
		vals := seen[key]
		out[key] = vals[len(vals)-1] // last value wins
		if len(vals) > 1 {
			duplicates[key] = vals
		}
	}

	return DedupResult{Output: out, Duplicates: duplicates}
}

// FormatDedupResult returns a human-readable summary of a DedupResult.
func FormatDedupResult(r DedupResult) string {
	if len(r.Duplicates) == 0 {
		return "No duplicate keys found.\n"
	}

	keys := make([]string, 0, len(r.Duplicates))
	for k := range r.Duplicates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d duplicate key(s) resolved:\n", len(keys))
	for _, k := range keys {
		vals := r.Duplicates[k]
		fmt.Fprintf(&sb, "  %s (%d occurrences, kept: %q)\n", k, len(vals), vals[len(vals)-1])
	}
	return sb.String()
}
