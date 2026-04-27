package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Channel represents a notification destination type.
type Channel string

const (
	ChannelSlack   Channel = "slack"
	ChannelWebhook Channel = "webhook"
)

// Config holds notification settings.
type Config struct {
	Channel    Channel           `json:"channel" yaml:"channel"`
	WebhookURL string            `json:"webhook_url" yaml:"webhook_url"`
	Headers    map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// Payload is the message sent to a notification channel.
type Payload struct {
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	DriftCount int      `json:"drift_count"`
	Timestamp time.Time `json:"timestamp"`
}

// Sender dispatches notifications.
type Sender struct {
	cfg    Config
	client HTTPClient
}

// HTTPClient is an interface for making HTTP requests (enables testing).
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// New creates a Sender with the given config.
func New(cfg Config) *Sender {
	return &Sender{cfg: cfg, client: &http.Client{Timeout: 10 * time.Second}}
}

// NewWithClient creates a Sender with a custom HTTP client.
func NewWithClient(cfg Config, client HTTPClient) *Sender {
	return &Sender{cfg: cfg, client: client}
}

// Send dispatches the payload to the configured channel.
func (s *Sender) Send(p Payload) error {
	if s.cfg.WebhookURL == "" {
		return fmt.Errorf("notify: webhook_url is required")
	}
	p.Timestamp = time.Now().UTC()

	var body []byte
	var err error

	switch s.cfg.Channel {
	case ChannelSlack:
		body, err = buildSlackPayload(p)
	case ChannelWebhook, "":
		body, err = json.Marshal(p)
	default:
		return fmt.Errorf("notify: unsupported channel %q", s.cfg.Channel)
	}
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.cfg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range s.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("notify: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func buildSlackPayload(p Payload) ([]byte, error) {
	slack := map[string]interface{}{
		"text": fmt.Sprintf("*%s*\n%s\nDrift count: %d", p.Title, p.Message, p.DriftCount),
	}
	return json.Marshal(slack)
}
