package ethclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestAsHTTPError_RPCHTTPError(t *testing.T) {
	raw := rpc.HTTPError{
		StatusCode: 429,
		Status:     "429 Too Many Requests",
		Body:       []byte("rate limited"),
	}
	wrapped := fmt.Errorf("call failed: %w", raw)

	he, ok := AsHTTPError(wrapped)
	if !ok {
		t.Fatal("expected AsHTTPError to return true")
	}
	if he.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", he.StatusCode)
	}
	if he.Header != nil {
		t.Fatal("expected nil header for plain rpc.HTTPError")
	}
}

func TestAsHTTPError_WrappedHTTPError(t *testing.T) {
	he := &HTTPError{
		StatusCode: 503,
		Status:     "503 Service Unavailable",
		Body:       []byte("try again"),
		Header:     http.Header{"Retry-After": {"30"}},
	}
	wrapped := fmt.Errorf("outer: %w", he)

	got, ok := AsHTTPError(wrapped)
	if !ok {
		t.Fatal("expected AsHTTPError to return true")
	}
	if got.StatusCode != 503 {
		t.Fatalf("expected 503, got %d", got.StatusCode)
	}
	if got.Header.Get("Retry-After") != "30" {
		t.Fatalf("expected Retry-After=30, got %q", got.Header.Get("Retry-After"))
	}
}

func TestAsHTTPError_NonHTTPError(t *testing.T) {
	_, ok := AsHTTPError(errors.New("some random error"))
	if ok {
		t.Fatal("expected AsHTTPError to return false for non-HTTP error")
	}
}

