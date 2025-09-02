package contractcall

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/samber/lo"
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

type InsufficientBalanceError struct {
	Balance *big.Int
	Err     *EvmError
}

func (e *InsufficientBalanceError) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return e.Err
}

func (e *InsufficientBalanceError) Error() string {
	if e.Balance == nil {
		return fmt.Sprintf("[insufficient eth balance (%s)]", e.Err)
	}
	return fmt.Sprintf("[insufficient eth balance (%s)]", e.Balance)
}

type EstimateGasError struct {
	Err error
}

func (e *EstimateGasError) Error() string {
	return e.Err.Error()
}
func (e *EstimateGasError) Unwrap() error {
	return e.Err
}

func IsEstimateGasError(err error) bool {
	var estimateGasError *EstimateGasError
	return errors.As(err, &estimateGasError)
}

type JsonError interface {
	Error() string
	ErrorCode() int
	ErrorData() any
}

type EvmError struct {
	Err     error
	Message string
	ErrCode int
	ErrData string
}

func (e *EvmError) Unwrap() error {
	return e.Err
}

// EvmError: -32000 execution reverted
// -32000 gas required exceeds allowance (gas insufficient)
// code=3, data=(string) (len=74) "0xab0b880c0000000000000000000000003ae83667b13f48f236278839b18155e078059542" contract execute failed
func (e *EvmError) Error() string {
	var args []string
	args = append(args, fmt.Sprintf("code=%d", e.ErrCode))
	args = append(args, fmt.Sprintf("message=%q", e.Message))
	if e.IsExecutionReverted() {
		args = append(args, fmt.Sprintf("data=%q", e.ErrData))
	}
	if e.Err != nil {
		args = append(args, e.Err.Error())
	}
	return strings.Join(args, ",")
}

///////////////////////////////// ETH ERROR /////////////////////////////////

// GETH erros: https://github.com/ethereum/go-ethereum/blob/fe24d22a622926c024030480773faf83b78b1319/core/error.go
// errors: https://github.com/erigontech/erigon/blob/release/2.60/rpc/errors.go

// IsOutOfGas
// reth:
// 578   │     /// Gas limit was exceeded during execution.
// 579   │     /// Contains the gas limit.
// 580   │     #[error("out of gas: gas required exceeds: {0}")]
// 581   │     BasicOutOfGas(u64),
// 582   │     /// Gas limit was exceeded during memory expansion.
// 583   │     /// Contains the gas limit.
// 584   │     #[error("out of gas: gas exhausted during memory expansion: {0}")]
// 585   │     MemoryOutOfGas(u64),
// 586   │     /// Gas limit was exceeded during precompile execution.
// 587   │     /// Contains the gas limit.
// 588   │     #[error("out of gas: gas exhausted during precompiled contract execution: {0}")]
// 589   │     PrecompileOutOfGas(u64),
// 590   │     /// An operand to an opcode was invalid or out of range.
// 591   │     /// Contains the gas limit.
// 592   │     #[error("out of gas: invalid operand to an opcode: {0}")]
// 593   │     InvalidOperandOutOfGas(u64),
func (e *EvmError) IsOutOfGas() bool {
	msg := e.Message
	return strings.Contains(msg, "out of gas")
}

// IsGasTooLow low gasLimit
// geth: intrinsic gas too low: gas 2, minimum needed 21000
// reth: intrinsic gas too low
func (e *EvmError) IsGasTooLow() bool {
	msg := e.Message
	return strings.Contains(msg, "intrinsic gas too low")
}

// IsGasTooHigh
// geth: exceeds block gas limit
// reth: exceeds block gas limit
func (e *EvmError) IsGasTooHigh() bool {
	msg := e.Message
	return strings.Contains(msg, "exceeds block gas limit") ||
		strings.Contains(msg, "exceeds max transaction gas limit") || // reth: crates/rpc/rpc-eth-types/src/error/mod.rs
		// reth:
		/// Thrown when a new transaction is added to the pool, but then immediately discarded to
		/// respect the tx fee exceeds the configured cap
		// #[error("tx fee ({max_tx_fee_wei} wei) exceeds the configured cap ({tx_fee_cap_wei} wei)")]
		(strings.HasPrefix(msg, "tx fee") && strings.Contains(msg, "exceeds the configured cap")) || // reth
		strings.Contains(msg, "gas limit too high") || // reth
		strings.Contains(msg, "intrinsic gas too high") // reth
}

func (e *EvmError) IsNonceTooLow() bool {
	msg := e.Message
	return strings.Contains(msg, "nonce too low")
}

func (e *EvmError) IsNonceTooHigh() bool {
	msg := e.Message
	return strings.Contains(msg, "nonce too high")
}

func (e *EvmError) IsInsufficientBalance() bool {
	msg := e.Message
	// geth and reth
	if strings.Contains(msg, "insufficient funds for transfer") || strings.Contains(msg, "insufficient funds for gas") {
		return true
	}
	// reth
	if strings.Contains(msg, "gas required exceeds allowance") {
		return true
	}
	return false
}

// IsExecutionReverted Execution reverted
// geth: ✅https://github.com/ethereum/go-ethereum/blob/master/internal/ethapi/errors.go#L55
// reth: ✅
func (e *EvmError) IsExecutionReverted() bool {
	return e.ErrCode == 3
}

func (e *EvmError) RevertedReason() []byte {
	if !e.IsExecutionReverted() {
		return nil
	}
	bs, err := hex.DecodeString(strings.TrimPrefix(e.ErrData, "0x"))
	if err != nil {
		return nil
	}
	return bs
}

///////////////////////////////// ETH ERROR /////////////////////////////////

func ParseEvmError(err error) error {
	if err == nil {
		return nil
	}
	var ret = EvmError{
		Err: err,
	}
	var unwrapErr = err
	for {
		var v JsonError
		if errors.As(unwrapErr, &v) {
			ret.Message = v.Error()
			ret.ErrCode = v.ErrorCode()
			switch v := v.ErrorData().(type) {
			case string:
				ret.ErrData = v
			}
		}
		unwrapErr = unwrapOnce(unwrapErr)
		if lo.IsNil(unwrapErr) {
			break
		}
	}
	if ret.IsInsufficientBalance() {
		return &InsufficientBalanceError{
			Err: &ret,
		}
	}
	return &ret
}

func unwrapOnce(err error) (cause error) {
	switch e := err.(type) {
	case interface{ Cause() error }:
		return e.Cause()
	case interface{ Unwrap() error }:
		return e.Unwrap()
	}
	return nil
}
