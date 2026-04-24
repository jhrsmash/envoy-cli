package envfile

import "sort"

// SyncOptions controls the behaviour of Sync.
type SyncOptions struct {
	// AddMissing adds keys present in source but absent in target.
	AddMissing bool
	// RemoveExtra removes keys present in target but absent in source.
	RemoveExtra bool
	// Overwrite updates keys whose values differ between source and target.
	Overwrite bool
}

// SyncResult describes what changed after a Sync operation.
type SyncResult struct {
	Added   []string
	Removed []string
	Updated []string
	Output  map[string]string
}

// Sync reconciles target with source according to opts.
// It never mutates the input maps; it returns a new map in SyncResult.Output.
func Sync(source, target map[string]string, opts SyncOptions) SyncResult {
	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var added, removed, updated []string

	if opts.AddMissing || opts.Overwrite {
		for k, sv := range source {
			tv, exists := out[k]
			if !exists && opts.AddMissing {
				out[k] = sv
				added = append(added, k)
			} else if exists && sv != tv && opts.Overwrite {
				out[k] = sv
				updated = append(updated, k)
			}
		}
	}

	if opts.RemoveExtra {
		for k := range target {
			if _, inSource := source[k]; !inSource {
				delete(out, k)
				removed = append(removed, k)
			}
		}
	}

	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(updated)

	return SyncResult{
		Added:   added,
		Removed: removed,
		Updated: updated,
		Output:  out,
	}
}
