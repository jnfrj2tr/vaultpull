package env

import (
	"testing"
)

var sampleSecrets = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
}

func TestFilter_IncludePrefix(t *testing.T) {
	result, err := Filter(sampleSecrets, FilterOptions{IncludePrefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
}

func TestFilter_ExcludePrefix(t *testing.T) {
	result, err := Filter(sampleSecrets, FilterOptions{ExcludePrefix: "DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded")
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
}

func TestFilter_Pattern(t *testing.T) {
	result, err := Filter(sampleSecrets, FilterOptions{Pattern: "_HOST$"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := Filter(sampleSecrets, FilterOptions{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestFilter_NoOptions(t *testing.T) {
	result, err := Filter(sampleSecrets, FilterOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(sampleSecrets) {
		t.Fatalf("expected all keys, got %d", len(result))
	}
}

func TestFilter_CombinedPrefixAndPattern(t *testing.T) {
	result, err := Filter(sampleSecrets, FilterOptions{IncludePrefix: "DB_", Pattern: "PASSWORD"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD")
	}
}
