package envfile

import (
	"fmt"
	"strings"
)

// NormalizeOptions controls how normalization is applied.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to UPPERCASE.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// RemoveEmptyValues drops keys whose values are empty after trimming.
	RemoveEmptyValues bool
	// QuoteValues wraps values that contain spaces in double quotes.
	QuoteValues bool
}

// NormalizeResult holds the output of a Normalize call.
type NormalizeResult struct {
	Env     map[string]string
	Changes []NormalizeChange
}

// NormalizeChange describes a single normalization mutation.
type NormalizeChange struct {
	Key    string
	OldKey string // non-empty when the key itself was renamed
	OldVal string
	NewVal string
	Reason string
}

// Normalize applies the given options to env and returns a new map along with
// a list of changes that were made. The original map is never mutated.
func Normalize(env map[string]string, opts NormalizeOptions) NormalizeResult {
	out := make(map[string]string, len(env))
	var changes []NormalizeChange

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.UppercaseKeys {
			newKey = strings.ToUpper(k)
		}

		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		if opts.RemoveEmptyValues && newVal == "" {
			changes = append(changes, NormalizeChange{
				Key:    newKey,
				OldKey: k,
				OldVal: v,
				NewVal: "",
				Reason: "removed empty value",
			})
			continue
		}

		if opts.QuoteValues && strings.ContainsAny(newVal, " \t") {
			if !strings.HasPrefix(newVal, `"`) {
				newVal = fmt.Sprintf(`"%s"`, newVal)
			}
		}

		if newKey != k || newVal != v {
			oldKey := ""
			if newKey != k {
				oldKey = k
			}
			changes = append(changes, NormalizeChange{
				Key:    newKey,
				OldKey: oldKey,
				OldVal: v,
				NewVal: newVal,
				Reason: buildReason(k, newKey, v, newVal),
			})
		}

		out[newKey] = newVal
	}

	return NormalizeResult{Env: out, Changes: changes}
}

func buildReason(oldKey, newKey, oldVal, newVal string) string {
	var parts []string
	if oldKey != newKey {
		parts = append(parts, "key uppercased")
	}
	if oldVal != newVal {
		if strings.TrimSpace(oldVal) == strings.TrimSpace(newVal) {
			parts = append(parts, "value trimmed")
		} else {
			parts = append(parts, "value quoted")
		}
	}
	return strings.Join(parts, ", ")
}
