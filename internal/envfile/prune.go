package envfile

import "sort"

// PruneOptions controls which keys are removed during pruning.
type PruneOptions struct {
	// RemoveEmpty removes keys whose values are empty strings.
	RemoveEmpty bool
	// RemoveKeys is an explicit list of keys to remove.
	RemoveKeys []string
	// RemovePrefix removes all keys that start with the given prefix.
	RemovePrefix string
}

// PruneResult holds the output of a Prune operation.
type PruneResult struct {
	Output  map[string]string
	Removed []string
}

// Prune removes keys from env according to the given options.
// The original map is never mutated.
func Prune(env map[string]string, opts PruneOptions) PruneResult {
	explicit := make(map[string]bool, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		explicit[k] = true
	}

	out := make(map[string]string, len(env))
	var removed []string

	for k, v := range env {
		switch {
		case explicit[k]:
			removed = append(removed, k)
		case opts.RemoveEmpty && v == "":
			removed = append(removed, k)
		case opts.RemovePrefix != "" && len(k) >= len(opts.RemovePrefix) && k[:len(opts.RemovePrefix)] == opts.RemovePrefix:
			removed = append(removed, k)
		default:
			out[k] = v
		}
	}

	sort.Strings(removed)
	return PruneResult{Output: out, Removed: removed}
}
