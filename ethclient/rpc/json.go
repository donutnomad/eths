package rpc

import (
	"encoding/json"
	"fmt"
)

const vsn = "2.0"

// A value of this type can a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type jsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func (msg *jsonrpcMessage) isNotification() bool {
	return msg.hasValidVersion() && msg.ID == nil && msg.Method != ""
}

func (msg *jsonrpcMessage) isCall() bool {
	return msg.hasValidVersion() && msg.hasValidID() && msg.Method != ""
}

func (msg *jsonrpcMessage) isResponse() bool {
	return msg.hasValidVersion() && msg.hasValidID() && msg.Method == "" && msg.Params == nil && (msg.Result != nil || msg.Error != nil)
}

func (msg *jsonrpcMessage) hasValidID() bool {
	return len(msg.ID) > 0 && msg.ID[0] != '{' && msg.ID[0] != '['
}

func (msg *jsonrpcMessage) hasValidVersion() bool {
	return msg.Version == vsn
}

func (msg *jsonrpcMessage) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

func (msg *jsonrpcMessage) errorResponse(err error) *jsonrpcMessage {
	resp := errorMessage(err)
	resp.ID = msg.ID
	return resp
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err *jsonError) Error() string {
	msg := err.Message
	if msg == "" {
		msg = fmt.Sprintf("json-rpc error %d", err.Code)
	}
	if err.Data == nil {
		return fmt.Sprintf("code: %d, message: %s", err.Code, msg)
	}
	data, dataErr := json.Marshal(err.Data)
	if dataErr != nil {
		return fmt.Sprintf("code: %d, message: %s, data: %v", err.Code, msg, err.Data)
	}
	return fmt.Sprintf("code: %d, message: %s, data: %s", err.Code, msg, data)
}

//
//func (err *jsonError) Error() string {
//	if err.Message == "" {
//		return fmt.Sprintf("json-rpc error %d", err.Code)
//	}
//	return err.Message
//}

func (err *jsonError) ErrorCode() int {
	return err.Code
}

func (err *jsonError) ErrorData() any {
	return err.Data
}

// hasValidID reports whether id is a usable JSON-RPC request/response ID,
// i.e. not absent and not JSON null.
func hasValidID(id json.RawMessage) bool {
	return len(id) > 0 && string(id) != "null"
}
