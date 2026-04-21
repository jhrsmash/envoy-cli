package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportToFile writes the exported env content to a file at the given path.
// The format is inferred from the file extension if opts.Format is empty.
func ExportToFile(env map[string]string, path string, opts ExportOptions) error {
	if opts.Format == "" {
		opts.Format = inferFormat(path)
	}

	content, err := Export(env, opts)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// inferFormat guesses the ExportFormat from a file's extension.
func inferFormat(path string) ExportFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return FormatJSON
	case ".sh":
		return FormatShell
	default:
		// .env, .txt, or no extension → dotenv
		return FormatDotenv
	}
}
