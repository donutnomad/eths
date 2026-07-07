// Copyright 2016 The go-ethereum Authors
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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// DialHTTP creates a new RPC client for the given URL.
func DialHTTP(endpoint string, options ...ClientOption) (*HttpClient, error) {
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	cfg := new(clientConfig)
	for _, opt := range options {
		opt.applyOption(cfg)
	}
	headers := make(http.Header, 2+len(cfg.httpHeaders))
	headers.Set("accept", "application/json")
	headers.Set("content-type", "application/json")
	for key, values := range cfg.httpHeaders {
		headers[key] = values
	}

	client := cfg.httpClient
	if client == nil {
		client = new(http.Client)
	}

	return &HttpClient{
		client:  client,
		headers: headers,
		url:     endpoint,
		auth:    cfg.httpAuth,
		closeCh: make(chan any),
	}, nil
}

// HttpClient represents a connection to an RPC server over HTTP.
type HttpClient struct {
	idCounter atomic.Uint32
	client    *http.Client
	url       string
	closeOnce sync.Once
	closeCh   chan any
	mu        sync.Mutex // protects headers
	headers   http.Header
	auth      HTTPAuth
}

// SupportedModules calls the rpc_modules method, retrieving the list of
// APIs that are available on the server.
func (c *HttpClient) SupportedModules() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result map[string]string
	err := c.Call(ctx, &result, "rpc_modules")
	return result, err
}

// Close closes the client, aborting any in-flight requests.
func (c *HttpClient) Close() {
	c.closeOnce.Do(func() { close(c.closeCh) })
}

// SetHeader adds a custom HTTP header to the client's requests.
func (c *HttpClient) SetHeader(key, value string) {
	c.mu.Lock()
	c.headers.Set(key, value)
	c.mu.Unlock()
}

// Call performs a JSON-RPC call with the given arguments. If the context is
// canceled before the call has successfully returned, CallContext returns immediately.
//
// The result must be a pointer so that package json can unmarshal into it. You
// can also pass nil, in which case the result is ignored.
func (c *HttpClient) Call(ctx context.Context, result any, method string, args ...any) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}
	msg, err := c.newMessage(method, args...)
	if err != nil {
		return err
	}
	resp, err := c.sendHTTP(ctx, msg)
	if err != nil {
		return err
	} else if resp.Error != nil {
		return resp.Error
	} else if len(resp.Result) == 0 {
		return ErrNoResult
	} else if result == nil {
		return nil
	}
	return json.Unmarshal(resp.Result, result)
}

// BatchCall sends all given requests as a single batch and waits for the server
// to return a response for all of them. The wait duration is bounded by the
// context's deadline.
//
// In contrast to CallContext, BatchCallContext only returns errors that have occurred
// while sending the request. Any error specific to a request is reported through the
// Error field of the corresponding BatchElem.
//
// Note that batch calls may not be executed atomically on the server side.
func (c *HttpClient) BatchCall(ctx context.Context, b []BatchElem) error {
	var (
		msgs = make([]*jsonrpcMessage, len(b))
		byID = make(map[string]int, len(b))
	)
	for i, elem := range b {
		msg, err := c.newMessage(elem.Method, elem.Args...)
		if err != nil {
			return err
		}
		msgs[i] = msg
		byID[string(msg.ID)] = i
	}

	batchresp, err := c.sendBatchHTTP(ctx, msgs)
	if err != nil {
		return err
	}

	applyBatchResponse(b, byID, batchresp)
	return nil
}

// applyBatchResponse matches each message in batchresp to the BatchElem that
// requested it (via byID) and assigns Result/Error accordingly. Elements
// that never got a matching response are left with ErrMissingBatchResponse.
//
// Some non-conforming servers reject the whole batch upfront and respond
// with a single error message that isn't tied to any request ID (see
// unmarshalBatchResponse's fallback). In that case there is no point
// matching it against individual elements by ID - the actual error is
// propagated to every element instead of being masked by
// ErrMissingBatchResponse.
func applyBatchResponse(b []BatchElem, byID map[string]int, batchresp []*jsonrpcMessage) {
	if len(batchresp) == 1 && batchresp[0] != nil && batchresp[0].Error != nil && !hasValidID(batchresp[0].ID) {
		batchErr := batchresp[0].Error
		for i := range b {
			b[i].Error = batchErr
		}
		return
	}

	for _, resp := range batchresp {
		if resp == nil {
			// Ignore null responses. These can happen for batches sent via HTTP.
			continue
		}

		// Find the element corresponding to this response.
		index, ok := byID[string(resp.ID)]
		if !ok {
			continue
		}
		delete(byID, string(resp.ID))

		// Assign result and error.
		elem := &b[index]
		switch {
		case resp.Error != nil:
			elem.Error = resp.Error
		case resp.Result == nil:
			elem.Error = ErrNoResult
		default:
			elem.Error = json.Unmarshal(resp.Result, elem.Result)
		}
	}

	// Check that all expected responses have been received.
	for _, index := range byID {
		elem := &b[index]
		elem.Error = ErrMissingBatchResponse
	}
}

// Notify sends a notification, i.e. a method call that doesn't expect a response.
func (c *HttpClient) Notify(ctx context.Context, method string, args ...any) error {
	msg, err := c.newMessage(method, args...)
	if err != nil {
		return err
	}
	msg.ID = nil

	_, err = c.sendHTTP(ctx, msg)
	return err
}

