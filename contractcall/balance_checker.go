package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type IBalanceChecker interface {
	CheckBalance(ctx context.Context, from common.Address, data []byte, to *common.Address, gasPrice *GasPrice) error
}

type BalanceCheckerImpl struct {
	client IBalance
}

func NewBalanceCheckerImpl(client IBalance) *BalanceCheckerImpl {
	return &BalanceCheckerImpl{client: client}
}

func (b *BalanceCheckerImpl) CheckBalance(ctx context.Context, from common.Address, data []byte, to *common.Address, gasPrice *GasPrice) error {
	balance, err := b.client.BalanceAt(ctx, from, nil)
	if err != nil {
		return err
	}
	if balance.Sign() <= 0 {
		return &InsufficientBalanceError{Balance: balance}
	}
	var IsIstanbul = true
	var IsShanghai = true
	intrGas, err := IntrinsicGas(data, nil, nil, to == nil, true, IsIstanbul, IsShanghai)
	if err != nil {
		// overflow uint64
	} else {
		if gasPrice.LegacyGas != nil {
			var maxGas = new(big.Int).Mul(
				new(big.Int).SetUint64(intrGas), /*gas limit*/
				gasPrice.LegacyGas.GasPrice,
			)
			if maxGas.Cmp(balance) == 1 {
				return &InsufficientBalanceError{Balance: balance}
			}
		} else if gasPrice.DynamicGas != nil {
			var maxGas = new(big.Int).Mul(
				new(big.Int).SetUint64(intrGas), /*gas limit*/
				gasPrice.DynamicGas.MaxFeePerGas,
			)
			if maxGas.Cmp(balance) == 1 {
				return &InsufficientBalanceError{Balance: balance}
			}
		}
	}
	return nil
}
