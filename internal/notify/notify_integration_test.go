package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/internal/notify"
)

func TestIntegration_WebhookReceivesPayload(t *testing.T) {
	var received notify.Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := notify.Config{
		Channel:    notify.ChannelWebhook,
		WebhookURL: ts.URL,
	}
	sender := notify.New(cfg)

	p := notify.Payload{
		Title:      "Integration Test",
		Message:    "3 resources drifted",
		DriftCount: 3,
	}
	if err := sender.Send(p); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if received.Title != p.Title {
		t.Errorf("title: got %q, want %q", received.Title, p.Title)
	}
	if received.DriftCount != p.DriftCount {
		t.Errorf("drift_count: got %d, want %d", received.DriftCount, p.DriftCount)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestIntegration_SlackPayloadSent(t *testing.T) {
	var raw map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&raw)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := notify.Config{
		Channel:    notify.ChannelSlack,
		WebhookURL: ts.URL,
	}
	sender := notify.New(cfg)

	if err := sender.Send(notify.Payload{Title: "Slack Test", Message: "ok", DriftCount: 0}); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if _, ok := raw["text"]; !ok {
		t.Error("slack payload missing 'text' key")
	}
}

func TestIntegration_ServerError_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sender := notify.New(notify.Config{Channel: notify.ChannelWebhook, WebhookURL: ts.URL})
	err := sender.Send(notify.Payload{Title: "t", Message: "m", DriftCount: 1})
	if err == nil {
		t.Fatal("expected error from 500 response")
	}
}

func TestIntegration_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Use a very short timeout client to trigger timeout behaviour.
	client := &http.Client{Timeout: 1 * time.Millisecond}
	sender := notify.NewWithClient(
		notify.Config{Channel: notify.ChannelWebhook, WebhookURL: ts.URL},
		client,
	)
	err := sender.Send(notify.Payload{Title: "t", Message: "m", DriftCount: 0})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
