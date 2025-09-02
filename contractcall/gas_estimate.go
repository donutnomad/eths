package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

type GasEstimateImpl struct {
	client ethereum.GasEstimator
	logger ILogger
}

func NewGasEstimateImpl(client ethereum.GasEstimator, logger ILogger) *GasEstimateImpl {
	return &GasEstimateImpl{client: client, logger: logger}
}

func (i *GasEstimateImpl) EstimateGas(ctx context.Context, chainId *big.Int, msg ethereum.CallMsg) (*big.Int, error) {
	logCallMsg(i.logger, &msg)
	gas, err := i.client.EstimateGas(ctx, msg)
	if err != nil {
		return big.NewInt(0), ParseEvmError(err)
	}
	return new(big.Int).SetUint64(gas), nil
}

type FixedGasLimit struct {
	Value *big.Int
}

func (i *FixedGasLimit) EstimateGas(_ context.Context, _ *big.Int, _ ethereum.CallMsg) (*big.Int, error) {
	return i.Value, nil
}

type GasEstimateFuncImpl func(ctx context.Context, chainId *big.Int, msg ethereum.CallMsg) (*big.Int, error)

func (i GasEstimateFuncImpl) EstimateGas(ctx context.Context, chainId *big.Int, msg ethereum.CallMsg) (*big.Int, error) {
	return i(ctx, chainId, msg)
}
