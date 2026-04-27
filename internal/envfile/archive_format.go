package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// FormatArchiveResult returns a human-readable summary of an Archive operation.
func FormatArchiveResult(r ArchiveResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Archived %d key(s) under label %q\n", r.Keys, r.Label))
	sb.WriteString(fmt.Sprintf("  Stored at: %s\n", r.Stored))
	return sb.String()
}

// FormatArchiveEntry returns a human-readable view of a loaded ArchiveEntry.
func FormatArchiveEntry(e ArchiveEntry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Archive: %s\n", e.Label))
	sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", e.Timestamp.Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString(fmt.Sprintf("  Keys (%d):\n", len(e.Env)))

	keys := make([]string, 0, len(e.Env))
	for k := range e.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("    %s=%s\n", k, e.Env[k]))
	}
	return sb.String()
}

// FormatArchiveList returns a formatted list of available archive labels.
func FormatArchiveList(labels []string) string {
	if len(labels) == 0 {
		return "No archives found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Available archives (%d):\n", len(labels)))
	for _, l := range labels {
		sb.WriteString(fmt.Sprintf("  - %s\n", l))
	}
	return sb.String()
}
