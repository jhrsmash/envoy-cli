package envfile

import (
	"fmt"
	"sort"
	"time"
)

// RollbackOptions controls how a rollback is performed.
type RollbackOptions struct {
	// Label selects a specific archive entry to roll back to.
	// If empty, the most recent archive entry is used.
	Label string

	// DryRun reports what would change without applying it.
	DryRun bool
}

// RollbackResult holds the outcome of a rollback operation.
type RollbackResult struct {
	Label     string
	Timestamp time.Time
	Applied   map[string]string // final env after rollback
	Diff      []DiffEntry       // changes relative to current env
	DryRun    bool
}

// Rollback restores a previously archived snapshot of the env.
// current is the live environment; archives is the list of saved archives.
// The target archive is selected by opts.Label (or most-recent if blank).
func Rollback(current map[string]string, archives []ArchiveEntry, opts RollbackOptions) (RollbackResult, error) {
	if len(archives) == 0 {
		return RollbackResult{}, fmt.Errorf("rollback: no archives available")
	}

	// Sort archives newest-first.
	sorted := make([]ArchiveEntry, len(archives))
	copy(sorted, archives)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})

	var target *ArchiveEntry
	if opts.Label == "" {
		target = &sorted[0]
	} else {
		for i := range sorted {
			if sorted[i].Label == opts.Label {
				target = &sorted[i]
				break
			}
		}
		if target == nil {
			return RollbackResult{}, fmt.Errorf("rollback: archive label %q not found", opts.Label)
		}
	}

	diffEntries := Diff(current, target.Env)

	applied := target.Env
	if opts.DryRun {
		applied = current
	}

	return RollbackResult{
		Label:     target.Label,
		Timestamp: target.Timestamp,
		Applied:   applied,
		Diff:      diffEntries,
		DryRun:    opts.DryRun,
	}, nil
}