func TestAsHTTPError_Nil(t *testing.T) {
	_, ok := AsHTTPError(nil)
	if ok {
		t.Fatal("expected AsHTTPError to return false for nil")
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "429 rpc.HTTPError",
			err:  rpc.HTTPError{StatusCode: 429, Status: "429 Too Many Requests"},
			want: true,
		},
		{
			name: "429 HTTPError",
			err:  &HTTPError{StatusCode: 429, Status: "429 Too Many Requests"},
			want: true,
		},
		{
			name: "500 error",
			err:  rpc.HTTPError{StatusCode: 500, Status: "500 Internal Server Error"},
			want: false,
		},
		{
			name: "non-HTTP error",
			err:  errors.New("connection refused"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimited(tt.err); got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code int
		want bool
	}{
		{
			name: "match 502",
			err:  rpc.HTTPError{StatusCode: 502, Status: "502 Bad Gateway"},
			code: 502,
			want: true,
		},
		{
			name: "no match",
			err:  rpc.HTTPError{StatusCode: 503, Status: "503 Service Unavailable"},
			code: 502,
			want: false,
		},
		{
			name: "wrapped match",
			err:  fmt.Errorf("rpc: %w", rpc.HTTPError{StatusCode: 403, Status: "403 Forbidden"}),
			code: 403,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHTTPStatus(tt.err, tt.code); got != tt.want {
				t.Errorf("IsHTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryAfter(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want time.Duration
	}{
		{
			name: "with Retry-After header",
			err: &HTTPError{
				StatusCode: 429,
				Status:     "429 Too Many Requests",
				Header:     http.Header{"Retry-After": {"5"}},
			},
			want: 5 * time.Second,
		},
		{
			name: "no header (plain rpc.HTTPError)",
			err:  rpc.HTTPError{StatusCode: 429, Status: "429 Too Many Requests"},
			want: 0,
		},
		{
			name: "non-429 with Retry-After",
			err: &HTTPError{
				StatusCode: 503,
				Status:     "503 Service Unavailable",
				Header:     http.Header{"Retry-After": {"10"}},
			},
			want: 0,
		},
		{
			name: "invalid Retry-After value",
			err: &HTTPError{
				StatusCode: 429,
				Status:     "429 Too Many Requests",
				Header:     http.Header{"Retry-After": {"not-a-number"}},
			},
			want: 0,
		},
		{
			name: "zero Retry-After",
			err: &HTTPError{
				StatusCode: 429,
				Status:     "429 Too Many Requests",
				Header:     http.Header{"Retry-After": {"0"}},
			},
			want: 0,
		},
		{
			name: "nil error",
			err:  nil,
			want: 0,
		},
		{
			name: "past HTTP-date",
			err: &HTTPError{
				StatusCode: 429,
				Status:     "429 Too Many Requests",
				Header:     http.Header{"Retry-After": {"Mon, 01 Jan 2001 00:00:00 GMT"}},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetryAfter(tt.err); got != tt.want {
				t.Errorf("RetryAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryAfter_HTTPDate(t *testing.T) {
	future := time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)
	err := &HTTPError{
		StatusCode: 429,
		Status:     "429 Too Many Requests",
		Header:     http.Header{"Retry-After": {future}},
	}
	d := RetryAfter(err)
	if d < 25*time.Second || d > 31*time.Second {
		t.Fatalf("expected ~30s, got %v", d)
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *HTTPError
		want string
	}{
		{
			name: "with body",
			err:  &HTTPError{StatusCode: 429, Status: "429 Too Many Requests", Body: []byte("slow down")},
			want: "429 Too Many Requests: slow down",
		},
		{
			name: "without body",
			err:  &HTTPError{StatusCode: 429, Status: "429 Too Many Requests"},
			want: "429 Too Many Requests",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTTPError_Unwrap(t *testing.T) {
	he := &HTTPError{StatusCode: 429, Status: "429 Too Many Requests", Body: []byte("x")}
	var rpcErr rpc.HTTPError
	if !errors.As(he, &rpcErr) {
		t.Fatal("expected Unwrap to produce rpc.HTTPError")
	}
	if rpcErr.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", rpcErr.StatusCode)
	}
}

func TestWrapErr(t *testing.T) {
	t.Run("nil rt", func(t *testing.T) {
		origErr := fmt.Errorf("call: %w", rpc.HTTPError{StatusCode: 429})
		got := wrapErr(nil, origErr)
		if got.Error() != origErr.Error() {
			t.Fatal("expected no-op for nil rt")
		}
	})

	t.Run("nil err", func(t *testing.T) {
		rt := &headerCapture{base: http.DefaultTransport}
		if wrapErr(rt, nil) != nil {
			t.Fatal("expected nil for nil err")
		}
	})

	t.Run("non-HTTP error", func(t *testing.T) {
		rt := &headerCapture{base: http.DefaultTransport}
		err := errors.New("something else")
		if wrapErr(rt, err) != err {
			t.Fatal("expected original error returned")
		}
	})

	t.Run("wraps rpc.HTTPError with captured header", func(t *testing.T) {
		rt := &headerCapture{base: http.DefaultTransport}
		rt.mu.Lock()
		rt.last = http.Header{"Retry-After": {"10"}}
		rt.mu.Unlock()

		err := fmt.Errorf("call: %w", rpc.HTTPError{StatusCode: 429, Status: "429 Too Many Requests"})
		wrapped := wrapErr(rt, err)

		he, ok := AsHTTPError(wrapped)
		if !ok {
			t.Fatal("expected HTTPError")
		}
		if he.StatusCode != 429 {
			t.Fatalf("expected 429, got %d", he.StatusCode)
		}
		if he.Header.Get("Retry-After") != "10" {
			t.Fatalf("expected Retry-After=10, got %q", he.Header.Get("Retry-After"))
		}

		// header should be consumed
		if rt.last != nil {
			t.Fatal("expected header to be consumed")
		}
	})
}

func TestHeaderCapture_RoundTrip(t *testing.T) {
	// Use a fake RoundTripper
	hc := &headerCapture{
		base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 429,
				Header:     http.Header{"Retry-After": {"15"}},
				Body:       http.NoBody,
			}, nil
		}),
	}

	req, _ := http.NewRequest("POST", "http://localhost", nil)
	resp, err := hc.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", resp.StatusCode)
	}

	h := hc.consumeHeader()
	if h.Get("Retry-After") != "15" {
		t.Fatalf("expected Retry-After=15, got %q", h.Get("Retry-After"))
	}

	// Second consume should return nil
	if hc.consumeHeader() != nil {
		t.Fatal("expected nil on second consume")
	}
}

func TestHeaderCapture_2xxNotCaptured(t *testing.T) {
	hc := &headerCapture{
		base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"X-Custom": {"val"}},
				Body:       http.NoBody,
			}, nil
		}),
	}

	req, _ := http.NewRequest("POST", "http://localhost", nil)
	_, err := hc.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	if hc.consumeHeader() != nil {
		t.Fatal("expected no header captured for 2xx response")
	}
}

// roundTripFunc is a helper for testing http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRateLimitEndToEnd(t *testing.T) {
	// Start a fake RPC server that always returns 429 with Retry-After header.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("rate limited"))
	}))
	defer srv.Close()

	ctx := t.Context()
	ec, err := DialContext(ctx, srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer ec.Close()

	_, err = ec.ChainID(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsRateLimited(err) {
		t.Fatalf("expected rate limited, got: %v", err)
	}

	he, ok := AsHTTPError(err)
	if !ok {
		t.Fatal("expected AsHTTPError to return true")
	}
	if he.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", he.StatusCode)
	}
	if he.Header.Get("Retry-After") != "5" {
		t.Fatalf("expected Retry-After=5, got %q", he.Header.Get("Retry-After"))
	}
	if string(he.Body) != "rate limited" {
		t.Fatalf("expected body 'rate limited', got %q", string(he.Body))
	}

	d := RetryAfter(err)
	if d != 5*time.Second {
		t.Fatalf("expected RetryAfter=5s, got %v", d)
	}
}

