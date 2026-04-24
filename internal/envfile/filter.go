package envfile

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// FilterOptions controls how Filter selects keys from an env map.
type FilterOptions struct {
	// Prefix retains only keys that start with the given prefix.
	Prefix string
	// Suffix retains only keys that end with the given suffix.
	Suffix string
	// Pattern retains only keys matching the given regular expression.
	Pattern string
	// Keys retains only the explicitly listed keys (takes precedence over
	// Prefix/Suffix/Pattern when non-empty).
	Keys []string
}

// FilterResult holds the output of a Filter operation.
type FilterResult struct {
	Matched map[string]string
	Dropped []string
}

// Filter returns a new map containing only the entries from env that satisfy
// the criteria in opts. All criteria are ANDed together (each must match).
func Filter(env map[string]string, opts FilterOptions) (FilterResult, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return FilterResult{}, fmt.Errorf("invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	explicit := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = true
	}

	matched := make(map[string]string)
	var dropped []string

	for k, v := range env {
		if len(explicit) > 0 {
			if explicit[k] {
				matched[k] = v
			} else {
				dropped = append(dropped, k)
			}
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			dropped = append(dropped, k)
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, opts.Suffix) {
			dropped = append(dropped, k)
			continue
		}
		if re != nil && !re.MatchString(k) {
			dropped = append(dropped, k)
			continue
		}
		matched[k] = v
	}

	sort.Strings(dropped)
	return FilterResult{Matched: matched, Dropped: dropped}, nil
}

// FormatFilterResult returns a human-readable summary of a FilterResult.
func FormatFilterResult(r FilterResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Matched %d key(s):\n", len(r.Matched))
	keys := make([]string, 0, len(r.Matched))
	for k := range r.Matched {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&sb, "  + %s\n", k)
	}
	if len(r.Dropped) > 0 {
		fmt.Fprintf(&sb, "Dropped %d key(s):\n", len(r.Dropped))
		for _, k := range r.Dropped {
			fmt.Fprintf(&sb, "  - %s\n", k)
		}
	}
	return sb.String()
}
