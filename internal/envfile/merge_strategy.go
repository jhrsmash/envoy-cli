package envfile

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging two env maps.
type MergeStrategy int

const (
	// StrategyPreferBase keeps the base value on conflict.
	StrategyPreferBase MergeStrategy = iota
	// StrategyPreferOverride replaces base values with override values on conflict.
	StrategyPreferOverride
	// StrategyErrorOnConflict returns an error if any key exists in both maps with different values.
	StrategyErrorOnConflict
	// StrategyKeepBoth appends a suffix to the conflicting override key and retains both values.
	StrategyKeepBoth
)

// MergeStrategyResult holds the output of a strategy-aware merge.
type MergeStrategyResult struct {
	// Merged is the resulting env map after applying the strategy.
	Merged map[string]string
	// Conflicts lists keys that had differing values in base and override.
	Conflicts []string
	// Added lists keys that were new from the override map.
	Added []string
}

// MergeWithStrategy merges override into base using the provided strategy.
// base and override are not mutated.
func MergeWithStrategy(base, override map[string]string, strategy MergeStrategy) (MergeStrategyResult, error) {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	var conflicts []string
	var added []string

	for k, ov := range override {
		bv, exists := result[k]
		if !exists {
			result[k] = ov
			added = append(added, k)
			continue
		}
		if bv == ov {
			// No conflict — values are identical.
			continue
		}
		// Conflict detected.
		conflicts = append(conflicts, k)
		switch strategy {
		case StrategyPreferBase:
			// Keep base value; do nothing.
		case StrategyPreferOverride:
			result[k] = ov
		case StrategyErrorOnConflict:
			return MergeStrategyResult{}, fmt.Errorf(
				"merge conflict on key %q: base=%q override=%q", k, bv, ov,
			)
		case StrategyKeepBoth:
			// Retain base value under original key; store override under a suffixed key.
			suffixed := k + "__override"
			result[suffixed] = ov
			added = append(added, suffixed)
		}
	}

	return MergeStrategyResult{
		Merged:    result,
		Conflicts: sortedKeys(sliceToMap(conflicts)),
		Added:     sortedKeys(sliceToMap(added)),
	}, nil
}

// FormatMergeStrategyResult returns a human-readable summary of a strategy merge.
func FormatMergeStrategyResult(r MergeStrategyResult) string {
	if len(r.Conflicts) == 0 && len(r.Added) == 0 {
		return "merge complete: no conflicts, no new keys\n"
	}
	out := fmt.Sprintf("merge complete: %d conflict(s), %d new key(s)\n",
		len(r.Conflicts), len(r.Added))
	for _, k := range r.Conflicts {
		out += fmt.Sprintf("  ~ %s (conflict)\n", k)
	}
	for _, k := range r.Added {
		out += fmt.Sprintf("  + %s (added)\n", k)
	}
	return out
}

// sliceToMap converts a string slice to a set-like map for use with sortedKeys.
func sliceToMap(keys []string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = ""
	}
	return m
}
