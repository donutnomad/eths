package contractcall

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type CallManager struct {
	GasValidator IGasPriceValidator
	GasPricer    IGasPricer
	GasEstimate  IEstimateGas
	NonceManager IGetNonce
	Balance      IBalanceChecker
}

func NewDefaultCallManager(client IMyClient, logger ILogger) *CallManager {
	return &CallManager{
		GasValidator: NewDefaultGasValidator(client),
		GasPricer:    NewGasPricerDefault(client),
		GasEstimate:  NewGasEstimateImpl(client, logger),
		NonceManager: NewDefaultNonceManager(client),
		Balance:      NewBalanceCheckerImpl(client),
	}
}

type SendTxOption struct {
	Strict bool // true: 不允许向合约发送空data（因为，如果合约没有设置payable方法，则不能接受ETH，返回的只是execution reverted, 让人无法理解
}

func SendTx(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	data []byte,
	to common.Address,
	payer ISigner,
	callManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
	opt ...SendTxOption,
) (*ethTypes.Transaction, error) {
	return send(ctx, client, chainId, nil, data, &to, payer, callManager, beforeSend, false, true, opt...)
}

func SendTxE(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	to *common.Address,
	payer ISigner,
	callManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
	noSend bool,
	toIsContract bool,
	opt ...SendTxOption,
) (*ethTypes.Transaction, error) {
	return send(ctx, client, chainId, value, data, to, payer, callManager, beforeSend, noSend, toIsContract, opt...)
}

func EstimateTxE(
	ctx context.Context,
	client ICodeAt,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	from common.Address,
	to *common.Address,
	callManager *CallManager,
	toIsContract bool,
	opt ...SendTxOption,
) (*GasPrice, *big.Int, error) {
	builder, err := sendFn[*TxBuilder](ctx, &noOpTransactionSender{client}, chainId, value, data, to, NewNoOpSigner(from, nil), callManager, toIsContract, func(txBuilder *TxBuilder) (*TxBuilder, error) {
		return txBuilder, nil
	}, opt...)
	if err != nil {
		return nil, nil, err
	}
	return builder.gasPrice, builder.gasLimit, builder.err
}

func SendTxBuilder(
	ctx context.Context,
	txBuilder *TxBuilder,
	client ethereum.TransactionSender,
	payer ISigner,
	noSend bool,
	beforeSend func(tx *ethTypes.Transaction) error,
) (*ethTypes.Transaction, error) {
	txWrapper, err := txBuilder.SetFrom(payer.Address()).Build()
	if err != nil {
		return nil, err
	}
	if err = txWrapper.Sign(payer); err != nil {
		return nil, err
	}
	tx := txWrapper.ToTransaction()
	if beforeSend != nil {
		if err = beforeSend(tx); err != nil {
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
			Err: ParseEvmError(err),
		}
		return nil, fmt.Errorf("ethereum send transaction failed: %w,%w", errSend, EthereumRPCErr)
	}
	return tx, nil
}

func send(
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	to *common.Address,
	payer ISigner,
	callManager *CallManager,
	beforeSend func(tx *ethTypes.Transaction) error,
	noSend bool,
	checkContract bool,
	opt ...SendTxOption,
) (*ethTypes.Transaction, error) {
	return sendFn[*ethTypes.Transaction](ctx, client, chainId, value, data, to, payer, callManager, checkContract, func(txBuilder *TxBuilder) (*ethTypes.Transaction, error) {
		return SendTxBuilder(ctx, txBuilder, client, payer, noSend, beforeSend)
	}, opt...)
}

func sendFn[T any](
	ctx context.Context,
	client ISendTxClient,
	chainId *big.Int,
	value *big.Int,
	data []byte,
	to *common.Address,
	payer ISigner,
	callManager *CallManager,
	checkContract bool,
	fn func(txBuilder *TxBuilder) (T, error),
	opt ...SendTxOption,
) (T, error) {
	txBuilder := NewTxBuilder(ctx, chainId)
	if to != nil {
		txBuilder = txBuilder.SetTo(*to, checkContract)
	}
	txBuilder.
		SetFrom(payer.Address()).
		SetValue(value).
		SetData(data).
		SetNonceBy(callManager.NonceManager).
		SetGasPriceBy(callManager.GasPricer).
		BalanceCheck(callManager.Balance).
		SetGasLimitBy(callManager.GasEstimate).
		Check(client, callManager.GasValidator, opt...)
	if txBuilder.err != nil {
		return *new(T), txBuilder.err
	}
	return fn(txBuilder)
}

type noOpTransactionSender struct {
	ICodeAt
}

func (e *noOpTransactionSender) SendTransaction(ctx context.Context, tx *ethTypes.Transaction) error {
	return nil
}
