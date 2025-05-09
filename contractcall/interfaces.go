package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum"
)

// GasPrice represents the gas price configuration
type GasPrice struct {
	LegacyGas  *LegacyGas
	DynamicGas *DynamicGas
}

func NewGasPriceLegacy(price *big.Int) *GasPrice {
	return &GasPrice{LegacyGas: &LegacyGas{GasPrice: price}}
}
func NewGasPrice(maxPriorityFeePerGas, maxFeePerGas *big.Int) *GasPrice {
	return &GasPrice{DynamicGas: &DynamicGas{MaxPriorityFeePerGas: maxPriorityFeePerGas, MaxFeePerGas: maxFeePerGas}}
}

// LegacyGas represents traditional gas price
type LegacyGas struct {
	GasPrice *big.Int
}

// DynamicGas represents EIP-1559 gas price
type DynamicGas struct {
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

type IGasPricer interface {
	GetGasPrice(ctx context.Context, chainId *big.Int) (*GasPrice, error)
}

type IGasPriceValidator interface {
	ValidateGasPrice(ctx context.Context, chainId *big.Int, gasPrice *GasPrice) error
}

// IEstimateGas 计算gas总量
// OP:
// 估算执行 L2 交易所需的L1 数据 gas + L2 gas的数量。
// estimateL1Gas它是(L1 Gas) 和estimateGas(L2 Gas)的总和。
// Arb:
type IEstimateGas interface {
	EstimateGas(ctx context.Context, chainId *big.Int, msg ethereum.CallMsg) (*big.Int, error)
}

type IGetNonce interface {
	GetNonce(ctx context.Context, account common.Address, isPending bool) (uint64, error)
}

type IHeaderByNumber interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
}
type IPendingNonceAt interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}
