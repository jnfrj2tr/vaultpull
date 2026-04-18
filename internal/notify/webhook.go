package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to a webhook URL.
type WebhookPayload struct {
	Profile   string    `json:"profile"`
	Operation string    `json:"operation"`
	Added     int       `json:"added"`
	Updated   int       `json:"updated"`
	Removed   int       `json:"removed"`
	Timestamp time.Time `json:"timestamp"`
}

// WebhookNotifier posts event payloads to an HTTP endpoint.
type WebhookNotifier struct {
	URL    string
	client *http.Client
}

// NewWebhook creates a WebhookNotifier for the given URL.
func NewWebhook(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends the event as a JSON POST request.
func (w *WebhookNotifier) Notify(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	payload := WebhookPayload{
		Profile:   e.Profile,
		Operation: e.Operation,
		Added:     e.Added,
		Updated:   e.Updated,
		Removed:   e.Removed,
		Timestamp: e.Timestamp,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
