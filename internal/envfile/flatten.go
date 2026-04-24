package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key segments are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join key segments (default: "_").
	Separator string
	// Uppercase converts all resulting keys to uppercase.
	Uppercase bool
	// Prefix is prepended to every key in the output.
	Prefix string
}

// FlattenResult holds the output of a Flatten operation.
type FlattenResult struct {
	// Output is the flattened env map.
	Output map[string]string
	// Renamed maps original keys to their new flattened keys (only entries that changed).
	Renamed map[string]string
	// Unchanged lists keys that required no transformation.
	Unchanged []string
}

// Flatten normalises keys in env by replacing any occurrence of the given
// separator alternatives ("-", ".", " ") with opts.Separator, optionally
// uppercasing keys and prepending a prefix. It never mutates the input map.
func Flatten(env map[string]string, opts FlattenOptions) FlattenResult {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	output := make(map[string]string, len(env))
	renamed := make(map[string]string)
	var unchanged []string

	for origKey, val := range env {
		newKey := origKey

		// Replace common delimiters with the target separator.
		for _, alt := range []string{"-", ".", " "} {
			if alt != opts.Separator {
				newKey = strings.ReplaceAll(newKey, alt, opts.Separator)
			}
		}

		if opts.Uppercase {
			newKey = strings.ToUpper(newKey)
		}

		if opts.Prefix != "" {
			newKey = opts.Prefix + newKey
		}

		output[newKey] = val

		if newKey != origKey {
			renamed[origKey] = newKey
		} else {
			unchanged = append(unchanged, origKey)
		}
	}

	sort.Strings(unchanged)
	return FlattenResult{Output: output, Renamed: renamed, Unchanged: unchanged}
}

// FormatFlattenResult returns a human-readable summary of a FlattenResult.
func FormatFlattenResult(r FlattenResult) string {
	var sb strings.Builder

	if len(r.Renamed) == 0 {
		sb.WriteString("flatten: no keys required renaming\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("flatten: %d key(s) renamed, %d unchanged\n", len(r.Renamed), len(r.Unchanged)))

	// Print renames in sorted order for deterministic output.
	origKeys := make([]string, 0, len(r.Renamed))
	for k := range r.Renamed {
		origKeys = append(origKeys, k)
	}
	sort.Strings(origKeys)

	for _, orig := range origKeys {
		sb.WriteString(fmt.Sprintf("  %s -> %s\n", orig, r.Renamed[orig]))
	}

	return sb.String()
}
