package etherscan

import (
	"encoding/json"
	"testing"
)

// On a rate-limit error Etherscan returns the reason as a string in `result`;
// it must be captured in ErrorResult, while the `result:[]` empty-error case and
// the success case keep working.
func TestResponseCapturesErrorResult(t *testing.T) {
	const reason = "Max calls per sec rate limit reached (5/sec applies to this API key)"

	var rateLimited Response[[]json.RawMessage]
	if err := json.Unmarshal([]byte(`{"status":"0","message":"NOTOK","result":"`+reason+`"}`), &rateLimited); err != nil {
		t.Fatalf("unmarshal rate-limited: %v", err)
	}
	if rateLimited.ErrorResult != reason {
		t.Fatalf("ErrorResult = %q, want %q", rateLimited.ErrorResult, reason)
	}

	var noRecords Response[[]json.RawMessage]
	if err := json.Unmarshal([]byte(`{"status":"0","message":"No records found","result":[]}`), &noRecords); err != nil {
		t.Fatalf("unmarshal no-records (result is []) must not error: %v", err)
	}
	if noRecords.ErrorResult != "" {
		t.Fatalf("ErrorResult = %q, want empty for array result", noRecords.ErrorResult)
	}

	var ok Response[[]int]
	if err := json.Unmarshal([]byte(`{"status":"1","message":"OK","result":[1,2,3]}`), &ok); err != nil {
		t.Fatalf("unmarshal success: %v", err)
	}
	if len(ok.Result) != 3 || ok.ErrorResult != "" {
		t.Fatalf("Result = %v, ErrorResult = %q", ok.Result, ok.ErrorResult)
	}
}
