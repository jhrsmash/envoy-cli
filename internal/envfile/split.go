package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// SplitOptions controls how an env map is split into named buckets.
type SplitOptions struct {
	// Prefixes maps a bucket name to a key prefix (e.g. {"db": "DB_", "app": "APP_"}).
	Prefixes map[string]string
	// StripPrefix removes the matched prefix from keys in each bucket.
	StripPrefix bool
	// CatchAll is the bucket name for keys that match no prefix.
	// If empty, unmatched keys are dropped.
	CatchAll string
}

// SplitResult holds the output of a Split operation.
type SplitResult struct {
	Buckets  map[string]map[string]string
	Unmatched []string
}

// Split partitions env into named buckets based on key prefixes.
// Keys are matched in deterministic order (sorted bucket names).
func Split(env map[string]string, opts SplitOptions) (SplitResult, error) {
	if len(opts.Prefixes) == 0 {
		return SplitResult{}, fmt.Errorf("split: at least one prefix must be specified")
	}

	// Build sorted bucket names for deterministic matching.
	names := make([]string, 0, len(opts.Prefixes))
	for name := range opts.Prefixes {
		names = append(names, name)
	}
	sort.Strings(names)

	// Initialise buckets.
	buckets := make(map[string]map[string]string, len(names))
	for _, name := range names {
		buckets[name] = make(map[string]string)
	}
	if opts.CatchAll != "" {
		if _, exists := buckets[opts.CatchAll]; !exists {
			buckets[opts.CatchAll] = make(map[string]string)
		}
	}

	var unmatched []string

	for k, v := range env {
		matched := false
		for _, name := range names {
			prefix := opts.Prefixes[name]
			if strings.HasPrefix(k, prefix) {
				key := k
				if opts.StripPrefix {
					key = strings.TrimPrefix(k, prefix)
				}
				buckets[name][key] = v
				matched = true
				break
			}
		}
		if !matched {
			if opts.CatchAll != "" {
				buckets[opts.CatchAll][k] = v
			} else {
				unmatched = append(unmatched, k)
			}
		}
	}

	sort.Strings(unmatched)
	return SplitResult{Buckets: buckets, Unmatched: unmatched}, nil
}

// FormatSplitResult returns a human-readable summary of a SplitResult.
func FormatSplitResult(r SplitResult) string {
	var sb strings.Builder

	names := make([]string, 0, len(r.Buckets))
	for name := range r.Buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		bucket := r.Buckets[name]
		sb.WriteString(fmt.Sprintf("[%s] %d key(s)\n", name, len(bucket)))
		keys := make([]string, 0, len(bucket))
		for k := range bucket {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", k, bucket[k]))
		}
	}

	if len(r.Unmatched) > 0 {
		sb.WriteString(fmt.Sprintf("unmatched: %s\n", strings.Join(r.Unmatched, ", ")))
	}

	return sb.String()
}
