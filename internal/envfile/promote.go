package envfile

import "fmt"

// PromoteResult holds the outcome of promoting keys from one environment to another.
type PromoteResult struct {
	// Promoted contains keys that were copied from source to destination.
	Promoted []string
	// Skipped contains keys that already existed in destination and were not overwritten.
	Skipped []string
	// Missing contains keys requested for promotion that were absent in source.
	Missing []string
}

// Promote copies selected keys from src into dst.
// If keys is empty, all keys from src are promoted.
// Existing keys in dst are skipped unless overwrite is true.
func Promote(src, dst map[string]string, keys []string, overwrite bool) (map[string]string, PromoteResult) {
	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	var pr PromoteResult

	candidates := keys
	if len(candidates) == 0 {
		candidates = sortedKeys(src)
	}

	for _, k := range candidates {
		v, ok := src[k]
		if !ok {
			pr.Missing = append(pr.Missing, k)
			continue
		}
		if _, exists := result[k]; exists && !overwrite {
			pr.Skipped = append(pr.Skipped, k)
			continue
		}
		result[k] = v
		pr.Promoted = append(pr.Promoted, k)
	}

	return result, pr
}

// FormatPromoteResult returns a human-readable summary of a PromoteResult.
func FormatPromoteResult(pr PromoteResult, srcLabel, dstLabel string) string {
	out := fmt.Sprintf("Promote: %s → %s\n", srcLabel, dstLabel)

	if len(pr.Promoted) == 0 && len(pr.Skipped) == 0 && len(pr.Missing) == 0 {
		out += "  (nothing to promote)\n"
		return out
	}

	for _, k := range pr.Promoted {
		out += fmt.Sprintf("  + promoted: %s\n", k)
	}
	for _, k := range pr.Skipped {
		out += fmt.Sprintf("  ~ skipped:  %s (already exists)\n", k)
	}
	for _, k := range pr.Missing {
		out += fmt.Sprintf("  ! missing:  %s (not found in source)\n", k)
	}

	return out
}
