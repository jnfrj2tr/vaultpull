// Package template provides .env file rendering from Go text/template strings.
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Data holds the values available inside a template.
type Data struct {
	Secrets map[string]string
	Env     map[string]string
}

// Renderer renders a template file using secret and environment values.
type Renderer struct {
	data Data
}

// New creates a Renderer with the provided secrets and ambient env vars.
func New(secrets map[string]string) *Renderer {
	env := make(map[string]string)
	for _, pair := range os.Environ() {
		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				env[pair[:i]] = pair[i+1:]
				break
			}
		}
	}
	return &Renderer{data: Data{Secrets: secrets, Env: env}}
}

// RenderFile reads a template from srcPath, executes it, and writes the result
// to dstPath, creating or truncating the destination file.
func (r *Renderer) RenderFile(srcPath, dstPath string) error {
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("template: read source: %w", err)
	}
	out, err := r.RenderString(string(src))
	if err != nil {
		return err
	}
	if err := os.WriteFile(dstPath, []byte(out), 0600); err != nil {
		return fmt.Errorf("template: write destination: %w", err)
	}
	return nil
}

// RenderString executes a template string and returns the rendered output.
func (r *Renderer) RenderString(tmpl string) (string, error) {
	t, err := template.New("").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template: parse: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, r.data); err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}
	return buf.String(), nil
}
