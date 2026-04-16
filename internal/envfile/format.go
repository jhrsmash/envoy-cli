package envfile

import (
	"fmt"
	"io"
	"sort"
)

// FormatDiff writes a human-readable diff to the provided writer.
// Lines are prefixed with '+' for added, '-' for removed, and '~' for changed.
func FormatDiff(w io.Writer, d DiffResult) {
	if d.IsEmpty() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(d.Added) > 0 {
		fmt.Fprintln(w, "Added:")
		for _, k := range sortedKeys(d.Added) {
			fmt.Fprintf(w, "  + %s=%s\n", k, d.Added[k])
		}
	}

	if len(d.Removed) > 0 {
		fmt.Fprintln(w, "Removed:")
		for _, k := range sortedKeys(d.Removed) {
			fmt.Fprintf(w, "  - %s=%s\n", k, d.Removed[k])
		}
	}

	if len(d.Changed) > 0 {
		fmt.Fprintln(w, "Changed:")
		keys := make([]string, 0, len(d.Changed))
		for k := range d.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pair := d.Changed[k]
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", k, pair[0], pair[1])
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
