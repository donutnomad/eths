// Copyright 2022 The go-ethereum Authors
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
	"net/http"

	"github.com/gorilla/websocket"
)

// WsClientOption is a configuration option for the RPC client.
type WsClientOption interface {
	applyOption(*wsclientConfig)
}

type wsclientConfig struct {
	httpAuth    HTTPAuth
	httpHeaders http.Header
	// WebSocket options
	wsDialer           *websocket.Dialer
	wsMessageSizeLimit *int64 // wsMessageSizeLimit nil = default, 0 = no limit
}

func (cfg *wsclientConfig) initHeaders() {
	if cfg.httpHeaders == nil {
		cfg.httpHeaders = make(http.Header)
	}
}

func (cfg *wsclientConfig) setHeader(key, value string) {
	cfg.initHeaders()
	cfg.httpHeaders.Set(key, value)
}

type wsOptionFunc func(*wsclientConfig)

func (fn wsOptionFunc) applyOption(opt *wsclientConfig) {
	fn(opt)
}

// WithWebsocketDialer configures the websocket.Dialer used by the RPC client.
func WithWebsocketDialer(dialer websocket.Dialer) WsClientOption {
	return wsOptionFunc(func(cfg *wsclientConfig) {
		cfg.wsDialer = &dialer
	})
}

// WithWebsocketMessageSizeLimit configures the websocket message size limit used by the RPC
// client. Passing a limit of 0 means no limit.
func WithWebsocketMessageSizeLimit(messageSizeLimit int64) WsClientOption {
	return wsOptionFunc(func(cfg *wsclientConfig) {
		cfg.wsMessageSizeLimit = &messageSizeLimit
	})
}

// WithWebsocketHTTPAuth configures HTTP request authentication. The given provider will be called
// whenever a request is made. Note that only one authentication provider can be active at
// any time.
func WithWebsocketHTTPAuth(a HTTPAuth) WsClientOption {
	if a == nil {
		panic("nil auth")
	}
	return wsOptionFunc(func(cfg *wsclientConfig) {
		cfg.httpAuth = a
	})
}
