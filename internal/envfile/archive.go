package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ArchiveEntry represents a single archived version of an env map.
type ArchiveEntry struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// ArchiveResult holds the outcome of an Archive operation.
type ArchiveResult struct {
	Label   string
	Keys    int
	Stored  string
}

// Archive saves a snapshot of env under a named label to the given directory.
// Each label creates or overwrites a file named <label>.archive.json.
func Archive(env map[string]string, label, dir string) (ArchiveResult, error) {
	if label == "" {
		return ArchiveResult{}, fmt.Errorf("archive label must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ArchiveResult{}, fmt.Errorf("create archive dir: %w", err)
	}

	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}

	entry := ArchiveEntry{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Env:       copy,
	}

	path := filepath.Join(dir, label+".archive.json")
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("marshal archive: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return ArchiveResult{}, fmt.Errorf("write archive: %w", err)
	}

	return ArchiveResult{Label: label, Keys: len(env), Stored: path}, nil
}

// LoadArchive reads a previously archived entry by label from dir.
func LoadArchive(label, dir string) (ArchiveEntry, error) {
	path := filepath.Join(dir, label+".archive.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return ArchiveEntry{}, fmt.Errorf("read archive %q: %w", label, err)
	}
	var entry ArchiveEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return ArchiveEntry{}, fmt.Errorf("parse archive %q: %w", label, err)
	}
	return entry, nil
}

// ListArchives returns all archive labels found in dir.
func ListArchives(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list archives: %w", err)
	}
	var labels []string
	for _, e := range entries {
		name := e.Name()
		if filepath.Ext(name) == ".json" && len(name) > len(".archive.json") {
			labels = append(labels, name[:len(name)-len(".archive.json")])
		}
	}
	return labels, nil
}
