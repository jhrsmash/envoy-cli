package envfile

import (
	"fmt"
	"strings"
)

// ReorderOptions controls how keys are reordered within an env map.
type ReorderOptions struct {
	// Keys defines the explicit ordering of keys. Keys not listed here
	// will be appended after the explicitly ordered keys (in their original
	// relative order, or alphabetically if AlphaTail is true).
	Keys []string

	// AlphaTail sorts any remaining keys alphabetically after the explicit list.
	AlphaTail bool
}

// ReorderResult holds the outcome of a Reorder operation.
type ReorderResult struct {
	Ordered   []string // final key order
	Missing   []string // keys listed in opts.Keys that were not found in env
	Unlisted  []string // keys in env that were not listed in opts.Keys
}

// Reorder returns a new map equal to env, and a ReorderResult describing
// the final key ordering. The map itself is unordered; callers should use
// result.Ordered to iterate keys in the desired sequence.
func Reorder(env map[string]string, opts ReorderOptions) (map[string]string, ReorderResult) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	seen := make(map[string]bool)
	var ordered []string
	var missing []string

	for _, k := range opts.Keys {
		if _, ok := env[k]; !ok {
			missing = append(missing, k)
			continue
		}
		if !seen[k] {
			ordered = append(ordered, k)
			seen[k] = true
		}
	}

	// Collect unlisted keys.
	var unlisted []string
	for k := range env {
		if !seen[k] {
			unlisted = append(unlisted, k)
		}
	}

	if opts.AlphaTail {
		unlisted = sortStrings(unlisted)
	}

	ordered = append(ordered, unlisted...)

	return out, ReorderResult{
		Ordered:  ordered,
		Missing:  missing,
		Unlisted: unlisted,
	}
}

// sortStrings returns a sorted copy of ss.
func sortStrings(ss []string) []string {
	copy_ := make([]string, len(ss))
	copy(copy_, ss)
	// insertion sort — small slices expected
	for i := 1; i < len(copy_); i++ {
		for j := i; j > 0 && copy_[j] < copy_[j-1]; j-- {
			copy_[j], copy_[j-1] = copy_[j-1], copy_[j]
		}
	}
	return copy_
}

// FormatReorderResult returns a human-readable summary of a ReorderResult.
func FormatReorderResult(r ReorderResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Reordered: %d keys\n", len(r.Ordered))
	if len(r.Missing) > 0 {
		fmt.Fprintf(&sb, "Missing (not in env): %s\n", strings.Join(r.Missing, ", "))
	}
	if len(r.Unlisted) > 0 {
		fmt.Fprintf(&sb, "Unlisted (appended):  %s\n", strings.Join(r.Unlisted, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
