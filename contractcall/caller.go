package contractcall

import (
	"context"
	"fmt"
	"math/big"

	"github.com/donutnomad/blockchain-alg/xecdsa"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CallManager struct {
	GasValidator IGasPriceValidator
	GasPricer    IGasPricer
	GasEstimate  IEstimateGas
	NonceManager IGetNonce
}

func NewDefaultCallManager(client *ethclient.Client, logger ILogger) *CallManager {
	return &CallManager{
		GasValidator: NewDefaultGasValidator(client),
		GasPricer:    NewGasPricerDefault(client),
		GasEstimate:  NewGasEstimateImpl(client, logger),
		NonceManager: NewDefaultNonceManager(client),
	}
}

func SendTx(
	ctx context.Context,
	client *ethclient.Client,
	chainId *big.Int,
	data []byte,
	to common.Address,
	payer *xecdsa.PrivateKey,
	gasManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
) (*ethTypes.Transaction, error) {
	txBuilder := NewTxBuilder(ctx, chainId).
		SetFromByKey(payer.PublicKey).
		SetTo(to, true).
		SetData(data).
		SetNonceBy(gasManager.NonceManager).
		SetGasPriceBy(gasManager.GasPricer).
		SetGasLimitBy(gasManager.GasEstimate).
		Check(client, gasManager.GasValidator)
	return SendTxBuilder(ctx, txBuilder, client, payer, beforeSend)
}

func SendTxBuilder(
	ctx context.Context,
	txBuilder *TxBuilder,
	client *ethclient.Client,
	payer *xecdsa.PrivateKey,
	beforeSend func(tx *ethTypes.Transaction) error,
) (*ethTypes.Transaction, error) {
	txWrapper, err := txBuilder.SetFromByKey(payer.PublicKey).Build()
	if err != nil {
		return nil, err
	}
	tx := txWrapper.Sign(payer).ToTransaction()
	if beforeSend != nil {
		err = beforeSend(tx)
		if err != nil {
			return nil, err
		}
	}
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		errSend := &SendTransactionError{
			Tx:  tx,
			Err: err,
		}
		return nil, fmt.Errorf("ethereum send transaction failed: %w,%w", errSend, EthereumRPCErr)
	}
	return tx, nil
}
