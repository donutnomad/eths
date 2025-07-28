package etherscan

import "github.com/go-resty/resty/v2"

// EtherscanClient 提供对 Etherscan Transaction API 的访问
type EtherscanClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// NewEtherscanClient 创建一个新的 TransactionClient 实例
func NewEtherscanClient(apiKey string) *EtherscanClient {
	return &EtherscanClient{
		client:  resty.New(),
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/v2/api",
	}
}

func NewEtherscanClientWith(baseUrl, apiKey string, client *resty.Client) *EtherscanClient {
	return &EtherscanClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseUrl,
	}
}

func (c *EtherscanClient) SetClient(client *resty.Client) {
	c.client = client
}
