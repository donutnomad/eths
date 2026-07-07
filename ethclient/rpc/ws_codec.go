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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"
)

// ServerCodec implements reading, parsing and writing RPC messages received over the
// connection. Implementations must be go-routine safe since the codec can be called in
// multiple go-routines concurrently.
type ServerCodec interface {
	readBatch() (msgs []*jsonrpcMessage, isBatch bool, err error)
	close()

	jsonWriter
}

// jsonWriter can write JSON messages to its underlying connection.
// Implementations must be safe for concurrent use.
type jsonWriter interface {
	// writeJSON writes a message to the connection.
	writeJSON(ctx context.Context, msg interface{}, isError bool) error

	// Closed returns a channel which is closed when the connection is closed.
	closed() <-chan interface{}
	// RemoteAddr returns the peer address of the connection.
	remoteAddr() string
}

const (
	subscribeMethodSuffix    = "_subscribe"
	unsubscribeMethodSuffix  = "_unsubscribe"
	notificationMethodSuffix = "_subscription"

	defaultWriteTimeout = 10 * time.Second // used if context has no deadline
)

var null = json.RawMessage("null")

type subscriptionResult struct {
	ID     string          `json:"subscription"`
	Result json.RawMessage `json:"result,omitempty"`
}

func errorMessage(err error) *jsonrpcMessage {
	msg := &jsonrpcMessage{Version: vsn, ID: null, Error: &jsonError{
		Code:    errcodeDefault,
		Message: err.Error(),
	}}
	ec, ok := err.(Error)
	if ok {
		msg.Error.Code = ec.ErrorCode()
	}
	de, ok := err.(DataError)
	if ok {
		msg.Error.Data = de.ErrorData()
	}
	return msg
}

type deadlineCloser interface {
	io.Closer
	SetWriteDeadline(time.Time) error
}

// ConnRemoteAddr wraps the RemoteAddr operation, which returns a description
// of the peer address of a connection. If a conn also implements ConnRemoteAddr, this
// description is used in log messages.
type ConnRemoteAddr interface {
	RemoteAddr() string
}

// jsonCodec reads and writes JSON-RPC messages to the underlying connection. It also has
// support for parsing arguments and serializing (result) objects.
type jsonCodec struct {
	remote  string
	closer  sync.Once  // close closed channel once
	closeCh chan any   // closed on Close
	decode  decodeFunc // decoder to allow multiple transports
	encMu   sync.Mutex // guards the encoder
	encode  encodeFunc // encoder to allow multiple transports
	conn    deadlineCloser
}

type encodeFunc = func(v any, isErrorResponse bool) error

type decodeFunc = func(v any) error

// NewFuncCodec creates a codec which uses the given functions to read and write. If conn
// implements ConnRemoteAddr, log messages will use it to include the remote address of
// the connection.
func NewFuncCodec(conn deadlineCloser, encode encodeFunc, decode decodeFunc) ServerCodec {
	codec := &jsonCodec{
		closeCh: make(chan any),
		encode:  encode,
		decode:  decode,
		conn:    conn,
	}
	if ra, ok := conn.(ConnRemoteAddr); ok {
		codec.remote = ra.RemoteAddr()
	}
	return codec
}

func (c *jsonCodec) remoteAddr() string {
	return c.remote
}

func (c *jsonCodec) readBatch() (messages []*jsonrpcMessage, batch bool, err error) {
	// Decode the next JSON object in the input stream.
	// This verifies basic syntax, etc.
	var rawmsg json.RawMessage
	if err := c.decode(&rawmsg); err != nil {
		return nil, false, err
	}
	messages, batch = parseMessage(rawmsg)
	for i, msg := range messages {
		if msg == nil {
			// Message is JSON 'null'. Replace with zero value so it
			// will be treated like any other invalid message.
			messages[i] = new(jsonrpcMessage)
		}
	}
	return messages, batch, nil
}

func (c *jsonCodec) writeJSON(ctx context.Context, v any, isErrorResponse bool) error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(defaultWriteTimeout)
	}
	c.conn.SetWriteDeadline(deadline)
	return c.encode(v, isErrorResponse)
}

func (c *jsonCodec) close() {
	c.closer.Do(func() {
		close(c.closeCh)
		c.conn.Close()
	})
}

// closed returns a channel which will be closed when Close is called
func (c *jsonCodec) closed() <-chan any {
	return c.closeCh
}

// parseMessage parses raw bytes as a (batch of) JSON-RPC message(s). There are no error
// checks in this function because the raw message has already been syntax-checked when it
// is called. Any non-JSON-RPC messages in the input return the zero value of
// jsonrpcMessage.
func parseMessage(raw json.RawMessage) ([]*jsonrpcMessage, bool) {
	if !isBatch(raw) {
		msgs := []*jsonrpcMessage{{}}
		json.Unmarshal(raw, &msgs[0])
		return msgs, false
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.Token() // skip '['
	var msgs []*jsonrpcMessage
	for dec.More() {
		msgs = append(msgs, new(jsonrpcMessage))
		dec.Decode(&msgs[len(msgs)-1])
	}
	return msgs, true
}

// isBatch returns true when the first non-whitespace characters is '['
func isBatch(raw json.RawMessage) bool {
	for _, c := range raw {
		// skip insignificant whitespace (http://www.ietf.org/rfc/rfc4627.txt)
		if c == 0x20 || c == 0x09 || c == 0x0a || c == 0x0d {
			continue
		}
		return c == '['
	}
	return false
}
