package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportOverrideResult writes the result of an OS override operation to a file.
// Supported formats: "json", "dotenv". Format is inferred from the file extension
// when format is empty.
func ExportOverrideResult(r EnvOverrideResult, path string, format string) error {
	if format == "" {
		format = inferOverrideFormat(path)
	}
	switch strings.ToLower(format) {
	case "json":
		return exportOverrideJSON(r, path)
	case "dotenv", "env":
		return exportOverrideDotenv(r, path)
	default:
		return fmt.Errorf("unsupported override export format: %q", format)
	}
}

func exportOverrideJSON(r EnvOverrideResult, path string) error {
	type jsonPayload struct {
		Env     map[string]string `json:"env"`
		Applied []EnvOverride     `json:"applied"`
		Skipped []string          `json:"skipped"`
	}
	payload := jsonPayload{
		Env:     r.Env,
		Applied: r.Applied,
		Skipped: r.Skipped,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal override result: %w", err)
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func exportOverrideDotenv(r EnvOverrideResult, path string) error {
	keys := make([]string, 0, len(r.Env))
	for k := range r.Env {
		keys = append(keys, k)
	}
	sortStrings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, r.Env[k]))
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func inferOverrideFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json"
	default:
		return "dotenv"
	}
}
