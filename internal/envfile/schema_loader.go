package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadSchema reads a schema definition file and returns a Schema.
//
// Each non-blank, non-comment line has the format:
//
//	KEY [required|optional] [type]
//
// Example:
//
//	DATABASE_URL required url
//	PORT         optional int
//	DEBUG        optional bool
func LoadSchema(path string) (Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return Schema{}, fmt.Errorf("schema: cannot open %q: %w", path, err)
	}
	defer f.Close()

	var schema Schema
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 1 {
			return Schema{}, fmt.Errorf("schema: line %d is malformed", lineNum)
		}

		field := SchemaField{Key: parts[0]}

		if len(parts) >= 2 {
			switch strings.ToLower(parts[1]) {
			case "required":
				field.Required = true
			case "optional":
				field.Required = false
			default:
				return Schema{}, fmt.Errorf("schema: line %d: unknown modifier %q (use 'required' or 'optional')", lineNum, parts[1])
			}
		}

		if len(parts) >= 3 {
			field.Pattern = strings.ToLower(parts[2])
		}

		schema.Fields = append(schema.Fields, field)
	}

	if err := scanner.Err(); err != nil {
		return Schema{}, fmt.Errorf("schema: read error: %w", err)
	}

	return schema, nil
}
