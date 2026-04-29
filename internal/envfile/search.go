package envfile

import (
	"regexp"
	"sort"
	"strings"
)

// SearchOptions controls how Search matches keys and values.
type SearchOptions struct {
	// Query is the search term (plain text or regex if UseRegex is true).
	Query string
	// UseRegex treats Query as a regular expression.
	UseRegex bool
	// CaseSensitive disables case-folding when false (default: case-insensitive).
	CaseSensitive bool
	// SearchKeys includes key names in the match scope.
	SearchKeys bool
	// SearchValues includes values in the match scope.
	SearchValues bool
}

// SearchMatch represents a single matched entry.
type SearchMatch struct {
	Key        string
	Value      string
	MatchedKey bool
	MatchedVal bool
}

// SearchResult is the output of a Search operation.
type SearchResult struct {
	Matches []SearchMatch
	Options SearchOptions
}

// Search scans env for keys/values matching opts and returns a SearchResult.
func Search(env map[string]string, opts SearchOptions) (SearchResult, error) {
	if !opts.SearchKeys && !opts.SearchValues {
		opts.SearchKeys = true
		opts.SearchValues = true
	}

	var re *regexp.Regexp
	if opts.UseRegex {
		pattern := opts.Query
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return SearchResult{}, err
		}
	}

	match := func(s string) bool {
		if re != nil {
			return re.MatchString(s)
		}
		q, t := opts.Query, s
		if !opts.CaseSensitive {
			q = strings.ToLower(q)
			t = strings.ToLower(t)
		}
		return strings.Contains(t, q)
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var matches []SearchMatch
	for _, k := range keys {
		v := env[k]
		mk := opts.SearchKeys && match(k)
		mv := opts.SearchValues && match(v)
		if mk || mv {
			matches = append(matches, SearchMatch{Key: k, Value: v, MatchedKey: mk, MatchedVal: mv})
		}
	}

	return SearchResult{Matches: matches, Options: opts}, nil
}
