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

// CreateXValues is an auto generated low-level Go binding around an user-defined struct.
type CreateXValues struct {
	ConstructorAmount *big.Int
	InitCallAmount    *big.Int
}

// CreatexMetaData contains all meta data concerning the Createx contract.
var CreatexMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"}],\"name\":\"FailedContractCreation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertData\",\"type\":\"bytes\"}],\"name\":\"FailedContractInitialisation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertData\",\"type\":\"bytes\"}],\"name\":\"FailedEtherTransfer\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"}],\"name\":\"InvalidNonceValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"}],\"name\":\"InvalidSalt\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"ContractCreation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractCreation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"Create3ProxyContractCreation\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"initCodeHash\",\"type\":\"bytes32\"}],\"name\":\"computeCreate2Address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"initCodeHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"computeCreate2Address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"computeCreate3Address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"computeCreate3Address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"computeCreateAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"computeCreateAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"computedAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deployCreate\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deployCreate2\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deployCreate2\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"}],\"name\":\"deployCreate2AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"}],\"name\":\"deployCreate2AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"}],\"name\":\"deployCreate2AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"}],\"name\":\"deployCreate2AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deployCreate2Clone\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deployCreate2Clone\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deployCreate3\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"}],\"name\":\"deployCreate3\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"}],\"name\":\"deployCreate3AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"}],\"name\":\"deployCreate3AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"}],\"name\":\"deployCreate3AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"}],\"name\":\"deployCreate3AndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"}],\"name\":\"deployCreateAndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"constructorAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initCallAmount\",\"type\":\"uint256\"}],\"internalType\":\"structCreateX.Values\",\"name\":\"values\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"refundAddress\",\"type\":\"address\"}],\"name\":\"deployCreateAndInit\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"deployCreateClone\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	ID:  "Createx",
}

// Createx is an auto generated Go binding around an Ethereum contract.
type Createx struct {
	abi abi.ABI
}

