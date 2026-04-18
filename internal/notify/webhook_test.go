package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhook_PostsJSON(t *testing.T) {
	var received WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wn := NewWebhook(ts.URL)
	e := Event{Profile: "prod", Operation: "rotate", Added: 3, Timestamp: time.Now()}
	if err := wn.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Profile != "prod" {
		t.Errorf("expected profile prod, got %s", received.Profile)
	}
	if received.Added != 3 {
		t.Errorf("expected added=3, got %d", received.Added)
	}
}

func TestWebhook_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wn := NewWebhook(ts.URL)
	err := wn.Notify(Event{Profile: "dev", Operation: "sync"})
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestWebhook_InvalidURLReturnsError(t *testing.T) {
	wn := NewWebhook("http://127.0.0.1:0/nope")
	err := wn.Notify(Event{Profile: "dev", Operation: "sync"})
	if err == nil {
		t.Error("expected connection error")
	}
}