func TestRateLimitHttpstatUs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external HTTP test")
	}

	urls := []string{
		"https://httpstat.us/429",
		"https://httpbin.org/status/429",
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var resp *http.Response
	for _, u := range urls {
		r, err := client.Post(u, "application/json", nil)
		if err == nil {
			resp = r
			t.Logf("using %s", u)
			break
		}
		t.Logf("%s unreachable: %v", u, err)
	}
	if resp == nil {
		t.Skip("no 429 endpoint reachable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 429 {
		t.Fatalf("expected 429, got %d", resp.StatusCode)
	}
	t.Logf("Status: %s", resp.Status)
	t.Logf("Retry-After: %q", resp.Header.Get("Retry-After"))
	t.Logf("Headers: %v", resp.Header)

	// Build an HTTPError from the response to test our parsing.
	he := &HTTPError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Header:     resp.Header,
	}

	if !IsRateLimited(he) {
		t.Fatal("expected IsRateLimited to return true")
	}

	// These endpoints may or may not include Retry-After; log what we get.
	d := RetryAfter(he)
	t.Logf("RetryAfter parsed: %v", d)
}

// https://www.alchemy.com/docs/reference/throughput
const alchemySepoliaRPC = "https://eth-sepolia.g.alchemy.com/v2/AxnmGEYn7VDkC4KqfNSFbSW9pHFR7PDO"

func TestRateLimitAlchemy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external RPC test")
	}

	ctx := t.Context()
	ec, err := DialContext(ctx, alchemySepoliaRPC)
	if err != nil {
		t.Skipf("cannot connect: %v", err)
	}
	defer ec.Close()

	// Verify the endpoint works.
	if _, err := ec.ChainID(ctx); err != nil {
		t.Skipf("endpoint unavailable: %v", err)
	}

	const goroutines = 200
	var rateLimited atomic.Int64
	var succeeded atomic.Int64
	var otherErr atomic.Int64

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			_, err := ec.BlockNumber(ctx)
			if err == nil {
				succeeded.Add(1)
				return
			}
			if IsRateLimited(err) {
				rateLimited.Add(1)
				he, ok := AsHTTPError(err)
				if !ok {
					t.Errorf("IsRateLimited=true but AsHTTPError=false")
				}
				t.Logf("429: status=%q retry-after=%q body=%q",
					he.Status, he.Header.Get("Retry-After"), string(he.Body))
				return
			}
			otherErr.Add(1)
			t.Logf("other error: %v", err)
		}()
	}
	wg.Wait()

	t.Logf("results: succeeded=%d rateLimited=%d otherErr=%d",
		succeeded.Load(), rateLimited.Load(), otherErr.Load())

	if rateLimited.Load() == 0 {
		t.Fatal("expected at least one 429, but got none — throughput limit may have changed")
	}
}

