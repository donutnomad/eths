//go:build etherscan_live

package etherscan

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/donutnomad/eths/common"
)

// Live: hits the real Etherscan API with enough concurrency to trip the per-second
// rate limit, then asserts the real reason surfaces through the full client path
// (resty -> Response.UnmarshalJSON -> client error), not just "NOTOK".
// Run: ETHERSCAN_API_KEY=xxx go test -tags etherscan_live ./etherscan/ -run TestLiveRateLimitReason -v
func TestLiveRateLimitReason(t *testing.T) {
	key := os.Getenv("ETHERSCAN_API_KEY")
	if key == "" {
		t.Skip("set ETHERSCAN_API_KEY")
	}
	c := NewEtherscanClient(key)
	addr := common.HexToAddress("0x0D45134a4847947304DB096A2d58083C2289FEC8")
	from := uint64(12151345)

	const N = 20
	errs := make([]error, N)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = c.GetLogs(context.Background(), 11155111, &addr, &from, nil, 1, 1000, GetLogsOptions{})
		}(i)
	}
	wg.Wait()

	var sawReason, sawErr bool
	for _, err := range errs {
		if err == nil {
			continue
		}
		sawErr = true
		t.Logf("err: %v", err)
		if strings.Contains(err.Error(), "Max calls per sec rate limit reached") {
			sawReason = true
		}
	}
	if !sawErr {
		t.Skip("no error triggered; increase concurrency or rerun")
	}
	if !sawReason {
		t.Fatal("got an error but the rate-limit reason was not surfaced; ErrorResult did not flow through")
	}
}
