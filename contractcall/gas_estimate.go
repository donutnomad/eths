package contractcall

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"math/big"
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
