// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts_pack

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// ERC20MetaData contains all meta data concerning the ERC20 contract.
var ERC20MetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	ID:  "ERC20",
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	abi abi.ABI
}

// NewERC20 creates a new instance of ERC20.
func NewERC20() *ERC20 {
	parsed, err := ERC20MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC20{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC20) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (eRC20 *ERC20) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := eRC20.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (eRC20 *ERC20) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return eRC20.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (eRC20 *ERC20) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (eRC20 *ERC20) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := eRC20.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (eRC20 *ERC20) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackApprove(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (eRC20 *ERC20) PackBalanceOf(account common.Address) []byte {
	enc, err := eRC20.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (eRC20 *ERC20) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return eRC20.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (eRC20 *ERC20) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) PackDecimals() []byte {
	enc, err := eRC20.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) TryPackDecimals() ([]byte, error) {
	return eRC20.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) UnpackDecimals(data []byte) (uint8, error) {
	out, err := eRC20.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) PackName() []byte {
	enc, err := eRC20.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) TryPackName() ([]byte, error) {
	return eRC20.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) UnpackName(data []byte) (string, error) {
	out, err := eRC20.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) PackSymbol() []byte {
	enc, err := eRC20.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) TryPackSymbol() ([]byte, error) {
	return eRC20.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) UnpackSymbol(data []byte) (string, error) {
	out, err := eRC20.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) PackTotalSupply() []byte {
	enc, err := eRC20.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) TryPackTotalSupply() ([]byte, error) {
	return eRC20.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := eRC20.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (eRC20 *ERC20) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := eRC20.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (eRC20 *ERC20) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackTransfer(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (eRC20 *ERC20) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := eRC20.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (eRC20 *ERC20) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return eRC20.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := eRC20.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ERC20ApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (ERC20Approval) ContractEventName() string {
	return ERC20ApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (eRC20 *ERC20) UnpackApprovalEvent(log *types.Log) (*ERC20Approval, error) {
	event := "Approval"
	if len(log.Topics) == 0 || log.Topics[0] != eRC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC20Approval)
	if len(log.Data) > 0 {
		if err := eRC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC20.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const ERC20TransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (ERC20Transfer) ContractEventName() string {
	return ERC20TransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (eRC20 *ERC20) UnpackTransferEvent(log *types.Log) (*ERC20Transfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 || log.Topics[0] != eRC20.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC20Transfer)
	if len(log.Data) > 0 {
		if err := eRC20.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC20.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (eRC20 *ERC20) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC20.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return eRC20.UnpackERC20InvalidSpenderError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ERC20ERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the ERC20 contract.
type ERC20ERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func ERC20ERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (eRC20 *ERC20) UnpackERC20InsufficientAllowanceError(raw []byte) (*ERC20ERC20InsufficientAllowance, error) {
	out := new(ERC20ERC20InsufficientAllowance)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC20ERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the ERC20 contract.
type ERC20ERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func ERC20ERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (eRC20 *ERC20) UnpackERC20InsufficientBalanceError(raw []byte) (*ERC20ERC20InsufficientBalance, error) {
	out := new(ERC20ERC20InsufficientBalance)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC20ERC20InvalidApprover represents a ERC20InvalidApprover error raised by the ERC20 contract.
type ERC20ERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func ERC20ERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (eRC20 *ERC20) UnpackERC20InvalidApproverError(raw []byte) (*ERC20ERC20InvalidApprover, error) {
	out := new(ERC20ERC20InvalidApprover)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC20ERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the ERC20 contract.
type ERC20ERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func ERC20ERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (eRC20 *ERC20) UnpackERC20InvalidReceiverError(raw []byte) (*ERC20ERC20InvalidReceiver, error) {
	out := new(ERC20ERC20InvalidReceiver)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC20ERC20InvalidSender represents a ERC20InvalidSender error raised by the ERC20 contract.
type ERC20ERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func ERC20ERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (eRC20 *ERC20) UnpackERC20InvalidSenderError(raw []byte) (*ERC20ERC20InvalidSender, error) {
	out := new(ERC20ERC20InvalidSender)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC20ERC20InvalidSpender represents a ERC20InvalidSpender error raised by the ERC20 contract.
type ERC20ERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func ERC20ERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (eRC20 *ERC20) UnpackERC20InvalidSpenderError(raw []byte) (*ERC20ERC20InvalidSpender, error) {
	out := new(ERC20ERC20InvalidSpender)
	if err := eRC20.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UnpackEvent unpacks event log based on topic0.
func (eRC20 *ERC20) UnpackEvent(log *types.Log) (interface {
	ContractEventName() string
	Topic0() common.Hash
}, error) {
	var mismatch = errors.New("event signature mismatch")
	if len(log.Topics) == 0 {
		return nil, mismatch
	}
	topic0 := log.Topics[0]
	if topic0 == eRC20.abi.Events["Approval"].ID {
		return eRC20.UnpackApprovalEvent(log)
	}
	if topic0 == eRC20.abi.Events["Transfer"].ID {
		return eRC20.UnpackTransferEvent(log)
	}
	return nil, mismatch
}

// ERC20ApprovalTopic0 returns the hash of the event signature.
//
// Solidity: event Approval(address owner, address spender, uint256 value)
func ERC20ApprovalTopic0() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

// ERC20TransferTopic0 returns the hash of the event signature.
//
// Solidity: event Transfer(address from, address to, uint256 value)
func ERC20TransferTopic0() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

// AllowanceMethodID returns the method ID for allowance (0xdd62ed3e).
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (eRC20 *ERC20) AllowanceMethodID() [4]byte {
	return [4]byte{0xdd, 0x62, 0xed, 0x3e}
}

// ApproveMethodID returns the method ID for approve (0x095ea7b3).
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (eRC20 *ERC20) ApproveMethodID() [4]byte {
	return [4]byte{0x09, 0x5e, 0xa7, 0xb3}
}

// BalanceOfMethodID returns the method ID for balanceOf (0x70a08231).
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (eRC20 *ERC20) BalanceOfMethodID() [4]byte {
	return [4]byte{0x70, 0xa0, 0x82, 0x31}
}

// DecimalsMethodID returns the method ID for decimals (0x313ce567).
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) DecimalsMethodID() [4]byte {
	return [4]byte{0x31, 0x3c, 0xe5, 0x67}
}

// NameMethodID returns the method ID for name (0x06fdde03).
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) NameMethodID() [4]byte {
	return [4]byte{0x06, 0xfd, 0xde, 0x03}
}

// SymbolMethodID returns the method ID for symbol (0x95d89b41).
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) SymbolMethodID() [4]byte {
	return [4]byte{0x95, 0xd8, 0x9b, 0x41}
}

// TotalSupplyMethodID returns the method ID for totalSupply (0x18160ddd).
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) TotalSupplyMethodID() [4]byte {
	return [4]byte{0x18, 0x16, 0x0d, 0xdd}
}

// TransferMethodID returns the method ID for transfer (0xa9059cbb).
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (eRC20 *ERC20) TransferMethodID() [4]byte {
	return [4]byte{0xa9, 0x05, 0x9c, 0xbb}
}

// TransferFromMethodID returns the method ID for transferFrom (0x23b872dd).
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (eRC20 *ERC20) TransferFromMethodID() [4]byte {
	return [4]byte{0x23, 0xb8, 0x72, 0xdd}
}

// Topic0 returns the hash of the event signature.
//
// Solidity: event Approval(address owner, address spender, uint256 value)
func (ERC20Approval) Topic0() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

// Topic0 returns the hash of the event signature.
//
// Solidity: event Transfer(address from, address to, uint256 value)
func (ERC20Transfer) Topic0() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

// Signature returns the event signature string.
//
// Solidity: event Approval(address owner, address spender, uint256 value)
func (ERC20Approval) Signature() string {
	return "Approval(address,address,uint256)"
}

// Signature returns the event signature string.
//
// Solidity: event Transfer(address from, address to, uint256 value)
func (ERC20Transfer) Signature() string {
	return "Transfer(address,address,uint256)"
}

// UnpackInputAllowance unpacks the input data for the allowance method.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (eRC20 *ERC20) UnpackInputAllowance(callData []byte) (owner common.Address, spender common.Address, err error) {
	method, ok := eRC20.abi.Methods["allowance"]
	if !ok {
		return common.Address{}, common.Address{}, errors.New("method 'allowance' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return common.Address{}, common.Address{}, errors.New("method signature mismatch")
	}
	arguments, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return common.Address{}, common.Address{}, err
	}
	owner = *abi.ConvertType(arguments[0], new(common.Address)).(*common.Address)
	spender = *abi.ConvertType(arguments[1], new(common.Address)).(*common.Address)
	return owner, spender, nil
}

// UnpackInputApprove unpacks the input data for the approve method.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackInputApprove(callData []byte) (spender common.Address, value *big.Int, err error) {
	method, ok := eRC20.abi.Methods["approve"]
	if !ok {
		return common.Address{}, nil, errors.New("method 'approve' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return common.Address{}, nil, errors.New("method signature mismatch")
	}
	arguments, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return common.Address{}, nil, err
	}
	spender = *abi.ConvertType(arguments[0], new(common.Address)).(*common.Address)
	value = abi.ConvertType(arguments[1], new(big.Int)).(*big.Int)
	return spender, value, nil
}

// UnpackInputBalanceOf unpacks the input data for the balanceOf method.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (eRC20 *ERC20) UnpackInputBalanceOf(callData []byte) (account common.Address, err error) {
	method, ok := eRC20.abi.Methods["balanceOf"]
	if !ok {
		return common.Address{}, errors.New("method 'balanceOf' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return common.Address{}, errors.New("method signature mismatch")
	}
	arguments, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return common.Address{}, err
	}
	account = *abi.ConvertType(arguments[0], new(common.Address)).(*common.Address)
	return account, nil
}

// UnpackInputDecimals unpacks the input data for the decimals method.
//
// Solidity: function decimals() view returns(uint8)
func (eRC20 *ERC20) UnpackInputDecimals(callData []byte) error {
	method, ok := eRC20.abi.Methods["decimals"]
	if !ok {
		return errors.New("method 'decimals' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return errors.New("method signature mismatch")
	}
	return nil
}

// UnpackInputName unpacks the input data for the name method.
//
// Solidity: function name() view returns(string)
func (eRC20 *ERC20) UnpackInputName(callData []byte) error {
	method, ok := eRC20.abi.Methods["name"]
	if !ok {
		return errors.New("method 'name' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return errors.New("method signature mismatch")
	}
	return nil
}

// UnpackInputSymbol unpacks the input data for the symbol method.
//
// Solidity: function symbol() view returns(string)
func (eRC20 *ERC20) UnpackInputSymbol(callData []byte) error {
	method, ok := eRC20.abi.Methods["symbol"]
	if !ok {
		return errors.New("method 'symbol' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return errors.New("method signature mismatch")
	}
	return nil
}

// UnpackInputTotalSupply unpacks the input data for the totalSupply method.
//
// Solidity: function totalSupply() view returns(uint256)
func (eRC20 *ERC20) UnpackInputTotalSupply(callData []byte) error {
	method, ok := eRC20.abi.Methods["totalSupply"]
	if !ok {
		return errors.New("method 'totalSupply' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return errors.New("method signature mismatch")
	}
	return nil
}

// UnpackInputTransfer unpacks the input data for the transfer method.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackInputTransfer(callData []byte) (to common.Address, value *big.Int, err error) {
	method, ok := eRC20.abi.Methods["transfer"]
	if !ok {
		return common.Address{}, nil, errors.New("method 'transfer' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return common.Address{}, nil, errors.New("method signature mismatch")
	}
	arguments, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return common.Address{}, nil, err
	}
	to = *abi.ConvertType(arguments[0], new(common.Address)).(*common.Address)
	value = abi.ConvertType(arguments[1], new(big.Int)).(*big.Int)
	return to, value, nil
}

// UnpackInputTransferFrom unpacks the input data for the transferFrom method.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (eRC20 *ERC20) UnpackInputTransferFrom(callData []byte) (from common.Address, to common.Address, value *big.Int, err error) {
	method, ok := eRC20.abi.Methods["transferFrom"]
	if !ok {
		return common.Address{}, common.Address{}, nil, errors.New("method 'transferFrom' not found")
	}
	if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
		return common.Address{}, common.Address{}, nil, errors.New("method signature mismatch")
	}
	arguments, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return common.Address{}, common.Address{}, nil, err
	}
	from = *abi.ConvertType(arguments[0], new(common.Address)).(*common.Address)
	to = *abi.ConvertType(arguments[1], new(common.Address)).(*common.Address)
	value = abi.ConvertType(arguments[2], new(big.Int)).(*big.Int)
	return from, to, value, nil
}
