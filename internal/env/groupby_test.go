package env

import (
	"testing"
)

func TestGroupBy_Prefix(t *testing.T) {
	input := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
		"NODELIM": "value",
	}
	result, err := GroupBy(input, GroupByOptions{Mode: GroupByPrefix})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(result["DB"]))
	}
	if len(result["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(result["APP"]))
	}
	if _, ok := result["other"]["NODELIM"]; !ok {
		t.Errorf("expected NODELIM in fallback group")
	}
}

func TestGroupBy_PrefixCustomDelimiter(t *testing.T) {
	input := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
		"app.env": "prod",
	}
	result, err := GroupBy(input, GroupByOptions{Mode: GroupByPrefix, Delimiter: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(result["db"]))
	}
}

func TestGroupBy_Delimiter(t *testing.T) {
	input := map[string]string{
		"aws_us_east_ACCESS_KEY": "AKIA",
		"aws_us_east_SECRET_KEY": "secret",
		"standalone":             "value",
	}
	result, err := GroupBy(input, GroupByOptions{Mode: GroupByDelimiter, Delimiter: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["aws_us_east"]; !ok {
		t.Errorf("expected group aws_us_east")
	}
	if _, ok := result["other"]["standalone"]; !ok {
		t.Errorf("expected standalone in fallback")
	}
}

func TestGroupBy_DelimiterMissingDelimiter(t *testing.T) {
	_, err := GroupBy(map[string]string{"KEY": "val"}, GroupByOptions{Mode: GroupByDelimiter})
	if err == nil {
		t.Error("expected error for missing delimiter")
	}
}

func TestGroupBy_Pattern(t *testing.T) {
	input := map[string]string{
		"PROD_API_KEY":     "abc",
		"PROD_DB_PASSWORD": "secret",
		"DEV_API_KEY":      "xyz",
		"UNMATCHED":        "val",
	}
	result, err := GroupBy(input, GroupByOptions{
		Mode:    GroupByPattern,
		Pattern: `^(PROD|DEV)_`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result["PROD"]) != 2 {
		t.Errorf("expected 2 PROD keys, got %d", len(result["PROD"]))
	}
	if len(result["DEV"]) != 1 {
		t.Errorf("expected 1 DEV key, got %d", len(result["DEV"]))
	}
	if _, ok := result["other"]["UNMATCHED"]; !ok {
		t.Errorf("expected UNMATCHED in fallback group")
	}
}

func TestGroupBy_PatternInvalid(t *testing.T) {
	_, err := GroupBy(map[string]string{"K": "v"}, GroupByOptions{
		Mode:    GroupByPattern,
		Pattern: `[invalid`,
	})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestGroupBy_UnknownMode(t *testing.T) {
	_, err := GroupBy(map[string]string{"K": "v"}, GroupByOptions{Mode: "bogus"})
	if err == nil {
		t.Error("expected error for unknown mode")
	}
}

func TestGroupBy_CustomFallback(t *testing.T) {
	input := map[string]string{"NOPREFIX": "val"}
	result, err := GroupBy(input, GroupByOptions{Mode: GroupByPrefix, Fallback: "misc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["misc"]["NOPREFIX"]; !ok {
		t.Errorf("expected NOPREFIX in custom fallback group 'misc'")
	}
}
