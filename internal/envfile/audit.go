package envfile

import (
	"fmt"
	"strings"
	"time"
)

// AuditAction represents the type of change recorded in an audit entry.
type AuditAction string

const (
	AuditAdded   AuditAction = "added"
	AuditRemoved AuditAction = "removed"
	AuditChanged AuditAction = "changed"
)

// AuditEntry records a single change to an environment variable.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    AuditAction
	OldValue  string
	NewValue  string
	Source    string // e.g. filename or environment label
}

// AuditLog is an ordered list of audit entries.
type AuditLog []AuditEntry

// Audit compares two env maps and returns an AuditLog describing every change.
// source is a human-readable label for the origin of the change (e.g. a filename).
func Audit(before, after map[string]string, source string) AuditLog {
	var log AuditLog
	now := time.Now().UTC()

	for key, newVal := range after {
		oldVal, existed := before[key]
		if !existed {
			log = append(log, AuditEntry{
				Timestamp: now,
				Key:       key,
				Action:    AuditAdded,
				NewValue:  newVal,
				Source:    source,
			})
		} else if oldVal != newVal {
			log = append(log, AuditEntry{
				Timestamp: now,
				Key:       key,
				Action:    AuditChanged,
				OldValue:  oldVal,
				NewValue:  newVal,
				Source:    source,
			})
		}
	}

	for key, oldVal := range before {
		if _, exists := after[key]; !exists {
			log = append(log, AuditEntry{
				Timestamp: now,
				Key:       key,
				Action:    AuditRemoved,
				OldValue:  oldVal,
				Source:    source,
			})
		}
	}

	return log
}

// FormatAuditLog returns a human-readable representation of an AuditLog.
func FormatAuditLog(log AuditLog) string {
	if len(log) == 0 {
		return "No changes recorded.\n"
	}

	var sb strings.Builder
	for _, entry := range log {
		ts := entry.Timestamp.Format(time.RFC3339)
		switch entry.Action {
		case AuditAdded:
			fmt.Fprintf(&sb, "[%s] %s ADDED %s=%q (source: %s)\n",
				ts, entry.Action, entry.Key, entry.NewValue, entry.Source)
		case AuditRemoved:
			fmt.Fprintf(&sb, "[%s] %s REMOVED %s (was %q) (source: %s)\n",
				ts, entry.Action, entry.Key, entry.OldValue, entry.Source)
		case AuditChanged:
			fmt.Fprintf(&sb, "[%s] %s CHANGED %s: %q -> %q (source: %s)\n",
				ts, entry.Action, entry.Key, entry.OldValue, entry.NewValue, entry.Source)
		}
	}
	return sb.String()
}
