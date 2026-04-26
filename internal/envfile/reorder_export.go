package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ReorderExportFormat represents the supported output formats for an exported
// reordered env file.
type ReorderExportFormat string

const (
	ReorderFormatDotenv ReorderExportFormat = "dotenv"
	ReorderFormatJSON   ReorderExportFormat = "json"
)

// ExportReordered writes env keys in the order defined by result.Ordered to
// a string in the requested format.
func ExportReordered(env map[string]string, result ReorderResult, format ReorderExportFormat) (string, error) {
	switch format {
	case ReorderFormatDotenv:
		return exportReorderedDotenv(env, result.Ordered), nil
	case ReorderFormatJSON:
		return exportReorderedJSON(env, result.Ordered)
	default:
		return "", fmt.Errorf("unsupported reorder export format: %q", format)
	}
}

// ExportReorderedToFile writes the reordered env to a file.
func ExportReorderedToFile(path string, env map[string]string, result ReorderResult, format ReorderExportFormat) error {
	contents, err := ExportReordered(env, result, format)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func exportReorderedDotenv(env map[string]string, order []string) string {
	var sb strings.Builder
	for _, k := range order {
		v, ok := env[k]
		if !ok {
			continue
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func exportReorderedJSON(env map[string]string, order []string) (string, error) {
	// Use a slice of objects to preserve order in JSON output.
	type kv struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	pairs := make([]kv, 0, len(order))
	for _, k := range order {
		if v, ok := env[k]; ok {
			pairs = append(pairs, kv{Key: k, Value: v})
		}
	}
	b, err := json.MarshalIndent(pairs, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

// inferReorderFormat guesses a ReorderExportFormat from a file extension.
func inferReorderFormat(path string) ReorderExportFormat {
	if strings.HasSuffix(path, ".json") {
		return ReorderFormatJSON
	}
	return ReorderFormatDotenv
}