func (c *HttpClient) newMessage(method string, paramsIn ...any) (*jsonrpcMessage, error) {
	msg := &jsonrpcMessage{Version: vsn, ID: c.nextID(), Method: method}
	if paramsIn != nil { // prevent sending "params":null
		var err error
		if msg.Params, err = json.Marshal(paramsIn); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func (c *HttpClient) nextID() json.RawMessage {
	id := c.idCounter.Add(1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (c *HttpClient) doRequest(ctx context.Context, body []byte) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, &RequestError{URL: c.url, RequestBody: string(body), Err: err}
	}
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	// set headers
	c.mu.Lock()
	req.Header = c.headers.Clone()
	c.mu.Unlock()
	setHeaders(req.Header, headersFromContext(ctx))

	if c.auth != nil {
		if err := c.auth(req.Header); err != nil {
			return nil, &RequestError{URL: c.url, RequestBody: string(body), Err: err}
		}
	}

	// do request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, &RequestError{URL: c.url, RequestBody: string(body), Err: err}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var respBody []byte
		if _, err := buf.ReadFrom(resp.Body); err == nil {
			respBody = buf.Bytes()
		}
		cleanlyCloseBody(resp.Body)
		return nil, &RequestError{
			URL:          c.url,
			RequestBody:  string(body),
			ResponseBody: string(respBody),
			Err: HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				Body:       respBody,
			},
		}
	}

	return resp.Body, nil
}

// peekBufSize is the number of leading response bytes kept around so decode
// errors can include a preview of what the server actually sent.
const peekBufSize = 300

func (c *HttpClient) sendHTTP(ctx context.Context, msg any) (*jsonrpcMessage, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	respBody, err := c.doRequest(ctx, body)
	if err != nil {
		return nil, err
	}
	defer cleanlyCloseBody(respBody)

	peek := newPeekReader(respBody, peekBufSize)
	var resp jsonrpcMessage
	if err := json.NewDecoder(peek).Decode(&resp); err != nil {
		return nil, &RequestError{URL: c.url, RequestBody: string(body), ResponseBody: peek.buffered(), Err: err}
	}
	return &resp, nil
}

func (c *HttpClient) sendBatchHTTP(ctx context.Context, msg any) ([]*jsonrpcMessage, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	respBody, err := c.doRequest(ctx, body)
	if err != nil {
		return nil, err
	}
	defer cleanlyCloseBody(respBody)

	peek := newPeekReader(respBody, peekBufSize)
	respmsgs, decodeErr := unmarshalBatchResponse(peek)
	if decodeErr != nil {
		return nil, &RequestError{URL: c.url, RequestBody: string(body), ResponseBody: peek.buffered(), Err: decodeErr}
	}
	return respmsgs, nil
}

// unmarshalBatchResponse decodes a batch JSON-RPC response, which is
// expected to be a JSON array of messages (`[{...}, {...}]`). Some
// non-conforming servers instead respond with a single JSON-RPC error object
// (`{...}`) when the whole batch is rejected upfront; this is accepted as a
// fallback and wrapped into a single-element slice. The response is decoded
// incrementally rather than buffered in full, so this is safe to use on
// arbitrarily large response bodies.
func unmarshalBatchResponse(r io.Reader) ([]*jsonrpcMessage, error) {
	br := bufio.NewReader(r)
	isObject, err := peekIsJSONObject(br)
	if err != nil {
		return nil, err
	}

	if isObject {
		var single jsonrpcMessage
		if err := json.NewDecoder(br).Decode(&single); err != nil {
			return nil, err
		}
		if single.Error != nil {
			return []*jsonrpcMessage{&single}, nil
		}
		return nil, errors.New("batch response is a single object without an error")
	}

	var respmsgs []*jsonrpcMessage
	if err := json.NewDecoder(br).Decode(&respmsgs); err != nil {
		return nil, err
	}
	return respmsgs, nil
}

// peekIsJSONObject reports whether the next JSON value on br is an object
// ('{') as opposed to an array ('['), without consuming any bytes beyond
// leading whitespace.
func peekIsJSONObject(br *bufio.Reader) (bool, error) {
	for {
		b, err := br.Peek(1)
		if err != nil {
			return false, err
		}
		switch b[0] {
		case ' ', '\t', '\n', '\r':
			br.Discard(1)
		default:
			return b[0] == '{', nil
		}
	}
}

// cleanlyCloseBody avoids sending unnecessary RST_STREAM and PING frames by
// ensuring the whole body is read before being closed.
// See https://blog.cloudflare.com/go-and-enhance-your-calm/#reading-bodies-in-go-can-be-unintuitive
func cleanlyCloseBody(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, body)
	_ = body.Close()
}

// peekReader wraps an io.Reader, remembering everything read through it so
// callers can fall back to re-parsing the full response after a decode
// error, without doing an eager io.ReadAll on the happy path. buffered()
// truncates the preview shown in error messages to limit bytes.
type peekReader struct {
	r     io.Reader
	buf   bytes.Buffer
	limit int
}

func newPeekReader(r io.Reader, limit int) *peekReader {
	return &peekReader{r: r, limit: limit}
}

func (p *peekReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	if n > 0 {
		p.buf.Write(b[:n])
	}
	return n, err
}

// bytesRead returns everything read through this reader so far.
func (p *peekReader) bytesRead() []byte {
	return p.buf.Bytes()
}

// buffered returns a preview of what was read, truncated to limit bytes.
func (p *peekReader) buffered() string {
	b := p.buf.Bytes()
	if len(b) > p.limit {
		b = b[:p.limit]
	}
	return string(b)
}
