package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportScope writes the scoped env map to a file.
// The format is inferred from the file extension (.json, .sh, or .env).
func ExportScope(r ScopeResult, destPath string) error {
	format := inferScopeFormat(destPath)

	var content string
	switch format {
	case "json":
		b, err := json.MarshalIndent(r.Scoped, "", "  ")
		if err != nil {
			return fmt.Errorf("scope export: json marshal: %w", err)
		}
		content = string(b) + "\n"
	case "shell":
		var sb strings.Builder
		for _, k := range sortedKeys(r.Scoped) {
			fmt.Fprintf(&sb, "export %s=%q\n", k, r.Scoped[k])
		}
		content = sb.String()
	default: // dotenv
		var sb strings.Builder
		for _, k := range sortedKeys(r.Scoped) {
			fmt.Fprintf(&sb, "%s=%s\n", k, r.Scoped[k])
		}
		content = sb.String()
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("scope export: mkdir: %w", err)
	}
	return os.WriteFile(destPath, []byte(content), 0o644)
}

func inferScopeFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json"
	case ".sh":
		return "shell"
	default:
		return "dotenv"
	}
}
