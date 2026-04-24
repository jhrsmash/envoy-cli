package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiffExportFormat enumerates supported output formats for exporting a diff.
type DiffExportFormat string

const (
	DiffExportJSON   DiffExportFormat = "json"
	DiffExportPatch  DiffExportFormat = "patch"
)

// ExportDiff serialises diff entries to the requested format string.
func ExportDiff(entries []DiffEntry, format DiffExportFormat) (string, error) {
	switch format {
	case DiffExportJSON:
		return exportDiffJSON(entries)
	case DiffExportPatch:
		return exportDiffPatch(entries), nil
	default:
		return "", fmt.Errorf("unsupported diff export format: %q", format)
	}
}

// ExportDiffToFile writes the diff to a file, inferring format from extension.
func ExportDiffToFile(entries []DiffEntry, path string) error {
	fmt := inferDiffFormat(path)
	out, err := ExportDiff(entries, fmt)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(out), 0o644)
}

func exportDiffJSON(entries []DiffEntry) (string, error) {
	type jsonEntry struct {
		Key      string `json:"key"`
		Status   string `json:"status"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
	}
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = jsonEntry{
			Key:      e.Key,
			Status:   string(e.Status),
			OldValue: e.OldValue,
			NewValue: e.NewValue,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func exportDiffPatch(entries []DiffEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, e.NewValue))
		case StatusRemoved:
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, e.OldValue))
		case StatusChanged:
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, e.OldValue))
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, e.NewValue))
		}
	}
	return sb.String()
}

func inferDiffFormat(path string) DiffExportFormat {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return DiffExportJSON
	default:
		return DiffExportPatch
	}
}
