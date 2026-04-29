package envfile

import (
	"fmt"
	"strings"
)

const pivotEmpty = "(unset)"

// FormatPivotResult renders a PivotResult as a plain-text table.
//
// Example output:
//
//	KEY          DEV          PROD
//	-----------  -----------  -----------
//	DB_HOST      localhost    db.example.com
//	DB_PORT      5432         5432
func FormatPivotResult(r PivotResult) string {
	if len(r.RowKeys) == 0 {
		return "(no matching keys)"
	}

	// Build column widths: first column is the row-key column.
	headers := append([]string{"KEY"}, r.Columns...)
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, rk := range r.RowKeys {
		if len(rk) > widths[0] {
			widths[0] = len(rk)
		}
		for ci, col := range r.Columns {
			v := pivotEmpty
			if row, ok := r.Rows[rk]; ok {
				if val, ok2 := row[col]; ok2 {
					v = val
				}
			}
			if len(v) > widths[ci+1] {
				widths[ci+1] = len(v)
			}
		}
	}

	fmtRow := func(cells []string) string {
		parts := make([]string, len(cells))
		for i, c := range cells {
			parts[i] = fmt.Sprintf("%-*s", widths[i], c)
		}
		return strings.Join(parts, "  ")
	}

	var sb strings.Builder
	sb.WriteString(fmtRow(headers))
	sb.WriteByte('\n')

	sep := make([]string, len(headers))
	for i, w := range widths {
		sep[i] = strings.Repeat("-", w)
	}
	sb.WriteString(fmtRow(sep))
	sb.WriteByte('\n')

	for _, rk := range r.RowKeys {
		cells := make([]string, len(headers))
		cells[0] = rk
		for ci, col := range r.Columns {
			v := pivotEmpty
			if row, ok := r.Rows[rk]; ok {
				if val, ok2 := row[col]; ok2 {
					v = val
				}
			}
			cells[ci+1] = v
		}
		sb.WriteString(fmtRow(cells))
		sb.WriteByte('\n')
	}

	return strings.TrimRight(sb.String(), "\n")
}
