package envfile

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// PreferBase keeps the base value on conflict.
	PreferBase MergeStrategy = iota
	// PreferOverride uses the override value on conflict.
	PreferOverride
)

// MergeResult holds the merged env map and metadata about the operation.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string // keys that had differing values
	Added     []string // keys only in override
}

// Merge combines base and override env maps according to the given strategy.
// Keys present only in base are always kept.
// Keys present only in override are always added.
// Conflicting keys are resolved by strategy.
func Merge(base, override map[string]string, strategy MergeStrategy) MergeResult {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []string
	var added []string

	for k, v := range override {
		baseVal, exists := merged[k]
		if !exists {
			merged[k] = v
			added = append(added, k)
			continue
		}
		if baseVal != v {
			conflicts = append(conflicts, k)
			if strategy == PreferOverride {
				merged[k] = v
			}
		}
	}

	return MergeResult{
		Merged:    merged,
		Conflicts: conflicts,
		Added:     added,
	}
}
