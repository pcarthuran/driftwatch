package azure_test

import (
	"testing"

	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/provider/azure"
)

func TestRegister_Success(t *testing.T) {
	registry := provider.NewRegistry()
	cfg := map[string]string{
		"subscription_id": "sub-abc-123",
	}

	if err := azure.Register(registry, cfg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	p, err := registry.Get(azure.ProviderName)
	if err != nil {
		t.Fatalf("expected provider to be registered, got: %v", err)
	}
	if p.Name() != azure.ProviderName {
		t.Errorf("expected provider name %q, got %q", azure.ProviderName, p.Name())
	}
}

func TestRegister_MissingSubscriptionID(t *testing.T) {
	registry := provider.NewRegistry()

	tests := []struct {
		name string
		cfg  map[string]string
	}{
		{"empty map", map[string]string{}},
		{"empty value", map[string]string{"subscription_id": ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := azure.Register(registry, tt.cfg); err == nil {
				t.Error("expected error for missing subscription_id, got nil")
			}
		})
	}
}

func TestRegister_Duplicate(t *testing.T) {
	registry := provider.NewRegistry()
	cfg := map[string]string{
		"subscription_id": "sub-abc-123",
	}

	if err := azure.Register(registry, cfg); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	if err := azure.Register(registry, cfg); err == nil {
		t.Error("expected error on duplicate registration, got nil")
	}
}
