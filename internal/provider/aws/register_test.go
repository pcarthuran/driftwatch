package aws_test

import (
	"testing"

	"github.com/user/driftwatch/internal/provider"
	"github.com/user/driftwatch/internal/provider/aws"
)

func TestRegister_Success(t *testing.T) {
	reg := provider.NewRegistry()
	err := aws.Register(reg, aws.Config{Region: "eu-west-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := reg.Get("aws")
	if !ok {
		t.Fatal("expected aws provider to be registered")
	}
	if p.Name() != "aws" {
		t.Errorf("expected name 'aws', got %q", p.Name())
	}
}

func TestRegister_MissingRegion(t *testing.T) {
	reg := provider.NewRegistry()
	err := aws.Register(reg, aws.Config{})
	if err == nil {
		t.Fatal("expected error for missing region, got nil")
	}
}

func TestRegister_Duplicate(t *testing.T) {
	reg := provider.NewRegistry()
	if err := aws.Register(reg, aws.Config{Region: "us-west-2"}); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	err := aws.Register(reg, aws.Config{Region: "us-east-1"})
	if err == nil {
		t.Fatal("expected error on duplicate registration, got nil")
	}
}
