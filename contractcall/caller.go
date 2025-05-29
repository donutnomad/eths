package contractcall

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type CallManager struct {
	GasValidator IGasPriceValidator
	GasPricer    IGasPricer
	GasEstimate  IEstimateGas
	NonceManager IGetNonce
}

func NewDefaultCallManager(client IMyClient, logger ILogger) *CallManager {
	return &CallManager{
		GasValidator: NewDefaultGasValidator(client),
		GasPricer:    NewGasPricerDefault(client),
		GasEstimate:  NewGasEstimateImpl(client, logger),
		NonceManager: NewDefaultNonceManager(client),
	}
}

func SendTx(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	data []byte,
	to common.Address,
	payer ISigner,
	gasManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
) (*ethTypes.Transaction, error) {
	return _send(ctx, client, chainId, nil, data, to, payer, gasManager, beforeSend, false, true)
}

func NoSendTx(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	data []byte,
	to common.Address,
	payer ISigner,
	gasManager *CallManager,
) (*ethTypes.Transaction, error) {
	return _send(ctx, client, chainId, nil, data, to, payer, gasManager, nil, true, true)
}

func SendTxE(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	to common.Address,
	payer ISigner,
	gasManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
	noSend bool,
	toIsContract bool,
) (*ethTypes.Transaction, error) {
	return _send(ctx, client, chainId, value, data, to, payer, gasManager, beforeSend, noSend, toIsContract)
}

func EstimateTxE(
	ctx context.Context,
	client ICodeAt,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	from common.Address,
	to common.Address,
	gasManager *CallManager,
	toIsContract bool,
) (*GasPrice, *big.Int, error) {
	txBuilder := NewTxBuilder(ctx, chainId).
		SetFrom(from).
		SetTo(to, toIsContract).
		SetValue(value).
		SetData(data).
		SetNonceBy(gasManager.NonceManager).
		SetGasPriceBy(gasManager.GasPricer).
		SetGasLimitBy(gasManager.GasEstimate).
		Check(client, nil)
	return txBuilder.gasPrice, txBuilder.gasLimit, txBuilder.err
}

func SendTxBuilder(
	ctx context.Context,
	txBuilder *TxBuilder,
	client ethereum.TransactionSender,
	payer ISigner,
	noSend bool,
	beforeSend func(tx *ethTypes.Transaction) error,
) (*ethTypes.Transaction, error) {
	txWrapper, err := txBuilder.SetFromByKey(payer.PublicKey()).Build()
	if err != nil {
		return nil, err
	}
	txWrapper, err = txWrapper.Sign(payer)
	if err != nil {
		return nil, err
	}
	tx := txWrapper.ToTransaction()
	if beforeSend != nil {
		err = beforeSend(tx)
		if err != nil {
			return nil, err
		}
	}
	if noSend {
		return tx, nil
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

func _send(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	to common.Address,
	payer ISigner,
	gasManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
	noSend bool,
	checkContract bool,
) (*ethTypes.Transaction, error) {
	txBuilder := NewTxBuilder(ctx, chainId).
		SetFromByKey(payer.PublicKey()).
		SetTo(to, checkContract).
		SetValue(value).
		SetData(data).
		SetNonceBy(gasManager.NonceManager).
		SetGasPriceBy(gasManager.GasPricer).
		SetGasLimitBy(gasManager.GasEstimate).
		Check(client, gasManager.GasValidator)
	return SendTxBuilder(ctx, txBuilder, client, payer, noSend, beforeSend)
}
