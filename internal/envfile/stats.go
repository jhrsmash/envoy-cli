package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Stats holds aggregate statistics about an env map.
type Stats struct {
	Total       int
	Empty       int
	Sensitive   int
	PrefixCounts map[string]int
	AvgValueLen float64
}

// ComputeStats analyses an env map and returns a Stats summary.
func ComputeStats(env map[string]string) Stats {
	if len(env) == 0 {
		return Stats{PrefixCounts: map[string]int{}}
	}

	var totalLen int
	prefixes := map[string]int{}
	var empty, sensitive int

	for k, v := range env {
		if v == "" {
			empty++
		}
		if isSensitive(k) {
			sensitive++
		}
		totalLen += len(v)

		if idx := strings.Index(k, "_"); idx > 0 {
			prefixes[k[:idx]]++
		}
	}

	return Stats{
		Total:        len(env),
		Empty:        empty,
		Sensitive:    sensitive,
		PrefixCounts: prefixes,
		AvgValueLen:  float64(totalLen) / float64(len(env)),
	}
}

// FormatStats returns a human-readable summary of Stats.
func FormatStats(s Stats) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total keys   : %d\n", s.Total))
	sb.WriteString(fmt.Sprintf("Empty values : %d\n", s.Empty))
	sb.WriteString(fmt.Sprintf("Sensitive    : %d\n", s.Sensitive))
	sb.WriteString(fmt.Sprintf("Avg val len  : %.1f chars\n", s.AvgValueLen))

	if len(s.PrefixCounts) > 0 {
		sb.WriteString("Prefix groups:\n")
		keys := make([]string, 0, len(s.PrefixCounts))
		for k := range s.PrefixCounts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %-20s %d\n", k+"_*", s.PrefixCounts[k]))
		}
	}

	return sb.String()
}
