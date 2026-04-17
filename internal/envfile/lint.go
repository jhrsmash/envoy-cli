package envfile

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning for an env file.
type LintIssue struct {
	Key     string
	Message string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("%s: %s", l.Key, l.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

// OK returns true if no issues were found.
func (r *LintResult) OK() bool {
	return len(r.Issues) == 0
}

// Lint checks an env map for common style and quality issues.
// It warns about:
//   - keys that are not UPPER_SNAKE_CASE
//   - empty values
//   - keys that shadow common reserved names
func Lint(env map[string]string) *LintResult {
	reserved := map[string]bool{
		"PATH": true, "HOME": true, "USER": true, "SHELL": true, "PWD": true,
	}

	result := &LintResult{}

	for key, val := range env {
		if key != strings.ToUpper(key) {
			result.Issues = append(result.Issues, LintIssue{
				Key:     key,
				Message: "key should be UPPER_SNAKE_CASE",
			})
		}

		if strings.TrimSpace(val) == "" {
			result.Issues = append(result.Issues, LintIssue{
				Key:     key,
				Message: "value is empty",
			})
		}

		if reserved[strings.ToUpper(key)] {
			result.Issues = append(result.Issues, LintIssue{
				Key:     key,
				Message: "key shadows a reserved system environment variable",
			})
		}
	}

	return result
}

// FormatLintResult returns a human-readable string of lint issues.
func FormatLintResult(r *LintResult) string {
	if r.OK() {
		return "No lint issues found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d lint issue(s) found:\n", len(r.Issues)))
	for _, issue := range r.Issues {
		sb.WriteString(fmt.Sprintf("  - %s\n", issue.String()))
	}
	return sb.String()
}
