package envfile

// DiffResult holds the result of comparing two env maps.
type DiffResult struct {
	Added   map[string]string // keys present in target but not in base
	Removed map[string]string // keys present in base but not in target
	Changed map[string][2]string // keys present in both but with different values [base, target]
}

// Diff compares two parsed env maps (base vs target) and returns a DiffResult.
func Diff(base, target map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	// Find removed and changed keys
	for k, baseVal := range base {
		if targetVal, ok := target[k]; !ok {
			result.Removed[k] = baseVal
		} else if baseVal != targetVal {
			result.Changed[k] = [2]string{baseVal, targetVal}
		}
	}

	// Find added keys
	for k, targetVal := range target {
		if _, ok := base[k]; !ok {
			result.Added[k] = targetVal
		}
	}

	return result
}

// IsEmpty returns true if there are no differences.
func (d DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}
