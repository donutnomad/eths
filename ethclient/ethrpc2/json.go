// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

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

func (msg *jsonrpcMessage) String() string {
	b, _ := json.Marshal(msg)
	return string(b)
}

// hasValidID reports whether id is a usable JSON-RPC request/response ID,
// i.e. not absent and not JSON null.
func hasValidID(id json.RawMessage) bool {
	return len(id) > 0 && string(id) != "null"
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

func (err *jsonError) ErrorCode() int {
	return err.Code
}

func (err *jsonError) ErrorData() interface{} {
	return err.Data
}
