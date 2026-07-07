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
	"context"
	"net/http"
)

// ClientOption is a configuration option for the RPC client.
type ClientOption interface {
	applyOption(*clientConfig)
}

type clientConfig struct {
	// HTTP settings
	httpClient  *http.Client
	httpHeaders http.Header
	httpAuth    HTTPAuth
}

func (cfg *clientConfig) initHeaders() {
	if cfg.httpHeaders == nil {
		cfg.httpHeaders = make(http.Header)
	}
}

func (cfg *clientConfig) setHeader(key, value string) {
	cfg.initHeaders()
	cfg.httpHeaders.Set(key, value)
}

type optionFunc func(*clientConfig)

func (fn optionFunc) applyOption(opt *clientConfig) {
	fn(opt)
}

// WithHeader configures HTTP headers set by the RPC client.
func WithHeader(key, value string) ClientOption {
	return optionFunc(func(cfg *clientConfig) {
		cfg.initHeaders()
		cfg.httpHeaders.Set(key, value)
	})
}

// WithHeaders configures HTTP headers set by the RPC client.
func WithHeaders(headers http.Header) ClientOption {
	return optionFunc(func(cfg *clientConfig) {
		cfg.initHeaders()
		for k, vs := range headers {
			cfg.httpHeaders[k] = vs
		}
	})
}

// WithHTTPClient configures the http.Client used by the RPC client.
func WithHTTPClient(c *http.Client) ClientOption {
	return optionFunc(func(cfg *clientConfig) {
		cfg.httpClient = c
	})
}

// WithHTTPAuth configures HTTP request authentication. The given provider will be called
// whenever a request is made. Note that only one authentication provider can be active at
// any time.
func WithHTTPAuth(a HTTPAuth) ClientOption {
	if a == nil {
		panic("nil auth")
	}
	return optionFunc(func(cfg *clientConfig) {
		cfg.httpAuth = a
	})
}

// A HTTPAuth function is called by the client whenever a HTTP request is sent.
// The function must be safe for concurrent use.
//
// Usually, HTTPAuth functions will call h.Set("authorization", "...") to add
// auth information to the request.
type HTTPAuth func(h http.Header) error

type mdHeaderKey struct{}

// NewContextWithHeaders wraps the given context, adding HTTP headers. These headers will
// be applied by Client when making a request using the returned context.
func NewContextWithHeaders(ctx context.Context, h http.Header) context.Context {
	if len(h) == 0 {
		// This check ensures the header map set in context will never be nil.
		return ctx
	}

	var ctxh http.Header
	prev, ok := ctx.Value(mdHeaderKey{}).(http.Header)
	if ok {
		ctxh = setHeaders(prev.Clone(), h)
	} else {
		ctxh = h.Clone()
	}
	return context.WithValue(ctx, mdHeaderKey{}, ctxh)
}

// headersFromContext is used to extract http.Header from context.
func headersFromContext(ctx context.Context) http.Header {
	source, _ := ctx.Value(mdHeaderKey{}).(http.Header)
	return source
}

// setHeaders sets all headers from src in dst.
func setHeaders(dst http.Header, src http.Header) http.Header {
	for key, values := range src {
		dst[http.CanonicalHeaderKey(key)] = values
	}
	return dst
}
