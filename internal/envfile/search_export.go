package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExportSearchResult writes a SearchResult to a file in the given format.
// Supported formats: "json", "text".
func ExportSearchResult(r SearchResult, path string, format string) error {
	if format == "" {
		format = inferSearchFormat(path)
	}

	var data []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		data, err = marshalSearchJSON(r)
		if err != nil {
			return err
		}
	case "text", "":
		data = []byte(FormatSearchResult(r))
	default:
		return fmt.Errorf("unsupported search export format: %q", format)
	}

	return os.WriteFile(path, data, 0o644)
}

func marshalSearchJSON(r SearchResult) ([]byte, error) {
	type jsonMatch struct {
		Key        string `json:"key"`
		Value      string `json:"value"`
		MatchedKey bool   `json:"matched_key"`
		MatchedVal bool   `json:"matched_value"`
	}
	type jsonResult struct {
		Query   string      `json:"query"`
		Count   int         `json:"count"`
		Matches []jsonMatch `json:"matches"`
	}

	out := jsonResult{
		Query:   r.Options.Query,
		Count:   len(r.Matches),
		Matches: make([]jsonMatch, len(r.Matches)),
	}
	for i, m := range r.Matches {
		out.Matches[i] = jsonMatch{Key: m.Key, Value: m.Value, MatchedKey: m.MatchedKey, MatchedVal: m.MatchedVal}
	}
	return json.MarshalIndent(out, "", "  ")
}

func inferSearchFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return "json"
	default:
		return "text"
	}
}
