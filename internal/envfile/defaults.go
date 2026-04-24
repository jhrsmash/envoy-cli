package envfile

import "fmt"

// DefaultsOptions configures how defaults are applied.
type DefaultsOptions struct {
	// Overwrite replaces existing keys with defaults when true.
	Overwrite bool
}

// DefaultsResult holds the outcome of applying defaults.
type DefaultsResult struct {
	Applied  map[string]string // keys that were set from defaults
	Skipped  map[string]string // keys that already existed and were not overwritten
}

// Defaults applies a map of default key-value pairs to env.
// Keys already present in env are skipped unless Overwrite is true.
// Neither env nor defaults is mutated; a new map is returned.
func Defaults(env, defaults map[string]string, opts DefaultsOptions) (map[string]string, DefaultsResult) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := DefaultsResult{
		Applied: make(map[string]string),
		Skipped: make(map[string]string),
	}

	for k, v := range defaults {
		if existing, exists := out[k]; exists && !opts.Overwrite {
			result.Skipped[k] = existing
			continue
		}
		out[k] = v
		result.Applied[k] = v
	}

	return out, result
}

// FormatDefaultsResult returns a human-readable summary of a DefaultsResult.
func FormatDefaultsResult(r DefaultsResult) string {
	if len(r.Applied) == 0 && len(r.Skipped) == 0 {
		return "defaults: nothing to apply\n"
	}

	var out string
	for _, k := range sortedKeys(r.Applied) {
		out += fmt.Sprintf("  applied  %s=%s\n", k, r.Applied[k])
	}
	for _, k := range sortedKeys(r.Skipped) {
		out += fmt.Sprintf("  skipped  %s (already set)\n", k)
	}
	return out
}
