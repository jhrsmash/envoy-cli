package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// PivotResult holds the output of a Pivot operation.
type PivotResult struct {
	// Rows maps row-key -> (column-key -> value).
	Rows map[string]map[string]string
	// Columns is the ordered list of unique column labels.
	Columns []string
	// RowKeys is the ordered list of row keys after stripping the prefix.
	RowKeys []string
}

// PivotOptions controls how Pivot behaves.
type PivotOptions struct {
	// Delimiter separates the prefix from the rest of the key (default "_").
	Delimiter string
	// Prefixes lists the env-key prefixes that become columns.
	// Each matching key is split into (prefix, remainder); the remainder
	// becomes the row key and the prefix becomes the column header.
	Prefixes []string
}

// Pivot reorganises env keys that share a common suffix (row key) but differ
// by a leading prefix (column) into a two-dimensional table.
//
// Example input:
//
//	DEV_DB_HOST=localhost  PROD_DB_HOST=db.example.com
//	DEV_DB_PORT=5432       PROD_DB_PORT=5432
//
// With prefixes ["DEV", "PROD"] and delimiter "_" the result is:
//
//	row      | DEV       | PROD
//	DB_HOST  | localhost | db.example.com
//	DB_PORT  | 5432      | 5432
func Pivot(env map[string]string, opts PivotOptions) (PivotResult, error) {
	if len(opts.Prefixes) == 0 {
		return PivotResult{}, fmt.Errorf("pivot: at least one prefix must be specified")
	}
	delim := opts.Delimiter
	if delim == "" {
		delim = "_"
	}

	prefixSet := make(map[string]struct{}, len(opts.Prefixes))
	for _, p := range opts.Prefixes {
		prefixSet[p] = struct{}{}
	}

	rows := make(map[string]map[string]string)
	rowOrder := make(map[string]struct{})

	for k, v := range env {
		for _, prefix := range opts.Prefixes {
			candidate := prefix + delim
			if strings.HasPrefix(k, candidate) {
				rowKey := k[len(candidate):]
				if _, ok := rows[rowKey]; !ok {
					rows[rowKey] = make(map[string]string)
				}
				rows[rowKey][prefix] = v
				rowOrder[rowKey] = struct{}{}
				break
			}
		}
	}

	rowKeys := make([]string, 0, len(rowOrder))
	for rk := range rowOrder {
		rowKeys = append(rowKeys, rk)
	}
	sort.Strings(rowKeys)

	cols := make([]string, len(opts.Prefixes))
	copy(cols, opts.Prefixes)

	return PivotResult{
		Rows:    rows,
		Columns: cols,
		RowKeys: rowKeys,
	}, nil
}
