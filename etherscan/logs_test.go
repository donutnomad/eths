package etherscan

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
)

var etherscanApiKey string
var etherscanClient *EtherscanClient

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("未找到 .env 文件，或加载失败。尝试从环境变量获取。")
	}
	apiKey := os.Getenv("ETHERSCAN_APIKEY")
	if apiKey == "" {
		panic("ETHERSCAN_APIKEY 环境变量未设置")
	}
	etherscanApiKey = apiKey
	client := NewEtherscanClient(etherscanApiKey)
	//client.SetDebug()
	etherscanClient = client
}

func TestGetLogs(t *testing.T) {
	chainID := uint64(1)
	tokenAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7") // mainnet USDT
	ctx := context.Background()
	logs, err := etherscanClient.GetLogs(ctx, chainID, &tokenAddress, nil, nil, 1, 10, GetLogsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	require.Nil(t, err)
	require.Greater(t, len(logs.Result), 0)
}

func TestGetLogs2(t *testing.T) {
	chainID := uint64(1)
	ctx := context.Background()
	logs, err := etherscanClient.GetLogs(ctx, chainID, nil, nil, nil, 1, 10, GetLogsOptions{
		Topic0:     mo.Some(common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")),
		Topic1:     mo.Some(common.HexToHash("0x00000000000000000000000036928500bc1dcd7af6a2b4008875cc336b927d51")),
		Topic01Opr: mo.Some(true),
	})
	require.Nil(t, err)
	require.Len(t, logs.Result, 0)
}

func TestGetLogs3(t *testing.T) {
	chainID := uint64(1)
	ctx := context.Background()
	logs, err := etherscanClient.GetLogs(ctx, chainID, nil, nil, nil, 1, 10, GetLogsOptions{
		Topic0:     mo.Some(common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")),
		Topic1:     mo.Some(common.HexToHash("0x00000000000000000000000036928500bc1dcd7af6a2b4008875cc336b927d57")),
		Topic01Opr: mo.Some(true),
	})
	require.Nil(t, err)
	require.Greater(t, len(logs.Result), 0)
}
