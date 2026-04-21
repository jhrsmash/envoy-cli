// Package envfile provides utilities for managing .env files.
//
// # Template Rendering
//
// The template module allows rendering text templates that contain ${KEY}
// placeholders against a parsed env map. This is useful for generating
// configuration files (e.g. docker-compose.yml, nginx.conf) from a base
// template and a specific environment's values.
//
// Basic usage:
//
//	env, _ := envfile.Parse("production.env")
//	result := envfile.RenderTemplate(templateStr, env)
//	if len(result.Missing) > 0 {
//	    fmt.Println("unresolved:", result.Missing)
//	}
//	fmt.Print(result.Rendered)
//
// To render directly from a file on disk:
//
//	result, err := envfile.RenderTemplateFile("nginx.conf.tmpl", env)
package envfile
