package contractcall

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type EthClientEnhance struct {
	client  *ethclient.Client
	ctx     context.Context
	timeout time.Duration
}

func NewEthClientEnhance(client *ethclient.Client, ctx context.Context, timeout time.Duration) *EthClientEnhance {
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

func (e *EthClientEnhance) CallContract(_ context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(e.ctx, e.timeout)
	defer cancel()
	return e.client.CallContract(ctx, call, blockNumber)
}

func SendTxAndWait(ctx context.Context, client *ethclient.Client, chainId *big.Int, contract common.Address, payer ISigner, data []byte, blockConfirmations uint64, callManager *CallManager, beforeSend func(tx *types.Transaction) error) error {
	tx, err := SendTx(ctx, client, chainId, data, contract, payer, callManager, beforeSend)
	if err != nil {
		return err
	}
	return Wait(ctx, client, tx.Hash(), blockConfirmations, nil)
}

func Wait(ctx context.Context, client *ethclient.Client, txHash common.Hash, blockConfirmations uint64, outReceipt *types.Receipt) error {
	return WaitRetry(ctx, client, txHash, 10, blockConfirmations, outReceipt)
}

func WaitRetry(ctx context.Context, client *ethclient.Client, txHash common.Hash, retryTimes int, blockConfirmations uint64, outReceipt *types.Receipt) error {
	var receipt *types.Receipt
	// wait receipt
	for i := 0; i < retryTimes; i++ {
		var err error
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				time.Sleep(3 * time.Second)
				continue
			} else {
				return errors.Wrap(EthereumRPCErr, err.Error())
			}
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
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
