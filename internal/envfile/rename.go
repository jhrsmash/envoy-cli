package envfile

import "fmt"

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	OldKey  string
	NewKey  string
	Renamed bool
	Reason  string
}

// RenameOptions controls the behaviour of Rename.
type RenameOptions struct {
	// Overwrite allows the new key to replace an existing key.
	Overwrite bool
}

// Rename renames oldKey to newKey in env.
// It returns a copy of the map with the rename applied and a RenameResult
// describing what happened.
func Rename(env map[string]string, oldKey, newKey string, opts RenameOptions) (map[string]string, RenameResult) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := RenameResult{OldKey: oldKey, NewKey: newKey}

	if oldKey == "" || newKey == "" {
		result.Reason = "key names must not be empty"
		return out, result
	}

	if oldKey == newKey {
		result.Reason = "old and new key are identical"
		return out, result
	}

	val, exists := out[oldKey]
	if !exists {
		result.Reason = fmt.Sprintf("key %q not found", oldKey)
		return out, result
	}

	if _, conflict := out[newKey]; conflict && !opts.Overwrite {
		result.Reason = fmt.Sprintf("key %q already exists; use Overwrite to replace it", newKey)
		return out, result
	}

	delete(out, oldKey)
	out[newKey] = val
	result.Renamed = true
	return out, result
}

// FormatRenameResult returns a human-readable summary of a RenameResult.
func FormatRenameResult(r RenameResult) string {
	if r.Renamed {
		return fmt.Sprintf("renamed %q → %q", r.OldKey, r.NewKey)
	}
	return fmt.Sprintf("rename skipped: %s", r.Reason)
}
