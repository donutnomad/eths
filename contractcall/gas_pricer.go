package contractcall

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"math/big"
)

type gasCaller interface {
	IHeaderByNumber
	ethereum.GasPricer
	ethereum.GasPricer1559
}

type GasPricerDefault struct {
	client                  gasCaller
	BaseFeeWiggleMultiplier int64
}

func NewGasPricerDefault(client gasCaller) *GasPricerDefault {
	return &GasPricerDefault{client: client, BaseFeeWiggleMultiplier: 2}
}

func (p *GasPricerDefault) GetGasPrice(ctx context.Context, chainId *big.Int) (*GasPrice, error) {
	// Only query for basefee if gasPrice not specified
	head, err := p.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(EthereumRPCErr, err.Error())
	}
	if head.BaseFee != nil {
		tip, err := p.client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, errors.Wrap(EthereumRPCErr, err.Error())
		}
		gasFeeCap := new(big.Int).Add(tip, new(big.Int).Mul(head.BaseFee, big.NewInt(p.BaseFeeWiggleMultiplier)))
		return NewGasPrice(tip, gasFeeCap), nil
	} else {
		gasPrice, err := p.client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		return NewGasPriceLegacy(gasPrice), nil
	}
}

func isUseDynamicTx(ctx context.Context, chainId *big.Int, client IHeaderByNumber) (bool, error) {
	if ok := isUseDynamicTx1(chainId); ok {
		return true, nil
	}
	// Only query for basefee if gasPrice not specified
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return false, errors.Wrap(EthereumRPCErr, err.Error())
	}
	return head.BaseFee != nil, nil
}

func isUseDynamicTx1(chainId *big.Int) bool {
	// Optimism: 10 (OP Mainnet)
	// Ethereum: 1, Sepolia(11155111)
	// Arbitrum One: 42161
	// Arbitrum Nova: 42170
	chainIDUint64 := chainId.Uint64()
	return chainIDUint64 == 1 || chainIDUint64 == 10 || chainIDUint64 == 11155111 || chainIDUint64 == 42161 || chainIDUint64 == 42170
}
