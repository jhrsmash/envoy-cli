package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Classification holds the result of classifying env keys by inferred purpose.
type Classification struct {
	// Groups maps a category label to the keys that belong to it.
	Groups map[string][]string
	// Unclassified holds keys that did not match any known category.
	Unclassified []string
}

// categoryRules maps a category name to a list of key substrings that imply membership.
var categoryRules = map[string][]string{
	"auth":     {"TOKEN", "SECRET", "PASSWORD", "PASSWD", "API_KEY", "APIKEY", "AUTH", "CREDENTIAL", "CERT", "PRIVATE"},
	"database": {"DB_", "DATABASE", "POSTGRES", "MYSQL", "MONGO", "REDIS", "DSN", "JDBC"},
	"network":  {"HOST", "PORT", "URL", "URI", "ADDR", "ENDPOINT", "PROXY", "DOMAIN"},
	"logging":  {"LOG", "LOGLEVEL", "LOG_LEVEL", "DEBUG", "VERBOSE", "TRACE"},
	"feature":  {"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_", "FF_"},
}

// Classify inspects each key in env and assigns it to one or more
// well-known categories based on naming conventions.
// Keys that match no rule are placed in Unclassified.
func Classify(env map[string]string) Classification {
	groups := make(map[string][]string)
	var unclassified []string

	for key := range env {
		upper := strings.ToUpper(key)
		matched := false

		for category, patterns := range categoryRules {
			for _, p := range patterns {
				if strings.Contains(upper, p) {
					groups[category] = append(groups[category], key)
					matched = true
					break
				}
			}
		}

		if !matched {
			unclassified = append(unclassified, key)
		}
	}

	for cat := range groups {
		sort.Strings(groups[cat])
	}
	sort.Strings(unclassified)

	return Classification{Groups: groups, Unclassified: unclassified}
}

// FormatClassifyResult returns a human-readable summary of a Classification.
func FormatClassifyResult(c Classification) string {
	var sb strings.Builder

	categories := make([]string, 0, len(c.Groups))
	for cat := range c.Groups {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		keys := c.Groups[cat]
		sb.WriteString(fmt.Sprintf("[%s] (%d key(s))\n", cat, len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	if len(c.Unclassified) > 0 {
		sb.WriteString(fmt.Sprintf("[unclassified] (%d key(s))\n", len(c.Unclassified)))
		for _, k := range c.Unclassified {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	if sb.Len() == 0 {
		return "no keys to classify\n"
	}
	return sb.String()
}
