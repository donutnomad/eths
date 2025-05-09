package contractcall

import (
	"net/url"
	"reflect"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var EthereumRPCErr = errors.New("ethereum rpc error")
var GasInvalidGasPriceErr = errors.New("invalid gas price")
var TxBuilderMissingRequiredFieldErr = errors.New("missing required field")

// SendTransactionError Ethereum SendTransaction Error
type SendTransactionError struct {
	Tx  *ethTypes.Transaction
	Err error
}

func (e *SendTransactionError) Error() string {
	return e.Err.Error()
}
func (e *SendTransactionError) Unwrap() error {
	return e.Err
}

type JSONError struct {
	JSONUnsupportedTypeError *struct{ Type reflect.Type }
}

type HTTPError struct {
	// "net/http: nil Context"
	NilContext *struct{ Err error }
	// "net/http: invalid method %q", method
	InvalidMethod *struct{ Err error }
	// url.Parse()
	URLError *struct{ Err url.Error }
}

type EthereumCallError struct {
	// call result parameter must be pointer or nil interface: %v
	InvalidCallResult *struct{ Err error }
	JSONError         *JSONError
	// node: NewJWTAuth
	// "failed to create JWT token: %w"
	CreateJWTTokenFailed *struct{ Err error }
}
