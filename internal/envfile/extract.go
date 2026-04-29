package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ExtractOptions controls how keys are extracted from an env map.
type ExtractOptions struct {
	// Keys is an explicit list of keys to extract. If empty, all keys are returned.
	Keys []string
	// Prefix filters keys to only those matching the given prefix.
	Prefix string
	// StripPrefix removes the prefix from the resulting keys.
	StripPrefix bool
	// FailOnMissing returns an error if any explicitly requested key is absent.
	FailOnMissing bool
}

// ExtractResult holds the outcome of an Extract operation.
type ExtractResult struct {
	Extracted map[string]string
	Missing   []string
}

// Extract pulls a subset of keys from env according to the given options.
// It returns an ExtractResult containing the matched key/value pairs and any
// missing keys that were explicitly requested.
func Extract(env map[string]string, opts ExtractOptions) (ExtractResult, error) {
	extracted := make(map[string]string)
	var missing []string

	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			v, ok := env[k]
			if !ok {
				missing = append(missing, k)
				continue
			}
			key := k
			if opts.StripPrefix && opts.Prefix != "" {
				key = strings.TrimPrefix(k, opts.Prefix)
			}
			extracted[key] = v
		}
	} else {
		for k, v := range env {
			if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
				continue
			}
			key := k
			if opts.StripPrefix && opts.Prefix != "" {
				key = strings.TrimPrefix(k, opts.Prefix)
			}
			extracted[key] = v
		}
	}

	if opts.FailOnMissing && len(missing) > 0 {
		sort.Strings(missing)
		return ExtractResult{}, fmt.Errorf("extract: missing keys: %s", strings.Join(missing, ", "))
	}

	return ExtractResult{
		Extracted: extracted,
		Missing:   missing,
	}, nil
}

// FormatExtractResult returns a human-readable summary of the extraction.
func FormatExtractResult(r ExtractResult) string {
	var sb strings.Builder

	keys := make([]string, 0, len(r.Extracted))
	for k := range r.Extracted {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(&sb, "  extracted: %s=%s\n", k, r.Extracted[k])
	}

	if len(r.Missing) > 0 {
		missing := append([]string(nil), r.Missing...)
		sort.Strings(missing)
		for _, k := range missing {
			fmt.Fprintf(&sb, "  missing:   %s\n", k)
		}
	}

	return sb.String()
}
