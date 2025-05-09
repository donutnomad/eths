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

// ERC1967ProxyMetaData contains all meta data concerning the ERC1967Proxy contract.
var ERC1967ProxyMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	ID:  "ERC1967Proxy",
}

// ERC1967Proxy is an auto generated Go binding around an Ethereum contract.
type ERC1967Proxy struct {
	abi abi.ABI
}

// NewERC1967Proxy creates a new instance of ERC1967Proxy.
func NewERC1967Proxy() *ERC1967Proxy {
	parsed, err := ERC1967ProxyMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ERC1967Proxy{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ERC1967Proxy) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address implementation, bytes _data) payable returns()
func (eRC1967Proxy *ERC1967Proxy) PackConstructor(implementation common.Address, _data []byte) []byte {
	enc, err := eRC1967Proxy.abi.Pack("", implementation, _data)
	if err != nil {
		panic(err)
	}
	return enc
}

// ERC1967ProxyUpgraded represents a Upgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const ERC1967ProxyUpgradedEventName = "Upgraded"

// ContractEventName returns the user-defined event name.
func (ERC1967ProxyUpgraded) ContractEventName() string {
	return ERC1967ProxyUpgradedEventName
}

// UnpackUpgradedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Upgraded(address indexed implementation)
func (eRC1967Proxy *ERC1967Proxy) UnpackUpgradedEvent(log *types.Log) (*ERC1967ProxyUpgraded, error) {
	event := "Upgraded"
	if log.Topics[0] != eRC1967Proxy.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(ERC1967ProxyUpgraded)
	if len(log.Data) > 0 {
		if err := eRC1967Proxy.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eRC1967Proxy.abi.Events[event].Inputs {
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
func (eRC1967Proxy *ERC1967Proxy) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["AddressEmptyCode"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackAddressEmptyCodeError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackERC1967InvalidImplementationError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["ERC1967NonPayable"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackERC1967NonPayableError(raw[4:])
	}
	if bytes.Equal(raw[:4], eRC1967Proxy.abi.Errors["FailedCall"].ID.Bytes()[:4]) {
		return eRC1967Proxy.UnpackFailedCallError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ERC1967ProxyAddressEmptyCode represents a AddressEmptyCode error raised by the ERC1967Proxy contract.
type ERC1967ProxyAddressEmptyCode struct {
	Target common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressEmptyCode(address target)
func ERC1967ProxyAddressEmptyCodeErrorID() common.Hash {
	return common.HexToHash("0x9996b315c842ff135b8fc4a08ad5df1c344efbc03d2687aecc0678050d2aac89")
}

// UnpackAddressEmptyCodeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressEmptyCode(address target)
func (eRC1967Proxy *ERC1967Proxy) UnpackAddressEmptyCodeError(raw []byte) (*ERC1967ProxyAddressEmptyCode, error) {
	out := new(ERC1967ProxyAddressEmptyCode)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "AddressEmptyCode", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyERC1967InvalidImplementation represents a ERC1967InvalidImplementation error raised by the ERC1967Proxy contract.
type ERC1967ProxyERC1967InvalidImplementation struct {
	Implementation common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func ERC1967ProxyERC1967InvalidImplementationErrorID() common.Hash {
	return common.HexToHash("0x4c9c8ce3ceb3130f17f7cdba48d89b5b0129f266a8bac114e6e315a41879b617")
}

// UnpackERC1967InvalidImplementationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967InvalidImplementation(address implementation)
func (eRC1967Proxy *ERC1967Proxy) UnpackERC1967InvalidImplementationError(raw []byte) (*ERC1967ProxyERC1967InvalidImplementation, error) {
	out := new(ERC1967ProxyERC1967InvalidImplementation)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "ERC1967InvalidImplementation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyERC1967NonPayable represents a ERC1967NonPayable error raised by the ERC1967Proxy contract.
type ERC1967ProxyERC1967NonPayable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC1967NonPayable()
func ERC1967ProxyERC1967NonPayableErrorID() common.Hash {
	return common.HexToHash("0xb398979fa84f543c8e222f17890372c487baf85e062276c127fef521eea7224b")
}

// UnpackERC1967NonPayableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC1967NonPayable()
func (eRC1967Proxy *ERC1967Proxy) UnpackERC1967NonPayableError(raw []byte) (*ERC1967ProxyERC1967NonPayable, error) {
	out := new(ERC1967ProxyERC1967NonPayable)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "ERC1967NonPayable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ERC1967ProxyFailedCall represents a FailedCall error raised by the ERC1967Proxy contract.
type ERC1967ProxyFailedCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedCall()
func ERC1967ProxyFailedCallErrorID() common.Hash {
	return common.HexToHash("0xd6bda27508c0fb6d8a39b4b122878dab26f731a7d4e4abe711dd3731899052a4")
}

// UnpackFailedCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedCall()
func (eRC1967Proxy *ERC1967Proxy) UnpackFailedCallError(raw []byte) (*ERC1967ProxyFailedCall, error) {
	out := new(ERC1967ProxyFailedCall)
	if err := eRC1967Proxy.abi.UnpackIntoInterface(out, "FailedCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}
