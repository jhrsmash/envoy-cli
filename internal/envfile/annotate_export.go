package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type annotateJSONEntry struct {
	Key     string `json:"key"`
	Comment string `json:"comment"`
}

// ExportAnnotations writes the annotation map to a file in dotenv or JSON format.
// Format is inferred from the file extension if not specified.
func ExportAnnotations(r AnnotateResult, path string, format string) error {
	if format == "" {
		format = inferAnnotateFormat(path)
	}
	switch strings.ToLower(format) {
	case "json":
		return exportAnnotationsJSON(r, path)
	case "dotenv", "env":
		return exportAnnotationsDotenv(r, path)
	default:
		return fmt.Errorf("unsupported annotation export format: %s", format)
	}
}

func exportAnnotationsJSON(r AnnotateResult, path string) error {
	keys := make([]string, 0, len(r.Annotations))
	for k := range r.Annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]annotateJSONEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, annotateJSONEntry{Key: k, Comment: r.Annotations[k]})
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func exportAnnotationsDotenv(r AnnotateResult, path string) error {
	keys := make([]string, 0, len(r.Annotations))
	for k := range r.Annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s\n%s=\n\n", r.Annotations[k], k))
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func inferAnnotateFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json"
	default:
		return "dotenv"
	}
}
