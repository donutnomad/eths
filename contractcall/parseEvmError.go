package contractcall

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/lo"
)

type IErrData interface {
	ErrorData() interface{}
}
type IErrCode interface {
	ErrorCode() int
}

type EvmError struct {
	Err     error
	ErrCode int
	ErrData interface{}
}

// EvmError: -32000 execution reverted
// -32000 gas required exceeds allowance (gas insufficient)
// code=3, data=(string) (len=74) "0xab0b880c0000000000000000000000003ae83667b13f48f236278839b18155e078059542" contract execute failed
func (e *EvmError) Error() string {
	var args []string
	if e.ErrCode != 0 {
		args = append(args, fmt.Sprintf("code=%d", e.ErrCode))
	}
	if !lo.IsNil(e.ErrData) {
		args = append(args, fmt.Sprintf("data=%s", spew.Sdump(e.ErrData)))
	}
	if e.Err != nil {
		args = append(args, e.Err.Error())
	}
	return strings.Join(args, ", ")
}

// IsOutGas
// detail: gas required exceeds allowance
func (e *EvmError) IsOutGas() bool {
	msg := e.Error()
	return strings.Contains(msg, "gas required exceeds")
}

// IsInsufficientFunds
// only testnet
// code=ServerErr,msg=,err=code=-32003, insufficient funds for gas * price + value: have 0 want 4000000000000000000000
// detail: insufficient funds for gas * price + value: have 0 want 4000000000000000000000
func (e *EvmError) IsInsufficientFunds() bool {
	msg := e.Error()
	return strings.Contains(msg, "insufficient funds for gas")
}

// IsInsufficientFundsForTransfer
// only mainnet
// chain_id=1 error="code=ServerErr,msg=,err=code=-32000, insufficient funds for transfer\nnew tx"
// insufficient funds for transfer
func (e *EvmError) IsInsufficientFundsForTransfer() bool {
	msg := e.Error()
	return strings.Contains(msg, "insufficient funds for transfer")
}

func ParseEvmError(err error) error {
	if err == nil {
		return nil
	}
	var retErr EvmError
	retErr.Err = err
	if v, ok := err.(IErrData); ok && !lo.IsNil(v) {
		retErr.ErrData = v.ErrorData()
	}
	if v, ok := err.(IErrCode); ok {
		retErr.ErrCode = v.ErrorCode()
	}
	return &retErr
}
