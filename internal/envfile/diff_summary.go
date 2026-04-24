package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// DiffSummary holds aggregated statistics about a diff result.
type DiffSummary struct {
	Added    int
	Removed  int
	Changed  int
	Total    int
	AddedKeys   []string
	RemovedKeys []string
	ChangedKeys []string
}

// SummarizeDiff computes a DiffSummary from a slice of DiffEntry values
// returned by Diff.
func SummarizeDiff(entries []DiffEntry) DiffSummary {
	var s DiffSummary
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
			s.AddedKeys = append(s.AddedKeys, e.Key)
		case StatusRemoved:
			s.Removed++
			s.RemovedKeys = append(s.RemovedKeys, e.Key)
		case StatusChanged:
			s.Changed++
			s.ChangedKeys = append(s.ChangedKeys, e.Key)
		}
	}
	sort.Strings(s.AddedKeys)
	sort.Strings(s.RemovedKeys)
	sort.Strings(s.ChangedKeys)
	s.Total = s.Added + s.Removed + s.Changed
	return s
}

// FormatDiffSummary returns a human-readable summary string for a DiffSummary.
func FormatDiffSummary(s DiffSummary) string {
	if s.Total == 0 {
		return "No differences found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Diff summary: %d change(s) total\n", s.Total))
	if s.Added > 0 {
		sb.WriteString(fmt.Sprintf("  + Added   (%d): %s\n", s.Added, strings.Join(s.AddedKeys, ", ")))
	}
	if s.Removed > 0 {
		sb.WriteString(fmt.Sprintf("  - Removed (%d): %s\n", s.Removed, strings.Join(s.RemovedKeys, ", ")))
	}
	if s.Changed > 0 {
		sb.WriteString(fmt.Sprintf("  ~ Changed (%d): %s\n", s.Changed, strings.Join(s.ChangedKeys, ", ")))
	}
	return strings.TrimRight(sb.String(), "\n")
}
