package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestMerge_OverwritesExisting(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=old\nDB_PORT=5432\n")

	result, err := Merge(path, map[string]string{"DB_HOST": "newhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "newhost" {
		t.Errorf("expected DB_HOST=newhost, got %s", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432 preserved, got %s", result["DB_PORT"])
	}
}

func TestMerge_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")

	result, err := Merge(path, map[string]string{"NEW_KEY": "value"})
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if result["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %s", result["NEW_KEY"])
	}
}

func TestReadEnvFile_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=val\n\nANOTHER=123\n")

	result, err := readEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(result), result)
	}
}
