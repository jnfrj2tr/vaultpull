package env

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(string) (string, error)

// Transformer applies named transformations to env maps.
type Transformer struct {
	funcs map[string]TransformFunc
}

// NewTransformer returns a Transformer with built-in transforms registered.
func NewTransformer() *Transformer {
	t := &Transformer{funcs: make(map[string]TransformFunc)}
	t.Register("upper", func(v string) (string, error) { return strings.ToUpper(v), nil })
	t.Register("lower", func(v string) (string, error) { return strings.ToLower(v), nil })
	t.Register("trim", func(v string) (string, error) { return strings.TrimSpace(v), nil })
	t.Register("base64", func(v string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	})
	return t
}

// Register adds a named transform function.
func (t *Transformer) Register(name string, fn TransformFunc) {
	t.funcs[name] = fn
}

// Apply applies the named transform to each value in the map.
// If only is non-empty, only those keys are transformed.
func (t *Transformer) Apply(secrets map[string]string, transform string, only []string) (map[string]string, error) {
	fn, ok := t.funcs[transform]
	if !ok {
		return nil, fmt.Errorf("unknown transform %q", transform)
	}
	filter := make(map[string]bool, len(only))
	for _, k := range only {
		filter[k] = true
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(filter) == 0 || filter[k] {
			var err error
			v, err = fn(v)
			if err != nil {
				return nil, fmt.Errorf("transform %q failed on key %q: %w", transform, k, err)
			}
		}
		out[k] = v
	}
	return out, nil
}