func TestOverloadDetector_Basic(t *testing.T) {
	od := newOverloadDetector(10, 0.5)

	// Empty window → not overloaded, rate = 0.
	if od.isOverloaded() {
		t.Fatal("empty window should not be overloaded")
	}
	if r := od.rate(); r != 0 {
		t.Fatalf("expected rate 0, got %f", r)
	}

	// Record 4 successes, 1 failure → rate = 0.2.
	for range 4 {
		od.record(false)
	}
	od.record(true)
	if od.isOverloaded() {
		t.Fatal("20% rate should not be overloaded at threshold 50%")
	}
	if r := od.rate(); r < 0.19 || r > 0.21 {
		t.Fatalf("expected rate ~0.2, got %f", r)
	}

	// Record 5 more 429s → 6/10 = 60% → overloaded.
	for range 5 {
		od.record(true)
	}
	if !od.isOverloaded() {
		t.Fatal("60% rate should be overloaded at threshold 50%")
	}
	if r := od.rate(); r < 0.59 || r > 0.61 {
		t.Fatalf("expected rate ~0.6, got %f", r)
	}
}

func TestOverloadDetector_RingWrap(t *testing.T) {
	od := newOverloadDetector(5, 0.5)

	// Fill window with 5 x 429.
	for range 5 {
		od.record(true)
	}
	if r := od.rate(); r != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", r)
	}

	// Overwrite all with successes.
	for range 5 {
		od.record(false)
	}
	if r := od.rate(); r != 0 {
		t.Fatalf("expected rate 0, got %f", r)
	}
	if od.isOverloaded() {
		t.Fatal("should not be overloaded after all successes")
	}
}

func TestOverloadDetector_ExactThreshold(t *testing.T) {
	od := newOverloadDetector(4, 0.5)

	// 2/4 = 0.5 → exactly at threshold → overloaded.
	od.record(true)
	od.record(true)
	od.record(false)
	od.record(false)

	if !od.isOverloaded() {
		t.Fatal("50% rate should be overloaded at threshold 50% (>=)")
	}
}

func TestOverloadDetector_Concurrent(t *testing.T) {
	od := newOverloadDetector(100, 0.5)

	var wg sync.WaitGroup
	wg.Add(200)
	for range 200 {
		go func() {
			defer wg.Done()
			od.record(true)
			od.isOverloaded()
			od.rate()
		}()
	}
	wg.Wait()

	// After 200 records into a 100-slot window, count should be 100 and all 429.
	if r := od.rate(); r != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", r)
	}
}

func TestOverloadEndToEnd(t *testing.T) {
	var reqCount atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := reqCount.Add(1)
		if n%2 == 0 {
			// Every other request returns 429.
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
	}))
	defer srv.Close()

	ctx := t.Context()
	ec, err := DialContext(ctx, srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer ec.Close()

	// Initially not overloaded.
	if ec.IsOverloaded() {
		t.Fatal("should not be overloaded before any requests")
	}

	// Fire enough requests to fill the window and reach ~50% rate.
	for range 60 {
		ec.ChainID(ctx)
	}

	rate := ec.OverloadRate()
	t.Logf("overload rate after 60 requests: %.2f", rate)
	if rate < 0.3 {
		t.Fatalf("expected rate >= 0.3, got %f", rate)
	}
}

func TestOverloadAlchemy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external RPC test")
	}

	ctx := t.Context()
	ec, err := DialContext(ctx, alchemySepoliaRPC)
	if err != nil {
		t.Skipf("cannot connect: %v", err)
	}
	defer ec.Close()

	if _, err := ec.ChainID(ctx); err != nil {
		t.Skipf("endpoint unavailable: %v", err)
	}

	// Hammer with concurrent calls to trigger 429s.
	const goroutines = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			ec.BlockNumber(ctx)
		}()
	}
	wg.Wait()

	rate := ec.OverloadRate()
	t.Logf("overload rate after %d concurrent calls: %.2f, isOverloaded=%v",
		goroutines, rate, ec.IsOverloaded())

	// We expect some 429s from Alchemy's free tier with 200 concurrent calls.
	if rate == 0 {
		t.Log("warning: no 429s detected — throughput limit may have changed")
	}
}

