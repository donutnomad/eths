package rpc

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestUnmarshalBatchResponse_StandardArray(t *testing.T) {
	body := `[{"jsonrpc":"2.0","id":"1","result":"0x1"},{"jsonrpc":"2.0","id":"2","result":"0x2"}]`

	msgs, err := unmarshalBatchResponse(strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if string(msgs[0].ID) != `"1"` || string(msgs[1].ID) != `"2"` {
		t.Fatalf("unexpected message IDs: %q, %q", msgs[0].ID, msgs[1].ID)
	}
}

func TestUnmarshalBatchResponse_SingleErrorObjectFallback(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"batch rejected"}}`

	msgs, err := unmarshalBatchResponse(strings.NewReader(body))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Error == nil {
		t.Fatalf("expected message to carry an error")
	}
	if msgs[0].Error.Code != -32600 || msgs[0].Error.Message != "batch rejected" {
		t.Fatalf("unexpected error content: %+v", msgs[0].Error)
	}
}

func TestUnmarshalBatchResponse_SingleSuccessObjectNotFalledBack(t *testing.T) {
	// A single successful (non-error) object is not a valid batch response
	// and must not be silently accepted as one.
	body := `{"jsonrpc":"2.0","id":"1","result":"0x1"}`

	_, err := unmarshalBatchResponse(strings.NewReader(body))
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestUnmarshalBatchResponse_InvalidJSON(t *testing.T) {
	_, err := unmarshalBatchResponse(strings.NewReader("not json"))
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestUnmarshalBatchResponse_EmptyArray(t *testing.T) {
	msgs, err := unmarshalBatchResponse(strings.NewReader("[]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(msgs) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(msgs))
	}
}

func TestApplyBatchResponse_WholeBatchRejected(t *testing.T) {
	// Simulates a non-conforming server that rejects the entire batch with a
	// single error object carrying a null ID, instead of a per-request array.
	b := []BatchElem{
		{Method: "eth_call", Result: new(string)},
		{Method: "eth_call", Result: new(string)},
	}
	byID := map[string]int{`"1"`: 0, `"2"`: 1}
	batchresp := []*jsonrpcMessage{
		{Version: vsn, ID: json.RawMessage("null"), Error: &jsonError{Code: -32600, Message: "batch too large"}},
	}

	applyBatchResponse(b, byID, batchresp)

	for i, elem := range b {
		if elem.Error == nil {
			t.Fatalf("element %d: expected an error, got nil", i)
		}
		rpcErr, ok := elem.Error.(rpc.Error)
		if !ok {
			t.Fatalf("element %d: expected an rpc.Error, got %T", i, elem.Error)
		}
		if rpcErr.ErrorCode() != -32600 || !strings.Contains(rpcErr.Error(), "batch too large") {
			t.Fatalf("element %d: unexpected error content: %v", i, elem.Error)
		}
		if errors.Is(elem.Error, ErrMissingBatchResponse) {
			t.Fatalf("element %d: real error must not be masked by ErrMissingBatchResponse", i)
		}
	}
}

func TestApplyBatchResponse_StandardMatchByID(t *testing.T) {
	var r1, r2 string
	b := []BatchElem{
		{Method: "eth_call", Result: &r1},
		{Method: "eth_call", Result: &r2},
	}
	byID := map[string]int{`"1"`: 0, `"2"`: 1}
	batchresp := []*jsonrpcMessage{
		{Version: vsn, ID: json.RawMessage(`"2"`), Result: json.RawMessage(`"0x2"`)},
		{Version: vsn, ID: json.RawMessage(`"1"`), Result: json.RawMessage(`"0x1"`)},
	}

	applyBatchResponse(b, byID, batchresp)

	if b[0].Error != nil || r1 != "0x1" {
		t.Fatalf("element 0: unexpected result/error: %q, %v", r1, b[0].Error)
	}
	if b[1].Error != nil || r2 != "0x2" {
		t.Fatalf("element 1: unexpected result/error: %q, %v", r2, b[1].Error)
	}
}

func TestApplyBatchResponse_MissingResponseKeepsOthers(t *testing.T) {
	var r1 string
	b := []BatchElem{
		{Method: "eth_call", Result: &r1},
		{Method: "eth_call", Result: new(string)},
	}
	byID := map[string]int{`"1"`: 0, `"2"`: 1}
	batchresp := []*jsonrpcMessage{
		{Version: vsn, ID: json.RawMessage(`"1"`), Result: json.RawMessage(`"0x1"`)},
	}

	applyBatchResponse(b, byID, batchresp)

	if b[0].Error != nil || r1 != "0x1" {
		t.Fatalf("element 0: unexpected result/error: %q, %v", r1, b[0].Error)
	}
	if !errors.Is(b[1].Error, ErrMissingBatchResponse) {
		t.Fatalf("element 1: expected ErrMissingBatchResponse, got %v", b[1].Error)
	}
}
