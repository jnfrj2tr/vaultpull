package env

import (
	"encoding/base64"
	"testing"
)

func TestApply_Upper(t *testing.T) {
	tr := NewTransformer()
	in := map[string]string{"KEY": "hello", "OTHER": "world"}
	out, err := tr.Apply(in, "upper", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "HELLO" || out["OTHER"] != "WORLD" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApply_Lower(t *testing.T) {
	tr := NewTransformer()
	in := map[string]string{"A": "FOO", "B": "BAR"}
	out, err := tr.Apply(in, "lower", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "foo" || out["B"] != "bar" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApply_Trim(t *testing.T) {
	tr := NewTransformer()
	in := map[string]string{"K": "  spaced  "}
	out, err := tr.Apply(in, "trim", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "spaced" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApply_Base64(t *testing.T) {
	tr := NewTransformer()
	in := map[string]string{"SECRET": "mysecret"}
	out, err := tr.Apply(in, "base64", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := base64.StdEncoding.EncodeToString([]byte("mysecret"))
	if out["SECRET"] != want {
		t.Errorf("got %q want %q", out["SECRET"], want)
	}
}

func TestApply_OnlySubset(t *testing.T) {
	tr := NewTransformer()
	in := map[string]string{"A": "hello", "B": "world"}
	out, err := tr.Apply(in, "upper", []string{"A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "HELLO" {
		t.Errorf("A should be uppercased, got %q", out["A"])
	}
	if out["B"] != "world" {
		t.Errorf("B should be unchanged, got %q", out["B"])
	}
}

func TestApply_UnknownTransform(t *testing.T) {
	tr := NewTransformer()
	_, err := tr.Apply(map[string]string{"K": "v"}, "rot13", nil)
	if err == nil {
		t.Fatal("expected error for unknown transform")
	}
}

func TestRegister_CustomTransform(t *testing.T) {
	tr := NewTransformer()
	tr.Register("exclaim", func(v string) (string, error) { return v + "!", nil })
	out, err := tr.Apply(map[string]string{"MSG": "hi"}, "exclaim", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MSG"] != "hi!" {
		t.Errorf("got %q", out["MSG"])
	}
}
