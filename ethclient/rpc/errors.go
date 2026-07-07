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
	"errors"
	"fmt"
)

var (
	ErrNoResult             = errors.New("JSON-RPC response has no result")
	ErrMissingBatchResponse = errors.New("response batch did not contain a response to this call")
)

// HTTPError is returned by client operations when the HTTP status code of the
// response is not a 2xx status.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (err HTTPError) Error() string {
	if len(err.Body) == 0 {
		return err.Status
	}
	return fmt.Sprintf("%v: %s", err.Status, err.Body)
}

// RequestError wraps an error that occurred while performing a JSON-RPC HTTP
// request, adding the request URL, request body and (if available) response
// body to aid debugging.
type RequestError struct {
	URL          string
	RequestBody  string
	ResponseBody string
	Err          error
}

func (e *RequestError) Error() string {
	if len(e.ResponseBody) > 0 {
		return fmt.Sprintf("rpc request to %s failed: %v\nrequest: %s\nresponse: %s", e.URL, e.Err, e.RequestBody, e.ResponseBody)
	}
	return fmt.Sprintf("rpc request to %s failed: %v\nrequest: %s", e.URL, e.Err, e.RequestBody)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// Error wraps RPC errors, which contain an error code in addition to the message.
type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

// A DataError contains some data in addition to the error message.
type DataError interface {
	Error() string          // returns the message
	ErrorData() interface{} // returns the error data
}
