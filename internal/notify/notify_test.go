package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	statusCode int
	captured   *http.Request
	body        string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.captured = req
	body := io.NopCloser(strings.NewReader(m.body))
	return &http.Response{StatusCode: m.statusCode, Body: body}, nil
}

func samplePayload() Payload {
	return Payload{
		Title:      "Drift Detected",
		Message:    "2 resources drifted",
		DriftCount: 2,
	}
}

func TestSend_WebhookSuccess(t *testing.T) {
	mc := &mockHTTPClient{statusCode: 200, body: "ok"}
	s := NewWithClient(Config{Channel: ChannelWebhook, WebhookURL: "http://example.com/hook"}, mc)
	if err := s.Send(samplePayload()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mc.captured == nil {
		t.Fatal("expected request to be made")
	}
	if mc.captured.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json content-type")
	}
}

func TestSend_SlackPayloadShape(t *testing.T) {
	mc := &mockHTTPClient{statusCode: 200, body: ""}
	s := NewWithClient(Config{Channel: ChannelSlack, WebhookURL: "http://slack.example.com"}, mc)
	if err := s.Send(samplePayload()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]interface{}
	body, _ := io.ReadAll(mc.captured.Body)
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if _, ok := got["text"]; !ok {
		t.Error("slack payload missing 'text' field")
	}
}

func TestSend_MissingWebhookURL(t *testing.T) {
	s := New(Config{Channel: ChannelWebhook})
	err := s.Send(samplePayload())
	if err == nil {
		t.Fatal("expected error for missing webhook_url")
	}
}

func TestSend_Non2xxStatus(t *testing.T) {
	mc := &mockHTTPClient{statusCode: 500, body: ""}
	s := NewWithClient(Config{Channel: ChannelWebhook, WebhookURL: "http://example.com"}, mc)
	err := s.Send(samplePayload())
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}

func TestSend_UnsupportedChannel(t *testing.T) {
	mc := &mockHTTPClient{statusCode: 200, body: ""}
	s := NewWithClient(Config{Channel: "pagerduty", WebhookURL: "http://example.com"}, mc)
	err := s.Send(samplePayload())
	if err == nil {
		t.Fatal("expected error for unsupported channel")
	}
}

func TestSend_CustomHeaders(t *testing.T) {
	mc := &mockHTTPClient{statusCode: 204, body: ""}
	cfg := Config{
		Channel:    ChannelWebhook,
		WebhookURL: "http://example.com",
		Headers:    map[string]string{"X-Token": "secret"},
	}
	s := NewWithClient(cfg, mc)
	if err := s.Send(samplePayload()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mc.captured.Header.Get("X-Token") != "secret" {
		t.Error("expected custom header X-Token to be set")
	}
}
