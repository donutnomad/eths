package etherscan

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/samber/mo"
)

// TransactionClient 提供对 Etherscan Transaction API 的访问
type TransactionClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// NewTransactionClient 创建一个新的 TransactionClient 实例
func NewTransactionClient(apiKey string) *TransactionClient {
	return &TransactionClient{
		client:  resty.New(),
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/v2/api",
	}
}

func NewTransactionClientWith(baseUrl, apiKey string, client *resty.Client) *TransactionClient {
	return &TransactionClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseUrl,
	}
}

func (c *TransactionClient) SetClient(client *resty.Client) {
	c.client = client
}

// ContractExecutionStatusResult 表示合约执行状态结果
type ContractExecutionStatusResult struct {
	IsError        string `json:"isError"`        // "0" 表示成功，"1" 表示失败
	ErrDescription string `json:"errDescription"` // 错误描述（如果有）
}

// ContractExecutionStatusResponse 表示合约执行状态的 API 响应
type ContractExecutionStatusResponse struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Result  ContractExecutionStatusResult `json:"result"`
}

// TransactionReceiptStatusResult 表示交易收据状态结果
type TransactionReceiptStatusResult struct {
	Status string `json:"status"` // "0" 表示失败，"1" 表示成功
}

// TransactionReceiptStatusResponse 表示交易收据状态的 API 响应
type TransactionReceiptStatusResponse struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Result  TransactionReceiptStatusResult `json:"result"`
}

// GetContractExecutionStatusOptions 用于获取合约执行状态的选项
type GetContractExecutionStatusOptions struct {
	TxHash  string         // 必需：交易哈希
	ChainID mo.Option[int] // 可选：链ID，默认为1（以太坊主网）
}

// GetTransactionReceiptStatusOptions 用于获取交易收据状态的选项
type GetTransactionReceiptStatusOptions struct {
	TxHash  string         // 必需：交易哈希
	ChainID mo.Option[int] // 可选：链ID，默认为1（以太坊主网）
}

// IsSuccess 检查合约执行是否成功
func (r ContractExecutionStatusResult) IsSuccess() bool {
	return r.IsError == "0"
}

// HasError 检查合约执行是否失败
func (r ContractExecutionStatusResult) HasError() bool {
	return r.IsError == "1"
}

// IsSuccess 检查交易是否成功
func (r TransactionReceiptStatusResult) IsSuccess() bool {
	return r.Status == "1"
}

// HasError 检查交易是否失败
func (r TransactionReceiptStatusResult) HasError() bool {
	return r.Status == "0"
}

// GetContractExecutionStatus 获取合约执行状态
// 返回合约执行的状态码，用于检查智能合约交易是否成功执行
func (c *TransactionClient) GetContractExecutionStatus(opts GetContractExecutionStatusOptions) (*ContractExecutionStatusResponse, error) {
	if opts.TxHash == "" {
		return nil, fmt.Errorf("txhash parameter is required")
	}

	queryParams := map[string]string{
		"module": "transaction",
		"action": "getstatus",
		"txhash": opts.TxHash,
		"apikey": c.apiKey,
	}

	if opts.ChainID.IsPresent() {
		queryParams["chainid"] = fmt.Sprintf("%d", opts.ChainID.MustGet())
	} else {
		queryParams["chainid"] = "1" // 默认为以太坊主网
	}

	var response ContractExecutionStatusResponse
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if response.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	return &response, nil
}

// GetTransactionReceiptStatus 获取交易收据状态
// 返回交易执行的状态码，仅适用于拜占庭分叉后的交易
func (c *TransactionClient) GetTransactionReceiptStatus(opts GetTransactionReceiptStatusOptions) (*TransactionReceiptStatusResponse, error) {
	if opts.TxHash == "" {
		return nil, fmt.Errorf("txhash parameter is required")
	}

	queryParams := map[string]string{
		"module": "transaction",
		"action": "gettxreceiptstatus",
		"txhash": opts.TxHash,
		"apikey": c.apiKey,
	}

	if opts.ChainID.IsPresent() {
		queryParams["chainid"] = fmt.Sprintf("%d", opts.ChainID.MustGet())
	} else {
		queryParams["chainid"] = "1" // 默认为以太坊主网
	}

	var response TransactionReceiptStatusResponse
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if response.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	return &response, nil
}
