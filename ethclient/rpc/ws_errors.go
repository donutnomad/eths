package rpc

// Error types defined below are the built-in JSON-RPC errors.

var _ Error = new(invalidRequestError)

const (
	errcodeDefault                        = -32000
	legacyErrcodeNotificationsUnsupported = -32001
)

type notificationsUnsupportedError struct{}

func (e notificationsUnsupportedError) Error() string {
	return "notifications not supported"
}

func (e notificationsUnsupportedError) ErrorCode() int { return -32601 }

// Is checks for equivalence to another error. Here we define that all errors with code
// -32601 (method not found) are equivalent to notificationsUnsupportedError. This is
// done to enable the following pattern:
//
//	sub, err := client.Subscribe(...)
//	if errors.Is(err, rpc.ErrNotificationsUnsupported) {
//		// server doesn't support subscriptions
//	}
func (e notificationsUnsupportedError) Is(other error) bool {
	if other == (notificationsUnsupportedError{}) {
		return true
	}
	rpcErr, ok := other.(Error)
	if ok {
		code := rpcErr.ErrorCode()
		return code == -32601 || code == legacyErrcodeNotificationsUnsupported
	}
	return false
}

// received message isn't a valid request
type invalidRequestError struct{ message string }

func (e *invalidRequestError) ErrorCode() int { return -32600 }

func (e *invalidRequestError) Error() string { return e.message }
