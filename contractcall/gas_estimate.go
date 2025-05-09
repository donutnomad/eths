package contractcall

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type GasEstimateImpl struct {
	client *ethclient.Client
	logger ILogger
}

func NewGasEstimateImpl(client *ethclient.Client, logger ILogger) *GasEstimateImpl {
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
