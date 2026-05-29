package contractcall

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethclient"
	"github.com/donutnomad/eths/ethtype"
	"github.com/pkg/errors"
)

type EthClientEnhance struct {
	client  IEthereumRPC
	ctx     context.Context
	timeout time.Duration
}

func NewEthClientEnhance(client IEthereumRPC, ctx context.Context, timeout time.Duration) *EthClientEnhance {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &EthClientEnhance{
		client:  client,
		ctx:     ctx,
		timeout: timeout,
	}
}

func (e *EthClientEnhance) PendingCodeAt(_ context.Context, account common.Address) ([]byte, error) {
	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()
	return e.client.PendingCodeAt(ctx, account)
}

func (e *EthClientEnhance) ChainID(_ context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()
	return e.client.ChainID(ctx)
}

func (e *EthClientEnhance) CodeAt(_ context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()
	return e.client.CodeAt(ctx, contract, blockNumber)
}

func (e *EthClientEnhance) CallContract(_ context.Context, call ethclient.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()
	return e.client.CallContract(ctx, call, blockNumber)
}

func SendTxAndWait[C interface {
	ethclient.TransactionReader
	ethclient.BlockNumberReader
	ethclient.TransactionSender
	ICodeAt
}](ctx context.Context, client C, chainId *big.Int, contract common.Address, payer ISigner, data []byte, blockConfirmations uint64, callManager *CallManager, beforeSend func(tx *ethtype.ETransaction) error) error {
	tx, err := SendTx(ctx, client, chainId, data, contract, payer, callManager, beforeSend)
	if err != nil {
		return err
	}
	return Wait(ctx, client, tx.Hash(), blockConfirmations, nil)
}

func Wait[C interface {
	ethclient.TransactionReader
	ethclient.BlockNumberReader
}](ctx context.Context, client C, txHash common.Hash, blockConfirmations uint64, outReceipt *ethtype.EReceipt) error {
	return WaitRetry(ctx, client, txHash, 10, blockConfirmations, outReceipt)
}

func WaitRetry[C interface {
	ethclient.TransactionReader
	ethclient.BlockNumberReader
}](ctx context.Context, client C, txHash common.Hash, retryTimes int, blockConfirmations uint64, outReceipt *ethtype.EReceipt) error {
	var receipt *ethtype.EReceipt
	// wait receipt
	for i := 0; i < retryTimes; i++ {
		var err error
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if errors.Is(err, ethclient.NotFound) {
				time.Sleep(3 * time.Second)
				continue
			} else {
				return errors.Wrap(EthereumRPCErr, err.Error())
			}
		}
		if receipt.Status != ethtype.ReceiptStatusSuccessful {
			return fmt.Errorf("transaction %s receipt status is not successful", txHash.Hex())
		}
		if outReceipt != nil {
			*outReceipt = *receipt
		}
	}

	if receipt == nil || blockConfirmations == 0 {
		return nil
	}

	// wait for block confirmations
	blockNumber := receipt.BlockNumber
	if blockNumber == nil {
		return fmt.Errorf("transaction %s block number is nil", txHash.Hex())
	}

	// Get current block number
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(EthereumRPCErr, err.Error())
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Wait until we have confirmations
	for currentBlock-blockNumber.Uint64() < blockConfirmations {
		// Wait for 3 seconds before checking again
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
		// Get new current block number
		currentBlock, err = client.BlockNumber(ctx)
		if err != nil {
			return errors.Wrap(EthereumRPCErr, err.Error())
		}
	}
	return nil
}
