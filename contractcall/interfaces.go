package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type IEthereumRPC interface {
	ethereum.ChainIDReader
	ethereum.BlockNumberReader
	ethereum.ContractCaller
	ethereum.PendingContractCaller
	ethereum.GasEstimator
	ethereum.PendingStateReader
	ethereum.GasPricer
	ethereum.GasPricer1559

	ethereum.TransactionSender
	ethereum.LogFilterer

	ethereum.TransactionReader
	ethereum.ChainStateReader
	ethereum.ChainReader
}

type gasCaller interface {
	IHeaderByNumber
	ethereum.GasPricer
	ethereum.GasPricer1559
}

type IMyClient interface {
	IHeaderByNumber
	ethereum.GasPricer
	ethereum.GasPricer1559
	ethereum.GasEstimator
	INonceAt
	IBalance
}

type ISendTxClient interface {
	ethereum.TransactionSender
	ICodeAt
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

type IBalance interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}
