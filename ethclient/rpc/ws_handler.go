// Copyright 2019 The go-ethereum Authors
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
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/donutnomad/eths/internal/log"
)

// handler processes JSON-RPC messages received over the connection. There is one
// handler per connection. Note that handler is not safe for concurrent use. Message
// handling never blocks indefinitely because calls are processed on background
// goroutines launched by handler.
//
// This client only ever sends requests/notifications and receives responses and
// subscription notifications in return - it never answers method calls made by the
// remote end, so there is no method dispatch table here.
//
// The entry points for incoming messages are:
//
//	h.handleMsg(message)
//	h.handleBatch(message)
//
// Outgoing calls use the requestOp struct. Register the request before sending it
// on the connection:
//
//	op := &requestOp{ids: ...}
//	h.addRequestOp(op)
//
// Now send the request, then wait for the reply to be delivered through handleMsg:
//
//	if err := op.wait(...); err != nil {
//		h.removeRequestOp(op) // timeout, etc.
//	}
type handler struct {
	respWait   map[string]*requestOp          // active client requests
	clientSubs map[string]*ClientSubscription // active client subscriptions
	callWG     sync.WaitGroup                 // pending call goroutines
	rootCtx    context.Context                // canceled by close()
	cancelRoot func()                         // cancel function for rootCtx
	conn       jsonWriter                     // where responses will be sent
	log        log.Logger
}

type callProc struct {
	ctx context.Context
}

func newHandler(connCtx context.Context, conn jsonWriter) *handler {
	rootCtx, cancelRoot := context.WithCancel(connCtx)
	h := &handler{
		conn:       conn,
		respWait:   make(map[string]*requestOp),
		clientSubs: make(map[string]*ClientSubscription),
		rootCtx:    rootCtx,
		cancelRoot: cancelRoot,
		log:        log.Root(),
	}
	if conn.remoteAddr() != "" {
		h.log = h.log.New("conn", conn.remoteAddr())
	}
	return h
}

// handleBatch processes a batch of incoming messages: responses are matched to
// pending requests, subscription notifications are delivered, and any call messages
// (which this client never answers) are dropped.
func (h *handler) handleBatch(msgs []*jsonrpcMessage) {
	// Emit error response for empty batches:
	if len(msgs) == 0 {
		h.startCallProc(func(cp *callProc) {
			resp := errorMessage(&invalidRequestError{"empty batch"})
			h.conn.writeJSON(cp.ctx, resp, true)
		})
		return
	}

	// Handle non-call messages first.
	// Here we need to find the requestOp that sent the request batch.
	calls := make([]*jsonrpcMessage, 0, len(msgs))
	h.handleResponses(msgs, func(msg *jsonrpcMessage) {
		calls = append(calls, msg)
	})
	if len(calls) == 0 {
		return
	}

	// This client never exposes any methods of its own, so any call arriving
	// from the remote end is simply dropped; log it for visibility.
	for _, msg := range calls {
		h.log.Debug("Dropping incoming call", "method", msg.Method)
	}
}

// handleMsg handles a single non-batch message.
func (h *handler) handleMsg(msg *jsonrpcMessage) {
	msgs := []*jsonrpcMessage{msg}
	h.handleResponses(msgs, func(msg *jsonrpcMessage) {
		h.log.Debug("Dropping incoming call", "method", msg.Method)
	})
}

// close cancels all requests except for inflightReq and waits for
// call goroutines to shut down.
func (h *handler) close(err error, inflightReq *requestOp) {
	h.cancelAllRequests(err, inflightReq)
	h.callWG.Wait()
	h.cancelRoot()
}

// addRequestOp registers a request operation.
func (h *handler) addRequestOp(op *requestOp) {
	for _, id := range op.ids {
		h.respWait[string(id)] = op
	}
}

// removeRequestOp stops waiting for the given request IDs.
func (h *handler) removeRequestOp(op *requestOp) {
	for _, id := range op.ids {
		delete(h.respWait, string(id))
	}
}

// cancelAllRequests unblocks and removes pending requests and active subscriptions.
func (h *handler) cancelAllRequests(err error, inflightReq *requestOp) {
	didClose := make(map[*requestOp]bool)
	if inflightReq != nil {
		didClose[inflightReq] = true
	}

	for id, op := range h.respWait {
		// Remove the op so that later calls will not close op.resp again.
		delete(h.respWait, id)

		if !didClose[op] {
			op.err = err
			close(op.resp)
			didClose[op] = true
		}
	}
	for id, sub := range h.clientSubs {
		delete(h.clientSubs, id)
		sub.close(err)
	}
}

// startCallProc runs fn in a new goroutine and starts tracking it in the h.calls wait group.
func (h *handler) startCallProc(fn func(*callProc)) {
	h.callWG.Add(1)
	go func() {
		ctx, cancel := context.WithCancel(h.rootCtx)
		defer h.callWG.Done()
		defer cancel()
		fn(&callProc{ctx: ctx})
	}()
}

// handleResponses processes method call responses.
func (h *handler) handleResponses(batch []*jsonrpcMessage, handleCall func(*jsonrpcMessage)) {
	var resolvedops []*requestOp
	handleResp := func(msg *jsonrpcMessage) {
		op := h.respWait[string(msg.ID)]
		if op == nil {
			h.log.Debug("Unsolicited RPC response", "reqid", idForLog{msg.ID})
			return
		}
		resolvedops = append(resolvedops, op)
		delete(h.respWait, string(msg.ID))

		// For subscription responses, start the subscription if the server
		// indicates success. EthSubscribe gets unblocked in either case through
		// the op.resp channel.
		if op.sub != nil {
			if msg.Error != nil {
				op.err = msg.Error
			} else {
				op.err = json.Unmarshal(msg.Result, &op.sub.subid)
				if op.err == nil {
					go op.sub.run()
					h.clientSubs[op.sub.subid] = op.sub
				}
			}
		}

		if !op.hadResponse {
			op.hadResponse = true
			op.resp <- batch
		}
	}

	for _, msg := range batch {
		start := time.Now()
		switch {
		case msg.isResponse():
			handleResp(msg)
			h.log.Trace("Handled RPC response", "reqid", idForLog{msg.ID}, "duration", time.Since(start))

		case msg.isNotification():
			if strings.HasSuffix(msg.Method, notificationMethodSuffix) {
				h.handleSubscriptionResult(msg)
				continue
			}
			handleCall(msg)

		default:
			handleCall(msg)
		}
	}

	for _, op := range resolvedops {
		h.removeRequestOp(op)
	}
}

// handleSubscriptionResult processes subscription notifications.
func (h *handler) handleSubscriptionResult(msg *jsonrpcMessage) {
	var result subscriptionResult
	if err := json.Unmarshal(msg.Params, &result); err != nil {
		h.log.Debug("Dropping invalid subscription message")
		return
	}
	if h.clientSubs[result.ID] != nil {
		h.clientSubs[result.ID].deliver(result.Result)
	}
}

type idForLog struct{ json.RawMessage }

func (id idForLog) String() string {
	if s, err := strconv.Unquote(string(id.RawMessage)); err == nil {
		return s
	}
	return string(id.RawMessage)
}
