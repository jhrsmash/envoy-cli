package envfile

import (
	"fmt"
	"strings"
)

// FormatSearchResult renders a SearchResult as a human-readable string.
func FormatSearchResult(r SearchResult) string {
	if len(r.Matches) == 0 {
		return "No matches found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d match(es) found:\n", len(r.Matches)))

	for _, m := range r.Matches {
		var tags []string
		if m.MatchedKey {
			tags = append(tags, "key")
		}
		if m.MatchedVal {
			tags = append(tags, "value")
		}
		sb.WriteString(fmt.Sprintf("  %-30s = %s  [matched: %s]\n",
			m.Key, m.Value, strings.Join(tags, ", ")))
	}
	return sb.String()
}

// FormatSearchSummary returns a one-line summary of the result.
func FormatSearchSummary(r SearchResult) string {
	if len(r.Matches) == 0 {
		return "search: no matches"
	}
	return fmt.Sprintf("search: %d match(es) for %q", len(r.Matches), r.Options.Query)
}
