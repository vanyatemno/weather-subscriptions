package util

import (
	"strings"
	"testing"
)

func TestGenerateCodeLength(t *testing.T) {
	code, err := GenerateCode(6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("expected length 6, got %d", len(code))
	}
}

func TestGenerateCodeDigitsOnly(t *testing.T) {
	code, err := GenerateCode(12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Trim(code, "0123456789") != "" {
		t.Fatalf("expected digits only, got %q", code)
	}
}
