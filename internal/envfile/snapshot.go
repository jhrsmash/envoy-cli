package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// NewSnapshot creates a new Snapshot from an env map.
func NewSnapshot(label string, env map[string]string) Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Env:       copy,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadSnapshot reads a Snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return s, nil
}
