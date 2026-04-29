package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const rollbackIndexFile = ".envoy_rollback_index.json"

// rollbackIndex is the on-disk format for the rollback archive list.
type rollbackIndex struct {
	Archives []ArchiveEntry `json:"archives"`
}

// SaveRollbackIndex persists the list of archive entries to dir.
func SaveRollbackIndex(dir string, entries []ArchiveEntry) error {
	idx := rollbackIndex{Archives: entries}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("rollback: marshal index: %w", err)
	}
	path := filepath.Join(dir, rollbackIndexFile)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("rollback: write index: %w", err)
	}
	return nil
}

// LoadRollbackIndex reads the archive list from dir.
// Returns an empty slice when no index exists yet.
func LoadRollbackIndex(dir string) ([]ArchiveEntry, error) {
	path := filepath.Join(dir, rollbackIndexFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("rollback: read index: %w", err)
	}
	var idx rollbackIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("rollback: unmarshal index: %w", err)
	}
	return idx.Archives, nil
}
