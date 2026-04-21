package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateResult holds the output of rendering a template against an env map.
type TemplateResult struct {
	Rendered string
	Missing  []string
}

var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// RenderTemplate replaces ${KEY} placeholders in the template string with
// values from env. Missing keys are collected in TemplateResult.Missing and
// left as-is in the rendered output.
func RenderTemplate(tmpl string, env map[string]string) TemplateResult {
	missing := []string{}
	seen := map[string]bool{}

	rendered := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := templateVarRe.FindStringSubmatch(match)[1]
		if val, ok := env[key]; ok {
			return val
		}
		if !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
		return match
	})

	return TemplateResult{Rendered: rendered, Missing: missing}
}

// RenderTemplateFile reads a template file from disk and renders it against env.
func RenderTemplateFile(path string, env map[string]string) (TemplateResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TemplateResult{}, fmt.Errorf("template: read %q: %w", path, err)
	}
	return RenderTemplate(string(data), env), nil
}

// FormatTemplateResult formats a TemplateResult for display.
func FormatTemplateResult(r TemplateResult) string {
	var sb strings.Builder
	sb.WriteString(r.Rendered)
	if len(r.Missing) > 0 {
		sb.WriteString("\n# WARNING: unresolved template variables: ")
		sb.WriteString(strings.Join(r.Missing, ", "))
		sb.WriteString("\n")
	}
	return sb.String()
}
