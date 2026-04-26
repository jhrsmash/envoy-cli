package envfile

import "sort"

// ScopeOptions controls how Scope filters and transforms an env map.
type ScopeOptions struct {
	// Prefixes is the list of key prefixes to include (e.g. ["APP_", "DB_"]).
	Prefixes []string
	// StripPrefix removes the matched prefix from the resulting keys.
	StripPrefix bool
	// Uppercase converts all resulting keys to uppercase.
	Uppercase bool
}

// ScopeResult holds the output of a Scope operation.
type ScopeResult struct {
	Scoped   map[string]string
	Included []string // keys included (after any renaming)
	Excluded []string // keys excluded
}

// Scope returns a new map containing only keys that match one of the given
// prefixes. If StripPrefix is true the matched prefix is removed from each key.
func Scope(env map[string]string, opts ScopeOptions) ScopeResult {
	scoped := make(map[string]string)
	included := []string{}
	excluded := []string{}

	for k, v := range env {
		matched, prefix := matchesPrefix(k, opts.Prefixes)
		if !matched {
			excluded = append(excluded, k)
			continue
		}
		newKey := k
		if opts.StripPrefix && prefix != "" {
			newKey = k[len(prefix):]
		}
		if opts.Uppercase {
			newKey = toUpper(newKey)
		}
		scoped[newKey] = v
		included = append(included, newKey)
	}

	sort.Strings(included)
	sort.Strings(excluded)
	return ScopeResult{Scoped: scoped, Included: included, Excluded: excluded}
}

func matchesPrefix(key string, prefixes []string) (bool, string) {
	if len(prefixes) == 0 {
		return true, ""
	}
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true, p
		}
	}
	return false, ""
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
