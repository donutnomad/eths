package contractcall

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func covertGweiToWei(gwei string) *big.Int {
	fromString, err := decimal.NewFromString(gwei)
	if err != nil {
		return new(big.Int).SetUint64(0)
	}
	fromString = fromString.Mul(decimal.New(1, 9))
	return fromString.BigInt()
}

type ILogger interface {
	Printf(string, ...interface{})
}

func logCallMsg(logger ILogger, msg *ethereum.CallMsg) {
	if logger == nil || msg == nil {
		return
	}
	logger.Printf(fmt.Sprintf("estimateGasLimit: gasPrice:%s, gasTipCap:%s, gasFeeCap:%s\n",
		bigIntToString(msg.GasPrice),
		bigIntToString(msg.GasTipCap),
		bigIntToString(msg.GasFeeCap)),
	)
	logger.Printf(fmt.Sprintf("from: %v, to:%v, value:%v, data:%x\n", msg.From, msg.To, msg.Value, msg.Data))
}

func bigIntToString(input *big.Int) string {
	if input == nil {
		return "<nil>"
	}
	return input.String()
}

type KnownError struct {
	Name      string
	Signature string
	Arguments []string
}

var ErrGasUintOverflow = errors.New("gas uint64 overflow")

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, accessList ethTypes.AccessList, authList []ethTypes.SetCodeAuthorization, isContractCreation, isHomestead, isEIP2028, isEIP3860 bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isContractCreation && isHomestead {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	dataLen := uint64(len(data))
	// Bump the required gas by the amount of transactional data
	if dataLen > 0 {
		// Zero and non-zero bytes are priced differently
		z := uint64(bytes.Count(data, []byte{0}))
		nz := dataLen - z

		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := params.TxDataNonZeroGasFrontier
		if isEIP2028 {
			nonZeroGas = params.TxDataNonZeroGasEIP2028
		}
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, ErrGasUintOverflow
		}
		gas += nz * nonZeroGas

		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * params.TxDataZeroGas

		if isContractCreation && isEIP3860 {
			lenWords := toWordSize(dataLen)
			if (math.MaxUint64-gas)/params.InitCodeWordGas < lenWords {
				return 0, ErrGasUintOverflow
			}
			gas += lenWords * params.InitCodeWordGas
		}
	}
	if accessList != nil {
		gas += uint64(len(accessList)) * params.TxAccessListAddressGas
		gas += uint64(accessList.StorageKeys()) * params.TxAccessListStorageKeyGas
	}
	if authList != nil {
		gas += uint64(len(authList)) * params.CallNewAccountGas
	}
	return gas, nil
}

// toWordSize returns the ceiled word size required for init code payment calculation.
func toWordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

func bigToBytes32[T *big.Int | *uint256.Int](i T) [32]byte {
	var bs [32]byte
	if lo.IsNil(i) {
		return bs
	}
	switch v := any(i).(type) {
	case *big.Int:
		v.FillBytes(bs[:])
		return bs
	case *uint256.Int:
		v.ToBig().FillBytes(bs[:])
		return bs
	default:
		panic(UNREACHABLE)
	}
}

func copyInt(i *big.Int) *big.Int {
	if i == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(i)
}

func newInt(i *big.Int) *uint256.Int {
	if i == nil {
		return uint256.NewInt(0)
	}
	d := new(uint256.Int)
	d.SetFromBig(i)
	return d
}

func newIntBy[I constraints.Integer](i I) *uint256.Int {
	d := new(uint256.Int)
	d.SetUint64(uint64(i))
	return d
}

type ethTransactionReflect struct {
	Inner ethTypes.TxData
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	vCopy := new(big.Int).Sub(v, big.NewInt(35))
	return vCopy.Rsh(vCopy, 1)
}

func bigIntOrIntToBigInt[N *big.Int | *uint256.Int | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](num N) *big.Int {
	switch v := any(num).(type) {
	case *big.Int:
		return v
	case *uint256.Int:
		return v.ToBig()
	case uint:
		return new(big.Int).SetUint64(uint64(v))
	case uint8:
		return new(big.Int).SetUint64(uint64(v))
	case uint16:
		return new(big.Int).SetUint64(uint64(v))
	case uint32:
		return new(big.Int).SetUint64(uint64(v))
	case uint64:
		return new(big.Int).SetUint64(v)
	case int:
		return new(big.Int).SetInt64(int64(v))
	case int8:
		return new(big.Int).SetInt64(int64(v))
	case int16:
		return new(big.Int).SetInt64(int64(v))
	case int32:
		return new(big.Int).SetInt64(int64(v))
	case int64:
		return new(big.Int).SetInt64(v)
	default:
		panic(UNREACHABLE)
	}
}
