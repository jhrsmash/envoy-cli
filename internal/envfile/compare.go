package envfile

import "sort"

// CompareResult holds the result of comparing two env maps by key presence.
type CompareResult struct {
	OnlyInA    []string
	OnlyInB    []string
	InBoth     []string
	TotalA     int
	TotalB     int
}

// Compare returns which keys exist only in a, only in b, or in both.
// It does not compare values — use Diff for value-level diffing.
func Compare(a, b map[string]string) CompareResult {
	aSet := make(map[string]bool, len(a))
	bSet := make(map[string]bool, len(b))

	for k := range a {
		aSet[k] = true
	}
	for k := range b {
		bSet[k] = true
	}

	var onlyA, onlyB, both []string

	for k := range aSet {
		if bSet[k] {
			both = append(both, k)
		} else {
			onlyA = append(onlyA, k)
		}
	}
	for k := range bSet {
		if !aSet[k] {
			onlyB = append(onlyB, k)
		}
	}

	sort.Strings(onlyA)
	sort.Strings(onlyB)
	sort.Strings(both)

	return CompareResult{
		OnlyInA: onlyA,
		OnlyInB: onlyB,
		InBoth:  both,
		TotalA:  len(a),
		TotalB:  len(b),
	}
}

// FormatCompareResult returns a human-readable summary of a CompareResult.
func FormatCompareResult(r CompareResult, labelA, labelB string) string {
	var out string
	out += fmt.Sprintf("Keys only in %s (%d):\n", labelA, len(r.OnlyInA))
	for _, k := range r.OnlyInA {
		out += fmt.Sprintf("  - %s\n", k)
	}
	out += fmt.Sprintf("Keys only in %s (%d):\n", labelB, len(r.OnlyInB))
	for _, k := range r.OnlyInB {
		out += fmt.Sprintf("  - %s\n", k)
	}
	out += fmt.Sprintf("Keys in both: %d\n", len(r.InBoth))
	return out
}
