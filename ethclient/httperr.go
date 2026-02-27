package ethclient

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// HTTPError wraps rpc.HTTPError with additional HTTP response headers.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
	Header     http.Header
}

func (e *HTTPError) Error() string {
	if len(e.Body) == 0 {
		return e.Status
	}
	return e.Status + ": " + string(e.Body)
}

func (e *HTTPError) Unwrap() error {
	return rpc.HTTPError{
		StatusCode: e.StatusCode,
		Status:     e.Status,
		Body:       e.Body,
	}
}

// AsHTTPError extracts an HTTPError from the error chain.
// If the client was created via DialContext, the returned error includes
// captured response headers. Otherwise only StatusCode/Status/Body are available.
func AsHTTPError(err error) (*HTTPError, bool) {
	var he *HTTPError
	if errors.As(err, &he) {
		return he, true
	}
	var rpcErr rpc.HTTPError
	if errors.As(err, &rpcErr) {
		return &HTTPError{
			StatusCode: rpcErr.StatusCode,
			Status:     rpcErr.Status,
			Body:       rpcErr.Body,
		}, true
	}
	return nil, false
}

// IsRateLimited reports whether err is an HTTP 429 (Too Many Requests) error.
func IsRateLimited(err error) bool {
	return IsHTTPStatus(err, http.StatusTooManyRequests)
}

// IsHTTPStatus reports whether err is an HTTP error with the given status code.
func IsHTTPStatus(err error, code int) bool {
	he, ok := AsHTTPError(err)
	return ok && he.StatusCode == code
}

// RetryAfter extracts the Retry-After duration from a 429 error.
// It returns 0 if the error is not a 429 or has no Retry-After header.
// Both "delay-seconds" (e.g. "120") and "HTTP-date" (e.g. "Wed, 21 Oct 2015 07:28:00 GMT")
// formats are supported per RFC 9110.
func RetryAfter(err error) time.Duration {
	he, ok := AsHTTPError(err)
	if !ok || he.StatusCode != http.StatusTooManyRequests || he.Header == nil {
		return 0
	}
	v := he.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	// Try delay-seconds first.
	if secs, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
		if secs <= 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	// Try HTTP-date format.
	if t, parseErr := http.ParseTime(v); parseErr == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// headerCapture is an http.RoundTripper that captures response headers
// from non-2xx responses. The rpc package converts non-2xx responses to
// rpc.HTTPError, but that type does not include headers.
type headerCapture struct {
	base http.RoundTripper
	mu   sync.Mutex
	last http.Header // headers from the most recent non-2xx response
}

func (hc *headerCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := hc.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		hc.mu.Lock()
		hc.last = resp.Header.Clone()
		hc.mu.Unlock()
	}
	return resp, nil
}

// consumeHeader returns and clears the last captured header.
func (hc *headerCapture) consumeHeader() http.Header {
	hc.mu.Lock()
	h := hc.last
	hc.last = nil
	hc.mu.Unlock()
	return h
}

// overloadDetector uses a sliding window to track the 429 error rate.
type overloadDetector struct {
	mu        sync.Mutex
	window    []bool  // ring buffer, true = 429
	pos       int     // current write position
	count     int     // number of filled slots
	limited   int     // number of 429s in the window
	threshold float64 // rate >= threshold means overloaded
}

func newOverloadDetector(size int, threshold float64) *overloadDetector {
	return &overloadDetector{
		window:    make([]bool, size),
		threshold: threshold,
	}
}

// record writes a sample into the ring buffer.
func (od *overloadDetector) record(is429 bool) {
	od.mu.Lock()
	defer od.mu.Unlock()

	// If the buffer is full, subtract the value being overwritten.
	if od.count == len(od.window) {
		if od.window[od.pos] {
			od.limited--
		}
	} else {
		od.count++
	}
	od.window[od.pos] = is429
	if is429 {
		od.limited++
	}
	od.pos = (od.pos + 1) % len(od.window)
}

// isOverloaded reports whether the 429 rate meets or exceeds the threshold.
func (od *overloadDetector) isOverloaded() bool {
	od.mu.Lock()
	defer od.mu.Unlock()
	if od.count == 0 {
		return false
	}
	return float64(od.limited)/float64(od.count) >= od.threshold
}

// rate returns the current 429 rate (0.0 ~ 1.0).
func (od *overloadDetector) rate() float64 {
	od.mu.Lock()
	defer od.mu.Unlock()
	if od.count == 0 {
		return 0
	}
	return float64(od.limited) / float64(od.count)
}

// wrapErr enriches an rpc.HTTPError with captured response headers.
// If rt is nil or err does not contain rpc.HTTPError, err is returned as-is.
func wrapErr(rt *headerCapture, err error) error {
	if err == nil || rt == nil {
		return err
	}
	var rpcErr rpc.HTTPError
	if !errors.As(err, &rpcErr) {
		return err
	}
	return &HTTPError{
		StatusCode: rpcErr.StatusCode,
		Status:     rpcErr.Status,
		Body:       rpcErr.Body,
		Header:     rt.consumeHeader(),
	}
}
