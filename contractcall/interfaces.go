package contractcall

import (
	"context"
	"math/big"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethclient"
	"github.com/donutnomad/eths/ethtype"
)

type IEthereumRPC interface {
	ethclient.ChainIDReader
	ethclient.BlockNumberReader
	ethclient.ContractCaller
	ethclient.PendingContractCaller
	ethclient.GasEstimator
	ethclient.PendingStateReader
	ethclient.GasPricer
	ethclient.GasPricer1559

	ethclient.TransactionSender
	ethclient.LogFilterer

	ethclient.TransactionReader
	ethclient.ChainStateReader
	ethclient.ChainReader
}

type gasCaller interface {
	IHeaderByNumber
	ethclient.GasPricer
	ethclient.GasPricer1559
}

type IMyClient interface {
	IHeaderByNumber
	ethclient.GasPricer
	ethclient.GasPricer1559
	ethclient.GasEstimator
	INonceAt
	IBalance
}

type ISendTxClient interface {
	ethclient.TransactionSender
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
	EstimateGas(ctx context.Context, chainId *big.Int, msg ethclient.CallMsg) (*big.Int, error)
}

type IGetNonce interface {
	GetNonce(ctx context.Context, account common.Address, isPending bool) (uint64, error)
}

type IHeaderByNumber interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethtype.EHeader, error)
}

type IBalance interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}
