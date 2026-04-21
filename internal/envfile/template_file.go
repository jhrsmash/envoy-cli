package envfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RenderTemplateToFile renders a template file and writes the result to an
// output path. If the output file already exists it will be overwritten.
// Returns the TemplateResult so callers can inspect missing keys.
func RenderTemplateToFile(tmplPath, outPath string, env map[string]string) (TemplateResult, error) {
	result, err := RenderTemplateFile(tmplPath, env)
	if err != nil {
		return TemplateResult{}, err
	}

	if dir := filepath.Dir(outPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return TemplateResult{}, fmt.Errorf("template: mkdir %q: %w", dir, err)
		}
	}

	if err := os.WriteFile(outPath, []byte(result.Rendered), 0644); err != nil {
		return TemplateResult{}, fmt.Errorf("template: write %q: %w", outPath, err)
	}

	return result, nil
}

// inferTemplateName derives a default output filename by stripping a .tmpl
// suffix from the template filename, if present.
func inferTemplateName(tmplPath string) string {
	base := filepath.Base(tmplPath)
	if strings.HasSuffix(base, ".tmpl") {
		base = strings.TrimSuffix(base, ".tmpl")
	}
	return filepath.Join(filepath.Dir(tmplPath), base)
}

// RenderTemplateAuto renders a .tmpl file and writes the output to the same
// directory with the .tmpl suffix removed (e.g. nginx.conf.tmpl -> nginx.conf).
func RenderTemplateAuto(tmplPath string, env map[string]string) (TemplateResult, string, error) {
	outPath := inferTemplateName(tmplPath)
	if outPath == tmplPath {
		outPath = tmplPath + ".rendered"
	}
	result, err := RenderTemplateToFile(tmplPath, outPath, env)
	if err != nil {
		return TemplateResult{}, "", err
	}
	return result, outPath, nil
}
