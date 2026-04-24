package envfile

import "fmt"

// CopyOptions controls the behaviour of a Copy operation.
type CopyOptions struct {
	// Keys is the explicit list of keys to copy. If empty, all keys are copied.
	Keys []string
	// Overwrite controls whether existing keys in dst are replaced.
	Overwrite bool
}

// CopyResult holds the outcome of a Copy operation.
type CopyResult struct {
	Copied   []string
	Skipped  []string
	Missing  []string
}

// Copy copies keys from src into dst according to opts.
// It never mutates src.
func Copy(src, dst map[string]string, opts CopyOptions) (map[string]string, CopyResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	keys := opts.Keys
	if len(keys) == 0 {
		keys = sortedKeys(src)
	}

	var result CopyResult
	for _, k := range keys {
		v, exists := src[k]
		if !exists {
			result.Missing = append(result.Missing, k)
			continue
		}
		if _, inDst := out[k]; inDst && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		out[k] = v
		result.Copied = append(result.Copied, k)
	}
	return out, result
}

// FormatCopyResult returns a human-readable summary of a CopyResult.
func FormatCopyResult(r CopyResult) string {
	if len(r.Copied) == 0 && len(r.Skipped) == 0 && len(r.Missing) == 0 {
		return "copy: nothing to do\n"
	}
	out := ""
	for _, k := range r.Copied {
		out += fmt.Sprintf("  copied:  %s\n", k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("  skipped: %s (already exists)\n", k)
	}
	for _, k := range r.Missing {
		out += fmt.Sprintf("  missing: %s (not in source)\n", k)
	}
	return out
}
