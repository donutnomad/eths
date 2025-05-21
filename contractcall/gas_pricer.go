package contractcall

import (
	"context"
	"github.com/pkg/errors"
	"math/big"
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
