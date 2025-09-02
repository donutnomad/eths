package etherscan

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTxStatus(t *testing.T) {
	ctx := context.Background()
	status, err := etherscanClient.GetTransactionReceiptStatus(
		ctx, common.HexToHash("0xc7b2959514f65cffc16066254e0744c944a1056bd2d9956416349fa1021300a1"), 1,
	)
	require.Nil(t, err)
	require.True(t, status.Result.IsSuccess())
}

func TestTxStatus2(t *testing.T) {
	//(*etherscan.Response[github.com/donutnomad/eths/etherscan.ContractExecutionStatusResult])(0xc00019c040)({
	// Status: (string) (len=1) "1",
	// Message: (string) (len=2) "OK",
	// Result: (etherscan.ContractExecutionStatusResult) {
	//  IsError: (string) (len=1) "0",
	//  ErrDescription: (string) ""
	// }
	//})
	ctx := context.Background()
	status, err := etherscanClient.GetContractExecutionStatus(
		ctx, common.HexToHash("0xc7b2959514f65cffc16066254e0744c944a1056bd2d9956416349fa1021300a1"), 1,
	)
	require.Nil(t, err)
	require.True(t, status.Result.IsSuccess())
}
