package rpc

import (
	"strconv"
	"strings"
	"testing"
)

func TestJsonError_Error_IncludesCodeMessageAndData(t *testing.T) {
	err := &jsonError{
		Code:    -32602,
		Message: "invalid argument 0",
		Data:    map[string]any{"argument": 0, "reason": "not a valid address"},
	}

	s := err.Error()
	if !strings.Contains(s, strconv.Itoa(err.Code)) {
		t.Fatalf("expected error string to contain code %d, got %q", err.Code, s)
	}
	if !strings.Contains(s, err.Message) {
		t.Fatalf("expected error string to contain message %q, got %q", err.Message, s)
	}
	if !strings.Contains(s, "not a valid address") {
		t.Fatalf("expected error string to contain data content, got %q", s)
	}
}

func TestJsonError_Error_NoData(t *testing.T) {
	err := &jsonError{Code: -32601, Message: "method not found"}

	s := err.Error()
	if !strings.Contains(s, strconv.Itoa(err.Code)) || !strings.Contains(s, err.Message) {
		t.Fatalf("expected error string to contain code and message, got %q", s)
	}
	if strings.Contains(s, "data:") {
		t.Fatalf("expected no data section when Data is nil, got %q", s)
	}
}

func TestJsonError_Error_NoMessageFallsBackToCode(t *testing.T) {
	err := &jsonError{Code: -32000}

	s := err.Error()
	if !strings.Contains(s, strconv.Itoa(err.Code)) {
		t.Fatalf("expected error string to contain code %d, got %q", err.Code, s)
	}
}
