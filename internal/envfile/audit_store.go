package envfile

import (
	"encoding/json"
	"errors"
	"os"
)

// SaveAuditLog persists an AuditLog to a JSON file at the given path.
// If the file already contains a valid log, new entries are appended.
func SaveAuditLog(path string, incoming AuditLog) error {
	if len(incoming) == 0 {
		return nil
	}

	existing, err := LoadAuditLog(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	combined := append(existing, incoming...)

	data, err := json.MarshalIndent(combined, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// LoadAuditLog reads and deserialises an AuditLog from a JSON file.
// Returns os.ErrNotExist if the file does not exist.
func LoadAuditLog(path string) (AuditLog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, err
	}

	return log, nil
}

// ClearAuditLog removes the audit log file at the given path.
// Returns nil if the file does not exist.
func ClearAuditLog(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
