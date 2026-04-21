package envfile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ExportFormat represents a supported export target format.
type ExportFormat string

const (
	FormatJSON   ExportFormat = "json"
	FormatDotenv ExportFormat = "dotenv"
	FormatShell  ExportFormat = "shell"
)

// ExportOptions controls how the export is rendered.
type ExportOptions struct {
	Format  ExportFormat
	Sorted  bool
	Redacted bool
}

// Export converts an env map to the specified format string.
func Export(env map[string]string, opts ExportOptions) (string, error) {
	data := env
	if opts.Redacted {
		data = Redact(env, nil)
	}

	switch opts.Format {
	case FormatJSON:
		return exportJSON(data, opts.Sorted)
	case FormatDotenv:
		return exportDotenv(data, opts.Sorted), nil
	case FormatShell:
		return exportShell(data, opts.Sorted), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportJSON(env map[string]string, sorted bool) (string, error) {
	var out interface{}
	if sorted {
		ordered := make(map[string]string, len(env))
		for k, v := range env {
			ordered[k] = v
		}
		out = ordered
	} else {
		out = env
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}

func exportDotenv(env map[string]string, sorted bool) string {
	keys := keysOf(env, sorted)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportShell(env map[string]string, sorted bool) string {
	keys := keysOf(env, sorted)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func keysOf(env map[string]string, sorted bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if sorted {
		sort.Strings(keys)
	}
	return keys
}
