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

// Multicall3Call is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Call struct {
	Target   common.Address
	CallData []byte
}

// Multicall3Call3 is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Call3 struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

// Multicall3Call3Value is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Call3Value struct {
	Target       common.Address
	AllowFailure bool
	Value        *big.Int
	CallData     []byte
}

// Multicall3Result is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Result struct {
	Success    bool
	ReturnData []byte
}

// MulticallMetaData contains all meta data concerning the Multicall contract.
var MulticallMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowFailure\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call3[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate3\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowFailure\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call3Value[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate3Value\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"blockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBasefee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"basefee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainid\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockCoinbase\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"coinbase\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"difficulty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockGasLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gaslimit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryAggregate\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"requireSuccess\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"tryBlockAndAggregate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structMulticall3.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	ID:  "Multicall",
	Bin: "0x6080604052348015600f57600080fd5b5061039b8061001f6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806301b069c814610030575b600080fd5b61004361003e3660046101c7565b610059565b604051610050919061023e565b60405180910390f35b6060818067ffffffffffffffff811115610075576100756102db565b6040519080825280602002602001820160405280156100a857816020015b60608152602001906001900390816100935790505b50915060005b818110156101bf576000308686848181106100cb576100cb6102f1565b90506020028101906100dd9190610307565b6040516100eb929190610355565b600060405180830381855af49150503d8060008114610126576040519150601f19603f3d011682016040523d82523d6000602084013e61012b565b606091505b5085848151811061013e5761013e6102f1565b60209081029190910101529050806101b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f4d756c746963616c6c333a2063616c6c206661696c6564000000000000000000604482015260640160405180910390fd5b506001016100ae565b505092915050565b600080602083850312156101da57600080fd5b823567ffffffffffffffff8111156101f157600080fd5b8301601f8101851361020257600080fd5b803567ffffffffffffffff81111561021957600080fd5b8560208260051b840101111561022e57600080fd5b6020919091019590945092505050565b6000602082016020835280845180835260408501915060408160051b86010192506020860160005b828110156102cf57868503603f190184528151805180875260005b8181101561029d57602081840181015189830182015201610281565b506000602082890101526020601f19601f83011688010196505050602082019150602084019350600181019050610266565b50929695505050505050565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000808335601e1984360301811261031e57600080fd5b83018035915067ffffffffffffffff82111561033957600080fd5b60200191503681900382131561034e57600080fd5b9250929050565b818382376000910190815291905056fea26469706673582212202a1202ead93629aec820f9e7766d7afa052025380ab1a4c9d4ac27a2c00a286e64736f6c634300081c0033}",
}

// Multicall is an auto generated Go binding around an Ethereum contract.
type Multicall struct {
	abi abi.ABI
}

// NewMulticall creates a new instance of Multicall.
func NewMulticall() *Multicall {
	parsed, err := MulticallMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Multicall{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Multicall) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAggregate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x252dba42.
//
// Solidity: function aggregate((address,bytes)[] calls) payable returns(uint256 blockNumber, bytes[] returnData)
func (multicall *Multicall) PackAggregate(calls []Multicall3Call) []byte {
	enc, err := multicall.abi.Pack("aggregate", calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// AggregateOutput serves as a container for the return parameters of contract
// method Aggregate.
type AggregateOutput struct {
	BlockNumber *big.Int
	ReturnData  [][]byte
}

// UnpackAggregate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x252dba42.
//
// Solidity: function aggregate((address,bytes)[] calls) payable returns(uint256 blockNumber, bytes[] returnData)
func (multicall *Multicall) UnpackAggregate(data []byte) (AggregateOutput, error) {
	out, err := multicall.abi.Unpack("aggregate", data)
	outstruct := new(AggregateOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.BlockNumber = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.ReturnData = *abi.ConvertType(out[1], new([][]byte)).(*[][]byte)
	return *outstruct, err

}

// PackAggregate3 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82ad56cb.
//
// Solidity: function aggregate3((address,bool,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) PackAggregate3(calls []Multicall3Call3) []byte {
	enc, err := multicall.abi.Pack("aggregate3", calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAggregate3 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x82ad56cb.
//
// Solidity: function aggregate3((address,bool,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) UnpackAggregate3(data []byte) ([]Multicall3Result, error) {
	out, err := multicall.abi.Unpack("aggregate3", data)
	if err != nil {
		return *new([]Multicall3Result), err
	}
	out0 := *abi.ConvertType(out[0], new([]Multicall3Result)).(*[]Multicall3Result)
	return out0, err
}

// PackAggregate3Value is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x174dea71.
//
// Solidity: function aggregate3Value((address,bool,uint256,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) PackAggregate3Value(calls []Multicall3Call3Value) []byte {
	enc, err := multicall.abi.Pack("aggregate3Value", calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackAggregate3Value is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x174dea71.
//
// Solidity: function aggregate3Value((address,bool,uint256,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) UnpackAggregate3Value(data []byte) ([]Multicall3Result, error) {
	out, err := multicall.abi.Unpack("aggregate3Value", data)
	if err != nil {
		return *new([]Multicall3Result), err
	}
	out0 := *abi.ConvertType(out[0], new([]Multicall3Result)).(*[]Multicall3Result)
	return out0, err
}

// PackBlockAndAggregate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3077fa9.
//
// Solidity: function blockAndAggregate((address,bytes)[] calls) payable returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (multicall *Multicall) PackBlockAndAggregate(calls []Multicall3Call) []byte {
	enc, err := multicall.abi.Pack("blockAndAggregate", calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// BlockAndAggregateOutput serves as a container for the return parameters of contract
// method BlockAndAggregate.
type BlockAndAggregateOutput struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	ReturnData  []Multicall3Result
}

// UnpackBlockAndAggregate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3077fa9.
//
// Solidity: function blockAndAggregate((address,bytes)[] calls) payable returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (multicall *Multicall) UnpackBlockAndAggregate(data []byte) (BlockAndAggregateOutput, error) {
	out, err := multicall.abi.Unpack("blockAndAggregate", data)
	outstruct := new(BlockAndAggregateOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.BlockNumber = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.BlockHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ReturnData = *abi.ConvertType(out[2], new([]Multicall3Result)).(*[]Multicall3Result)
	return *outstruct, err

}

// PackGetBasefee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e64a696.
//
// Solidity: function getBasefee() view returns(uint256 basefee)
func (multicall *Multicall) PackGetBasefee() []byte {
	enc, err := multicall.abi.Pack("getBasefee")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetBasefee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3e64a696.
//
// Solidity: function getBasefee() view returns(uint256 basefee)
func (multicall *Multicall) UnpackGetBasefee(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getBasefee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetBlockHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 blockNumber) view returns(bytes32 blockHash)
func (multicall *Multicall) PackGetBlockHash(blockNumber *big.Int) []byte {
	enc, err := multicall.abi.Pack("getBlockHash", blockNumber)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetBlockHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 blockNumber) view returns(bytes32 blockHash)
func (multicall *Multicall) UnpackGetBlockHash(data []byte) ([32]byte, error) {
	out, err := multicall.abi.Unpack("getBlockHash", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackGetBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256 blockNumber)
func (multicall *Multicall) PackGetBlockNumber() []byte {
	enc, err := multicall.abi.Pack("getBlockNumber")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256 blockNumber)
func (multicall *Multicall) UnpackGetBlockNumber(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetChainId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256 chainid)
func (multicall *Multicall) PackGetChainId() []byte {
	enc, err := multicall.abi.Pack("getChainId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetChainId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3408e470.
//
// Solidity: function getChainId() view returns(uint256 chainid)
func (multicall *Multicall) UnpackGetChainId(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getChainId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCurrentBlockCoinbase is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa8b0574e.
//
// Solidity: function getCurrentBlockCoinbase() view returns(address coinbase)
func (multicall *Multicall) PackGetCurrentBlockCoinbase() []byte {
	enc, err := multicall.abi.Pack("getCurrentBlockCoinbase")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentBlockCoinbase is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa8b0574e.
//
// Solidity: function getCurrentBlockCoinbase() view returns(address coinbase)
func (multicall *Multicall) UnpackGetCurrentBlockCoinbase(data []byte) (common.Address, error) {
	out, err := multicall.abi.Unpack("getCurrentBlockCoinbase", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackGetCurrentBlockDifficulty is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72425d9d.
//
// Solidity: function getCurrentBlockDifficulty() view returns(uint256 difficulty)
func (multicall *Multicall) PackGetCurrentBlockDifficulty() []byte {
	enc, err := multicall.abi.Pack("getCurrentBlockDifficulty")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentBlockDifficulty is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x72425d9d.
//
// Solidity: function getCurrentBlockDifficulty() view returns(uint256 difficulty)
func (multicall *Multicall) UnpackGetCurrentBlockDifficulty(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getCurrentBlockDifficulty", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCurrentBlockGasLimit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86d516e8.
//
// Solidity: function getCurrentBlockGasLimit() view returns(uint256 gaslimit)
func (multicall *Multicall) PackGetCurrentBlockGasLimit() []byte {
	enc, err := multicall.abi.Pack("getCurrentBlockGasLimit")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentBlockGasLimit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86d516e8.
//
// Solidity: function getCurrentBlockGasLimit() view returns(uint256 gaslimit)
func (multicall *Multicall) UnpackGetCurrentBlockGasLimit(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getCurrentBlockGasLimit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetCurrentBlockTimestamp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (multicall *Multicall) PackGetCurrentBlockTimestamp() []byte {
	enc, err := multicall.abi.Pack("getCurrentBlockTimestamp")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentBlockTimestamp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (multicall *Multicall) UnpackGetCurrentBlockTimestamp(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getCurrentBlockTimestamp", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetEthBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (multicall *Multicall) PackGetEthBalance(addr common.Address) []byte {
	enc, err := multicall.abi.Pack("getEthBalance", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetEthBalance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (multicall *Multicall) UnpackGetEthBalance(data []byte) (*big.Int, error) {
	out, err := multicall.abi.Unpack("getEthBalance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}

// PackGetLastBlockHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27e86d6e.
//
// Solidity: function getLastBlockHash() view returns(bytes32 blockHash)
func (multicall *Multicall) PackGetLastBlockHash() []byte {
	enc, err := multicall.abi.Pack("getLastBlockHash")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetLastBlockHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27e86d6e.
//
// Solidity: function getLastBlockHash() view returns(bytes32 blockHash)
func (multicall *Multicall) UnpackGetLastBlockHash(data []byte) ([32]byte, error) {
	out, err := multicall.abi.Unpack("getLastBlockHash", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, err
}

// PackTryAggregate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbce38bd7.
//
// Solidity: function tryAggregate(bool requireSuccess, (address,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) PackTryAggregate(requireSuccess bool, calls []Multicall3Call) []byte {
	enc, err := multicall.abi.Pack("tryAggregate", requireSuccess, calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackTryAggregate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbce38bd7.
//
// Solidity: function tryAggregate(bool requireSuccess, (address,bytes)[] calls) payable returns((bool,bytes)[] returnData)
func (multicall *Multicall) UnpackTryAggregate(data []byte) ([]Multicall3Result, error) {
	out, err := multicall.abi.Unpack("tryAggregate", data)
	if err != nil {
		return *new([]Multicall3Result), err
	}
	out0 := *abi.ConvertType(out[0], new([]Multicall3Result)).(*[]Multicall3Result)
	return out0, err
}

// PackTryBlockAndAggregate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x399542e9.
//
// Solidity: function tryBlockAndAggregate(bool requireSuccess, (address,bytes)[] calls) payable returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (multicall *Multicall) PackTryBlockAndAggregate(requireSuccess bool, calls []Multicall3Call) []byte {
	enc, err := multicall.abi.Pack("tryBlockAndAggregate", requireSuccess, calls)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryBlockAndAggregateOutput serves as a container for the return parameters of contract
// method TryBlockAndAggregate.
type TryBlockAndAggregateOutput struct {
	BlockNumber *big.Int
	BlockHash   [32]byte
	ReturnData  []Multicall3Result
}

// UnpackTryBlockAndAggregate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x399542e9.
//
// Solidity: function tryBlockAndAggregate(bool requireSuccess, (address,bytes)[] calls) payable returns(uint256 blockNumber, bytes32 blockHash, (bool,bytes)[] returnData)
func (multicall *Multicall) UnpackTryBlockAndAggregate(data []byte) (TryBlockAndAggregateOutput, error) {
	out, err := multicall.abi.Unpack("tryBlockAndAggregate", data)
	outstruct := new(TryBlockAndAggregateOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.BlockNumber = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.BlockHash = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ReturnData = *abi.ConvertType(out[2], new([]Multicall3Result)).(*[]Multicall3Result)
	return *outstruct, err

}
