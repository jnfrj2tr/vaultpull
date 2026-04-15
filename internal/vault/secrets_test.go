package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeKVv2Response(data map[string]interface{}) []byte {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
		},
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestGetSecrets_ReturnsStringValues(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(makeKVv2Response(map[string]interface{}{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		}))
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	secrets, err := client.GetSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %s", secrets["DB_PORT"])
	}
}

func TestGetSecrets_NonStringValueCoerced(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(makeKVv2Response(map[string]interface{}{
			"RETRIES": 3,
		}))
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	secrets, err := client.GetSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if secrets["RETRIES"] != "3" {
		t.Errorf("expected RETRIES=3, got %s", secrets["RETRIES"])
	}
}

func TestGetSecrets_404ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := client.GetSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetSecrets_UnexpectedStatusReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := client.GetSecrets("secret/data/forbidden")
	if err == nil {
		t.Fatal("expected error for 403, got nil")
	}
}
