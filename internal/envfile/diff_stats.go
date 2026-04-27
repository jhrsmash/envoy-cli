package envfile

import "fmt"

// DiffStats holds aggregate counts derived from a slice of DiffEntry values.
type DiffStats struct {
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
}

// ComputeDiffStats returns a DiffStats summary for the provided diff entries.
func ComputeDiffStats(entries []DiffEntry) DiffStats {
	var s DiffStats
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed + s.Unchanged
	return s
}

// FormatDiffStats returns a human-readable summary of a DiffStats value.
func FormatDiffStats(s DiffStats) string {
	if s.Total == 0 {
		return "diff stats: no entries"
	}
	changed := s.Added + s.Removed + s.Changed
	out := fmt.Sprintf("diff stats: %d total", s.Total)
	if changed == 0 {
		out += ", no changes"
		return out
	}
	if s.Added > 0 {
		out += fmt.Sprintf(", +%d added", s.Added)
	}
	if s.Removed > 0 {
		out += fmt.Sprintf(", -%d removed", s.Removed)
	}
	if s.Changed > 0 {
		out += fmt.Sprintf(", ~%d changed", s.Changed)
	}
	if s.Unchanged > 0 {
		out += fmt.Sprintf(", %d unchanged", s.Unchanged)
	}
	return out
}
