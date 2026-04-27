package envfile

import (
	"fmt"
	"strings"
)

// FormatWatchEvent returns a human-readable summary of a WatchEvent.
func FormatWatchEvent(e WatchEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s] %s changed\n",
		e.ChangedAt.Format("15:04:05"), e.File))
	sb.WriteString(fmt.Sprintf("  hash: %s → %s\n",
		shortHash(e.OldHash), shortHash(e.NewHash)))

	added, removed, changed, unchanged := 0, 0, 0, 0
	for _, d := range e.Diff {
		switch d.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusChanged:
			changed++
		case StatusUnchanged:
			unchanged++
		}
	}

	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("+%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("~%d changed", changed))
	}
	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", unchanged))
	}
	sb.WriteString(fmt.Sprintf("  diff: %s\n", strings.Join(parts, ", ")))

	for _, d := range e.Diff {
		switch d.Status {
		case StatusAdded:
			sb.WriteString(fmt.Sprintf("  + %s=%s\n", d.Key, d.NewValue))
		case StatusRemoved:
			sb.WriteString(fmt.Sprintf("  - %s=%s\n", d.Key, d.OldValue))
		case StatusChanged:
			sb.WriteString(fmt.Sprintf("  ~ %s: %q → %q\n", d.Key, d.OldValue, d.NewValue))
		}
	}

	return sb.String()
}

func shortHash(h string) string {
	if len(h) >= 12 {
		return h[:12]
	}
	return h
}