// TestOverloadAlchemyBlockReceipts uses eth_getBlockReceipts (500 CU per call on Alchemy)
// to reliably trigger overload, then verifies IsOverloaded works as a backpressure signal.
func TestOverloadAlchemyBlockReceipts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external RPC test")
	}

	ctx := t.Context()
	ec, err := DialContext(ctx, alchemySepoliaRPC)
	if err != nil {
		t.Skipf("cannot connect: %v", err)
	}
	defer ec.Close()

	// Verify the endpoint works and get a recent block number.
	blockNum, err := ec.BlockNumber(ctx)
	if err != nil {
		t.Skipf("endpoint unavailable: %v", err)
	}
	t.Logf("latest block: %d", blockNum)

	// Phase 1: Blast concurrent BlockReceipts to trigger overload.
	const goroutines = 100
	var rateLimited, succeeded, otherErr atomic.Int64

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func() {
			defer wg.Done()
			bn := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNum - uint64(i%10)))
			_, err := ec.BlockReceipts(ctx, bn)
			if err == nil {
				succeeded.Add(1)
				return
			}
			if IsRateLimited(err) {
				rateLimited.Add(1)
				return
			}
			otherErr.Add(1)
		}()
	}
	wg.Wait()

	overloaded := ec.IsOverloaded()
	rate := ec.OverloadRate()
	t.Logf("phase 1 (blast): succeeded=%d rateLimited=%d otherErr=%d rate=%.2f isOverloaded=%v",
		succeeded.Load(), rateLimited.Load(), otherErr.Load(), rate, overloaded)

	if rateLimited.Load() == 0 {
		t.Skip("no 429s triggered — Alchemy throughput limit may have changed, cannot test overload detection")
	}

	// Verify that IsOverloaded correctly detects the overload state.
	if !overloaded {
		t.Fatalf("expected IsOverloaded=true after %d 429s out of %d requests (rate=%.2f)",
			rateLimited.Load(), goroutines, rate)
	}

	// Phase 2: Simulate the backpressure loop — wait while overloaded, then resume.
	// After the blast, the 50-slot window is dominated by 429s. Each successful
	// call pushes out one 429. With Alchemy free-tier recovery we need patience.
	t.Log("phase 2: backpressure loop — waiting for overload to clear...")
	waitStart := time.Now()
	var backoffCalls int
	for ec.IsOverloaded() {
		backoffCalls++
		time.Sleep(500 * time.Millisecond)

		// Make a single cheap call to feed new samples into the window.
		_, callErr := ec.BlockNumber(ctx)
		if callErr != nil {
			t.Logf("  backoff iter %d: still rate limited (rate=%.2f)", backoffCalls, ec.OverloadRate())
		} else {
			t.Logf("  backoff iter %d: success (rate=%.2f)", backoffCalls, ec.OverloadRate())
		}

		if time.Since(waitStart) > 2*time.Minute {
			t.Fatalf("overload did not clear within 2min (rate=%.2f)", ec.OverloadRate())
		}
	}
	t.Logf("phase 2: overload cleared after %d backoff iterations (%v), rate=%.2f",
		backoffCalls, time.Since(waitStart).Round(time.Millisecond), ec.OverloadRate())

	// Phase 3: After overload clears, a normal call should succeed.
	_, err = ec.BlockNumber(ctx)
	if err != nil {
		t.Logf("phase 3: first call after backoff still failed: %v", err)
	} else {
		t.Log("phase 3: call succeeded after overload cleared")
	}
	t.Logf("final overload rate: %.2f, isOverloaded=%v", ec.OverloadRate(), ec.IsOverloaded())
}
