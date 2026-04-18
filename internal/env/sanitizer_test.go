package env

import (
	"testing"
)

func TestSanitizeKey_AlreadyValid(t *testing.T) {
	res := SanitizeKey("MY_KEY")
	if res.Changed {
		t.Errorf("expected no change, got %q", res.Fixed)
	}
}

func TestSanitizeKey_LowercaseConverted(t *testing.T) {
	res := SanitizeKey("my_key")
	if res.Fixed != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %q", res.Fixed)
	}
	if !res.Changed {
		t.Error("expected Changed=true")
	}
}

func TestSanitizeKey_HyphenReplaced(t *testing.T) {
	res := SanitizeKey("my-key")
	if res.Fixed != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %q", res.Fixed)
	}
}

func TestSanitizeKey_DotReplaced(t *testing.T) {
	res := SanitizeKey("app.secret")
	if res.Fixed != "APP_SECRET" {
		t.Errorf("expected APP_SECRET, got %q", res.Fixed)
	}
}

func TestSanitizeKey_LeadingDigitPrefixed(t *testing.T) {
	res := SanitizeKey("1KEY")
	if res.Fixed != "_1KEY" {
		t.Errorf("expected _1KEY, got %q", res.Fixed)
	}
}

func TestSanitizeKey_InvalidCharsReplaced(t *testing.T) {
	res := SanitizeKey("key@name!")
	if res.Fixed != "KEY_NAME_" {
		t.Errorf("expected KEY_NAME_, got %q", res.Fixed)
	}
}

func TestSanitizeMap_BasicConversion(t *testing.T) {
	in := map[string]string{
		"db-host": "localhost",
		"db-port": "5432",
	}
	out, results := SanitizeMap(in)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost")
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestSanitizeMap_NoChanges(t *testing.T) {
	in := map[string]string{"VALID_KEY": "value"}
	out, results := SanitizeMap(in)
	if out["VALID_KEY"] != "value" {
		t.Error("expected VALID_KEY preserved")
	}
	if results[0].Changed {
		t.Error("expected no change reported")
	}
}