// NewCreatex creates a new instance of Createx.
func NewCreatex() *Createx {
	parsed, err := CreatexMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Createx{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Createx) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackComputeCreate2Address is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x890c283b.
//
// Solidity: function computeCreate2Address(bytes32 salt, bytes32 initCodeHash) view returns(address computedAddress)
func (createx *Createx) PackComputeCreate2Address(salt [32]byte, initCodeHash [32]byte) []byte {
	enc, err := createx.abi.Pack("computeCreate2Address", salt, initCodeHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreate2Address is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x890c283b.
//
// Solidity: function computeCreate2Address(bytes32 salt, bytes32 initCodeHash) view returns(address computedAddress)
func (createx *Createx) UnpackComputeCreate2Address(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreate2Address", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackComputeCreate2Address0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd323826a.
//
// Solidity: function computeCreate2Address(bytes32 salt, bytes32 initCodeHash, address deployer) pure returns(address computedAddress)
func (createx *Createx) PackComputeCreate2Address0(salt [32]byte, initCodeHash [32]byte, deployer common.Address) []byte {
	enc, err := createx.abi.Pack("computeCreate2Address0", salt, initCodeHash, deployer)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreate2Address0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd323826a.
//
// Solidity: function computeCreate2Address(bytes32 salt, bytes32 initCodeHash, address deployer) pure returns(address computedAddress)
func (createx *Createx) UnpackComputeCreate2Address0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreate2Address0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackComputeCreate3Address is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42d654fc.
//
// Solidity: function computeCreate3Address(bytes32 salt, address deployer) pure returns(address computedAddress)
func (createx *Createx) PackComputeCreate3Address(salt [32]byte, deployer common.Address) []byte {
	enc, err := createx.abi.Pack("computeCreate3Address", salt, deployer)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreate3Address is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42d654fc.
//
// Solidity: function computeCreate3Address(bytes32 salt, address deployer) pure returns(address computedAddress)
func (createx *Createx) UnpackComputeCreate3Address(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreate3Address", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackComputeCreate3Address0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6cec2536.
//
// Solidity: function computeCreate3Address(bytes32 salt) view returns(address computedAddress)
func (createx *Createx) PackComputeCreate3Address0(salt [32]byte) []byte {
	enc, err := createx.abi.Pack("computeCreate3Address0", salt)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreate3Address0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6cec2536.
//
// Solidity: function computeCreate3Address(bytes32 salt) view returns(address computedAddress)
func (createx *Createx) UnpackComputeCreate3Address0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreate3Address0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackComputeCreateAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x28ddd046.
//
// Solidity: function computeCreateAddress(uint256 nonce) view returns(address computedAddress)
func (createx *Createx) PackComputeCreateAddress(nonce *big.Int) []byte {
	enc, err := createx.abi.Pack("computeCreateAddress", nonce)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreateAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x28ddd046.
//
// Solidity: function computeCreateAddress(uint256 nonce) view returns(address computedAddress)
func (createx *Createx) UnpackComputeCreateAddress(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreateAddress", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackComputeCreateAddress0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74637a7a.
//
// Solidity: function computeCreateAddress(address deployer, uint256 nonce) view returns(address computedAddress)
func (createx *Createx) PackComputeCreateAddress0(deployer common.Address, nonce *big.Int) []byte {
	enc, err := createx.abi.Pack("computeCreateAddress0", deployer, nonce)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackComputeCreateAddress0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x74637a7a.
//
// Solidity: function computeCreateAddress(address deployer, uint256 nonce) view returns(address computedAddress)
func (createx *Createx) UnpackComputeCreateAddress0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("computeCreateAddress0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27fe1822.
//
// Solidity: function deployCreate(bytes initCode) payable returns(address newContract)
func (createx *Createx) PackDeployCreate(initCode []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate", initCode)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27fe1822.
//
// Solidity: function deployCreate(bytes initCode) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x26307668.
//
// Solidity: function deployCreate2(bytes32 salt, bytes initCode) payable returns(address newContract)
func (createx *Createx) PackDeployCreate2(salt [32]byte, initCode []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate2", salt, initCode)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x26307668.
//
// Solidity: function deployCreate2(bytes32 salt, bytes initCode) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate2(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate20 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x26a32fc7.
//
// Solidity: function deployCreate2(bytes initCode) payable returns(address newContract)
func (createx *Createx) PackDeployCreate20(initCode []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate20", initCode)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate20 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x26a32fc7.
//
// Solidity: function deployCreate2(bytes initCode) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate20(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate20", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2AndInit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa7db93f2.
//
// Solidity: function deployCreate2AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) PackDeployCreate2AndInit(salt [32]byte, initCode []byte, data []byte, values CreateXValues, refundAddress common.Address) []byte {
	enc, err := createx.abi.Pack("deployCreate2AndInit", salt, initCode, data, values, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2AndInit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa7db93f2.
//
// Solidity: function deployCreate2AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate2AndInit(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2AndInit", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2AndInit0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3fe107b.
//
// Solidity: function deployCreate2AndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) PackDeployCreate2AndInit0(initCode []byte, data []byte, values CreateXValues) []byte {
	enc, err := createx.abi.Pack("deployCreate2AndInit0", initCode, data, values)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2AndInit0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3fe107b.
//
// Solidity: function deployCreate2AndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate2AndInit0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2AndInit0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2AndInit1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe437252a.
//
// Solidity: function deployCreate2AndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) PackDeployCreate2AndInit1(initCode []byte, data []byte, values CreateXValues, refundAddress common.Address) []byte {
	enc, err := createx.abi.Pack("deployCreate2AndInit1", initCode, data, values, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2AndInit1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe437252a.
//
// Solidity: function deployCreate2AndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate2AndInit1(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2AndInit1", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2AndInit2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe96deee4.
//
// Solidity: function deployCreate2AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) PackDeployCreate2AndInit2(salt [32]byte, initCode []byte, data []byte, values CreateXValues) []byte {
	enc, err := createx.abi.Pack("deployCreate2AndInit2", salt, initCode, data, values)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2AndInit2 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe96deee4.
//
// Solidity: function deployCreate2AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate2AndInit2(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2AndInit2", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2Clone is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2852527a.
//
// Solidity: function deployCreate2Clone(bytes32 salt, address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) PackDeployCreate2Clone(salt [32]byte, implementation common.Address, data []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate2Clone", salt, implementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2Clone is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2852527a.
//
// Solidity: function deployCreate2Clone(bytes32 salt, address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) UnpackDeployCreate2Clone(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2Clone", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate2Clone0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81503da1.
//
// Solidity: function deployCreate2Clone(address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) PackDeployCreate2Clone0(implementation common.Address, data []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate2Clone0", implementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate2Clone0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x81503da1.
//
// Solidity: function deployCreate2Clone(address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) UnpackDeployCreate2Clone0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate2Clone0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7f565360.
//
// Solidity: function deployCreate3(bytes initCode) payable returns(address newContract)
func (createx *Createx) PackDeployCreate3(initCode []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate3", initCode)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate3 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7f565360.
//
// Solidity: function deployCreate3(bytes initCode) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate3(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate3", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate30 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9c36a286.
//
// Solidity: function deployCreate3(bytes32 salt, bytes initCode) payable returns(address newContract)
func (createx *Createx) PackDeployCreate30(salt [32]byte, initCode []byte) []byte {
	enc, err := createx.abi.Pack("deployCreate30", salt, initCode)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate30 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9c36a286.
//
// Solidity: function deployCreate3(bytes32 salt, bytes initCode) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate30(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate30", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate3AndInit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00d84acb.
//
// Solidity: function deployCreate3AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) PackDeployCreate3AndInit(salt [32]byte, initCode []byte, data []byte, values CreateXValues) []byte {
	enc, err := createx.abi.Pack("deployCreate3AndInit", salt, initCode, data, values)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate3AndInit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00d84acb.
//
// Solidity: function deployCreate3AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate3AndInit(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate3AndInit", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate3AndInit0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f990e3f.
//
// Solidity: function deployCreate3AndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) PackDeployCreate3AndInit0(initCode []byte, data []byte, values CreateXValues) []byte {
	enc, err := createx.abi.Pack("deployCreate3AndInit0", initCode, data, values)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate3AndInit0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2f990e3f.
//
// Solidity: function deployCreate3AndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate3AndInit0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate3AndInit0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate3AndInit1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xddda0acb.
//
// Solidity: function deployCreate3AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) PackDeployCreate3AndInit1(salt [32]byte, initCode []byte, data []byte, values CreateXValues, refundAddress common.Address) []byte {
	enc, err := createx.abi.Pack("deployCreate3AndInit1", salt, initCode, data, values, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate3AndInit1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xddda0acb.
//
// Solidity: function deployCreate3AndInit(bytes32 salt, bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate3AndInit1(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate3AndInit1", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreate3AndInit2 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5745aba.
//
// Solidity: function deployCreate3AndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) PackDeployCreate3AndInit2(initCode []byte, data []byte, values CreateXValues, refundAddress common.Address) []byte {
	enc, err := createx.abi.Pack("deployCreate3AndInit2", initCode, data, values, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreate3AndInit2 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5745aba.
//
// Solidity: function deployCreate3AndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreate3AndInit2(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreate3AndInit2", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreateAndInit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31a7c8c8.
//
// Solidity: function deployCreateAndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) PackDeployCreateAndInit(initCode []byte, data []byte, values CreateXValues) []byte {
	enc, err := createx.abi.Pack("deployCreateAndInit", initCode, data, values)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreateAndInit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x31a7c8c8.
//
// Solidity: function deployCreateAndInit(bytes initCode, bytes data, (uint256,uint256) values) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreateAndInit(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreateAndInit", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreateAndInit0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x98e81077.
//
// Solidity: function deployCreateAndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) PackDeployCreateAndInit0(initCode []byte, data []byte, values CreateXValues, refundAddress common.Address) []byte {
	enc, err := createx.abi.Pack("deployCreateAndInit0", initCode, data, values, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreateAndInit0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x98e81077.
//
// Solidity: function deployCreateAndInit(bytes initCode, bytes data, (uint256,uint256) values, address refundAddress) payable returns(address newContract)
func (createx *Createx) UnpackDeployCreateAndInit0(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreateAndInit0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackDeployCreateClone is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf9664498.
//
// Solidity: function deployCreateClone(address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) PackDeployCreateClone(implementation common.Address, data []byte) []byte {
	enc, err := createx.abi.Pack("deployCreateClone", implementation, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackDeployCreateClone is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf9664498.
//
// Solidity: function deployCreateClone(address implementation, bytes data) payable returns(address proxy)
func (createx *Createx) UnpackDeployCreateClone(data []byte) (common.Address, error) {
	out, err := createx.abi.Unpack("deployCreateClone", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// CreatexContractCreation represents a ContractCreation event raised by the Createx contract.
type CreatexContractCreation struct {
	NewContract common.Address
	Salt        [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const CreatexContractCreationEventName = "ContractCreation"

// ContractEventName returns the user-defined event name.
func (CreatexContractCreation) ContractEventName() string {
	return CreatexContractCreationEventName
}

// UnpackContractCreationEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ContractCreation(address indexed newContract, bytes32 indexed salt)
func (createx *Createx) UnpackContractCreationEvent(log *types.Log) (*CreatexContractCreation, error) {
	event := "ContractCreation"
	if log.Topics[0] != createx.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CreatexContractCreation)
	if len(log.Data) > 0 {
		if err := createx.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range createx.abi.Events[event].Inputs {
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

// CreatexContractCreation0 represents a ContractCreation0 event raised by the Createx contract.
type CreatexContractCreation0 struct {
	NewContract common.Address
	Raw         *types.Log // Blockchain specific contextual infos
}

const CreatexContractCreation0EventName = "ContractCreation0"

// ContractEventName returns the user-defined event name.
func (CreatexContractCreation0) ContractEventName() string {
	return CreatexContractCreation0EventName
}

// UnpackContractCreation0Event is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ContractCreation(address indexed newContract)
func (createx *Createx) UnpackContractCreation0Event(log *types.Log) (*CreatexContractCreation0, error) {
	event := "ContractCreation0"
	if log.Topics[0] != createx.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CreatexContractCreation0)
	if len(log.Data) > 0 {
		if err := createx.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range createx.abi.Events[event].Inputs {
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

// CreatexCreate3ProxyContractCreation represents a Create3ProxyContractCreation event raised by the Createx contract.
type CreatexCreate3ProxyContractCreation struct {
	NewContract common.Address
	Salt        [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const CreatexCreate3ProxyContractCreationEventName = "Create3ProxyContractCreation"

// ContractEventName returns the user-defined event name.
func (CreatexCreate3ProxyContractCreation) ContractEventName() string {
	return CreatexCreate3ProxyContractCreationEventName
}

// UnpackCreate3ProxyContractCreationEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Create3ProxyContractCreation(address indexed newContract, bytes32 indexed salt)
func (createx *Createx) UnpackCreate3ProxyContractCreationEvent(log *types.Log) (*CreatexCreate3ProxyContractCreation, error) {
	event := "Create3ProxyContractCreation"
	if log.Topics[0] != createx.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CreatexCreate3ProxyContractCreation)
	if len(log.Data) > 0 {
		if err := createx.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range createx.abi.Events[event].Inputs {
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
func (createx *Createx) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], createx.abi.Errors["FailedContractCreation"].ID.Bytes()[:4]) {
		return createx.UnpackFailedContractCreationError(raw[4:])
	}
	if bytes.Equal(raw[:4], createx.abi.Errors["FailedContractInitialisation"].ID.Bytes()[:4]) {
		return createx.UnpackFailedContractInitialisationError(raw[4:])
	}
	if bytes.Equal(raw[:4], createx.abi.Errors["FailedEtherTransfer"].ID.Bytes()[:4]) {
		return createx.UnpackFailedEtherTransferError(raw[4:])
	}
	if bytes.Equal(raw[:4], createx.abi.Errors["InvalidNonceValue"].ID.Bytes()[:4]) {
		return createx.UnpackInvalidNonceValueError(raw[4:])
	}
	if bytes.Equal(raw[:4], createx.abi.Errors["InvalidSalt"].ID.Bytes()[:4]) {
		return createx.UnpackInvalidSaltError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// CreatexFailedContractCreation represents a FailedContractCreation error raised by the Createx contract.
type CreatexFailedContractCreation struct {
	Emitter common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedContractCreation(address emitter)
func CreatexFailedContractCreationErrorID() common.Hash {
	return common.HexToHash("0xc05cee7adec1c7022c70b91bddcde5124ca9bd2894bcd20bdbeb98c4ccd6ad31")
}

// UnpackFailedContractCreationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedContractCreation(address emitter)
func (createx *Createx) UnpackFailedContractCreationError(raw []byte) (*CreatexFailedContractCreation, error) {
	out := new(CreatexFailedContractCreation)
	if err := createx.abi.UnpackIntoInterface(out, "FailedContractCreation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CreatexFailedContractInitialisation represents a FailedContractInitialisation error raised by the Createx contract.
type CreatexFailedContractInitialisation struct {
	Emitter    common.Address
	RevertData []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedContractInitialisation(address emitter, bytes revertData)
func CreatexFailedContractInitialisationErrorID() common.Hash {
	return common.HexToHash("0xa57ca239dc21ebdb895858cd57c414f9c89f18ea5c815cb1e329c666d45236f0")
}

// UnpackFailedContractInitialisationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedContractInitialisation(address emitter, bytes revertData)
func (createx *Createx) UnpackFailedContractInitialisationError(raw []byte) (*CreatexFailedContractInitialisation, error) {
	out := new(CreatexFailedContractInitialisation)
	if err := createx.abi.UnpackIntoInterface(out, "FailedContractInitialisation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CreatexFailedEtherTransfer represents a FailedEtherTransfer error raised by the Createx contract.
type CreatexFailedEtherTransfer struct {
	Emitter    common.Address
	RevertData []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FailedEtherTransfer(address emitter, bytes revertData)
func CreatexFailedEtherTransferErrorID() common.Hash {
	return common.HexToHash("0xc2b3f4452c5ac36c715121b95d78a40ac33806494b2975a8238b27da8a77e1e1")
}

// UnpackFailedEtherTransferError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FailedEtherTransfer(address emitter, bytes revertData)
func (createx *Createx) UnpackFailedEtherTransferError(raw []byte) (*CreatexFailedEtherTransfer, error) {
	out := new(CreatexFailedEtherTransfer)
	if err := createx.abi.UnpackIntoInterface(out, "FailedEtherTransfer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CreatexInvalidNonceValue represents a InvalidNonceValue error raised by the Createx contract.
type CreatexInvalidNonceValue struct {
	Emitter common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidNonceValue(address emitter)
func CreatexInvalidNonceValueErrorID() common.Hash {
	return common.HexToHash("0x3c55ab3b3cc44087e1906d945b58a9ee2cdac44c0773594f81c831027ddd8bc4")
}

// UnpackInvalidNonceValueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidNonceValue(address emitter)
func (createx *Createx) UnpackInvalidNonceValueError(raw []byte) (*CreatexInvalidNonceValue, error) {
	out := new(CreatexInvalidNonceValue)
	if err := createx.abi.UnpackIntoInterface(out, "InvalidNonceValue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CreatexInvalidSalt represents a InvalidSalt error raised by the Createx contract.
type CreatexInvalidSalt struct {
	Emitter common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSalt(address emitter)
func CreatexInvalidSaltErrorID() common.Hash {
	return common.HexToHash("0x13b3a2a19cc002fe27dc4952e92fb58eb225aa1ce015e59c8ba9b607a2163fe9")
}

// UnpackInvalidSaltError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSalt(address emitter)
func (createx *Createx) UnpackInvalidSaltError(raw []byte) (*CreatexInvalidSalt, error) {
	out := new(CreatexInvalidSalt)
	if err := createx.abi.UnpackIntoInterface(out, "InvalidSalt", raw); err != nil {
		return nil, err
	}
	return out, nil
}
