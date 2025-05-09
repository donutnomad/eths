package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// DefaultGasValidator implements IGasPriceValidator
type DefaultGasValidator struct {
	client *ethclient.Client
}

func NewDefaultGasValidator(client *ethclient.Client) *DefaultGasValidator {
	return &DefaultGasValidator{client: client}
}

// ValidateGasPrice validates if the gas price is reasonable
func (v *DefaultGasValidator) ValidateGasPrice(ctx context.Context, chainId *big.Int, gasPrice *GasPrice) error {
	// Get the latest block header to check base fee
	header, err := v.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return errors.Wrap(EthereumRPCErr, err.Error())
	}
	// baseFee * 1.2
	minMaxFee := decimal.NewFromBigInt(header.BaseFee, 0).Mul(decimal.NewFromFloat(1.2)).BigInt()

	// If using legacy gas
	if gasPrice.LegacyGas != nil {
		if header.BaseFee != nil {
			// For EIP-1559 chains, legacy gas price should be at least base fee
			if gasPrice.LegacyGas.GasPrice.Cmp(minMaxFee) < 0 {
				return errors.Wrap(GasInvalidGasPriceErr, "legacy gas price is lower than base fee")
			}
		}
		return nil
	}

	// If using dynamic gas (EIP-1559)
	if gasPrice.DynamicGas != nil {
		if header.BaseFee == nil {
			return errors.Wrap(GasInvalidGasPriceErr, "chain does not support EIP-1559")
		}
		// Max fee should be at least base fee + priority fee
		if gasPrice.DynamicGas.MaxFeePerGas.Cmp(minMaxFee) < 0 {
			return errors.Wrap(GasInvalidGasPriceErr, "max fee per gas is too low")
		}
	}

	return nil
}
