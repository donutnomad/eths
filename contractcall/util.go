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
	"github.com/shopspring/decimal"
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
