package etherscan

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// GetContractExecutionStatus 获取合约执行状态
// 返回合约执行的状态码，用于检查智能合约交易是否成功执行
func (e *EtherscanClient) GetContractExecutionStatus(ctx context.Context, txHash common.Hash, chainID uint64) (*Response[ContractExecutionStatusResult], error) {
	return getResult[ContractExecutionStatusResult](
		ctx,
		e,
		map[string]string{
			"txhash": txHash.Hex(),
		},
		"transaction",
		"getstatus",
		chainID,
	)
}

// GetTransactionReceiptStatus 获取交易收据状态
// 返回交易执行的状态码，仅适用于拜占庭分叉后的交易
func (e *EtherscanClient) GetTransactionReceiptStatus(ctx context.Context, txHash common.Hash, chainID uint64) (*Response[TransactionReceiptStatusResult], error) {
	return getResult[TransactionReceiptStatusResult](
		ctx,
		e,
		map[string]string{
			"txhash": txHash.Hex(),
		},
		"transaction",
		"gettxreceiptstatus",
		chainID,
	)
}

type ContractExecutionStatusResult struct {
	IsError        string `json:"isError"`        // "0" 表示成功，"1" 表示失败
	ErrDescription string `json:"errDescription"` // 错误描述（如果有）
}

func (r ContractExecutionStatusResult) IsSuccess() bool {
	return r.IsError == "0"
}

type TransactionReceiptStatusResult struct {
	Status string `json:"status"` // "0" 表示失败，"1" 表示成功
}

func (r TransactionReceiptStatusResult) IsSuccess() bool {
	return r.Status == "1"
}
