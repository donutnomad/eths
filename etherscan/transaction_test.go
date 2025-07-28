package etherscan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 模拟的合约执行状态响应数据
var mockContractExecutionStatusResponse = ContractExecutionStatusResponse{
	Status:  "1",
	Message: "OK",
	Result: ContractExecutionStatusResult{
		IsError:        "1",
		ErrDescription: "Bad jump destination",
	},
}

var mockContractExecutionStatusSuccessResponse = ContractExecutionStatusResponse{
	Status:  "1",
	Message: "OK",
	Result: ContractExecutionStatusResult{
		IsError:        "0",
		ErrDescription: "",
	},
}

// 模拟的交易收据状态响应数据
var mockTransactionReceiptStatusResponse = TransactionReceiptStatusResponse{
	Status:  "1",
	Message: "OK",
	Result: TransactionReceiptStatusResult{
		Status: "1",
	},
}

var mockTransactionReceiptStatusFailResponse = TransactionReceiptStatusResponse{
	Status:  "1",
	Message: "OK",
	Result: TransactionReceiptStatusResult{
		Status: "0",
	},
}

var mockTransactionReceiptErrorResponse = TransactionReceiptStatusResponse{
	Status:  "0",
	Message: "NOTOK",
	Result: TransactionReceiptStatusResult{
		Status: "",
	},
}

var mockContractExecutionErrorResponse = ContractExecutionStatusResponse{
	Status:  "0",
	Message: "NOTOK",
	Result: ContractExecutionStatusResult{
		IsError:        "",
		ErrDescription: "",
	},
}

// createTransactionMockServer 创建交易API的模拟服务器
func createTransactionMockServer(responseData interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(responseData)
	}))
}

// createTransactionClientWithMockServer 创建使用模拟服务器的交易客户端
func createTransactionClientWithMockServer(server *httptest.Server) *TransactionClient {
	client := NewEtherscanClient("test-api-key")
	client.baseURL = server.URL
	return client
}

func TestNewTransactionClient(t *testing.T) {
	client := NewEtherscanClient("test-api-key")

	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.apiKey)
	assert.Equal(t, "https://api.etherscan.io/v2/api", client.baseURL)
	assert.NotNil(t, client.client)
}

func TestTransactionClient_GetContractExecutionStatus(t *testing.T) {
	tests := []struct {
		name           string
		options        GetContractExecutionStatusOptions
		mockResponse   interface{}
		mockStatusCode int
		expectedError  string
		expectedResult *ContractExecutionStatusResult
	}{
		{
			name: "成功获取失败的合约执行状态",
			options: GetContractExecutionStatusOptions{
				TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
			},
			mockResponse:   mockContractExecutionStatusResponse,
			mockStatusCode: 200,
			expectedError:  "",
			expectedResult: &ContractExecutionStatusResult{
				IsError:        "1",
				ErrDescription: "Bad jump destination",
			},
		},
		{
			name: "成功获取成功的合约执行状态",
			options: GetContractExecutionStatusOptions{
				TxHash:  "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
				ChainID: mo.Some(1),
			},
			mockResponse:   mockContractExecutionStatusSuccessResponse,
			mockStatusCode: 200,
			expectedError:  "",
			expectedResult: &ContractExecutionStatusResult{
				IsError:        "0",
				ErrDescription: "",
			},
		},
		{
			name: "交易哈希参数为空",
			options: GetContractExecutionStatusOptions{
				TxHash: "",
			},
			mockResponse:   nil,
			mockStatusCode: 200,
			expectedError:  "txhash parameter is required",
			expectedResult: nil,
		},
		{
			name: "API 返回错误状态",
			options: GetContractExecutionStatusOptions{
				TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
			},
			mockResponse:   mockContractExecutionErrorResponse,
			mockStatusCode: 200,
			expectedError:  "API returned error: NOTOK",
			expectedResult: nil,
		},
		{
			name: "HTTP 状态码错误",
			options: GetContractExecutionStatusOptions{
				TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
			},
			mockResponse:   mockContractExecutionStatusResponse,
			mockStatusCode: 500,
			expectedError:  "API request failed with status code: 500",
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "txhash parameter is required" {
				client := NewEtherscanClient("test-api-key")
				result, err := client.GetContractExecutionStatus(tt.options)

				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			server := createTransactionMockServer(tt.mockResponse, tt.mockStatusCode)
			defer server.Close()

			client := createTransactionClientWithMockServer(server)
			result, err := client.GetContractExecutionStatus(tt.options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Status)
				assert.Equal(t, "OK", result.Message)
				assert.Equal(t, *tt.expectedResult, result.Result)
			}
		})
	}
}

func TestTransactionClient_GetTransactionReceiptStatus(t *testing.T) {
	tests := []struct {
		name           string
		options        GetTransactionReceiptStatusOptions
		mockResponse   interface{}
		mockStatusCode int
		expectedError  string
		expectedResult *TransactionReceiptStatusResult
	}{
		{
			name: "成功获取成功的交易收据状态",
			options: GetTransactionReceiptStatusOptions{
				TxHash: "0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76",
			},
			mockResponse:   mockTransactionReceiptStatusResponse,
			mockStatusCode: 200,
			expectedError:  "",
			expectedResult: &TransactionReceiptStatusResult{
				Status: "1",
			},
		},
		{
			name: "成功获取失败的交易收据状态",
			options: GetTransactionReceiptStatusOptions{
				TxHash:  "0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76",
				ChainID: mo.Some(1),
			},
			mockResponse:   mockTransactionReceiptStatusFailResponse,
			mockStatusCode: 200,
			expectedError:  "",
			expectedResult: &TransactionReceiptStatusResult{
				Status: "0",
			},
		},
		{
			name: "交易哈希参数为空",
			options: GetTransactionReceiptStatusOptions{
				TxHash: "",
			},
			mockResponse:   nil,
			mockStatusCode: 200,
			expectedError:  "txhash parameter is required",
			expectedResult: nil,
		},
		{
			name: "API 返回错误状态",
			options: GetTransactionReceiptStatusOptions{
				TxHash: "0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76",
			},
			mockResponse:   mockTransactionReceiptErrorResponse,
			mockStatusCode: 200,
			expectedError:  "API returned error: NOTOK",
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "txhash parameter is required" {
				client := NewEtherscanClient("test-api-key")
				result, err := client.GetTransactionReceiptStatus(tt.options)

				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			server := createTransactionMockServer(tt.mockResponse, tt.mockStatusCode)
			defer server.Close()

			client := createTransactionClientWithMockServer(server)
			result, err := client.GetTransactionReceiptStatus(tt.options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Status)
				assert.Equal(t, "OK", result.Message)
				assert.Equal(t, *tt.expectedResult, result.Result)
			}
		})
	}
}

// 测试结果类型的便利方法
func TestContractExecutionStatusResult_HelperMethods(t *testing.T) {
	// 测试成功的情况
	successResult := ContractExecutionStatusResult{
		IsError:        "0",
		ErrDescription: "",
	}
	assert.True(t, successResult.IsSuccess())
	assert.False(t, successResult.HasError())

	// 测试失败的情况
	failResult := ContractExecutionStatusResult{
		IsError:        "1",
		ErrDescription: "Bad jump destination",
	}
	assert.False(t, failResult.IsSuccess())
	assert.True(t, failResult.HasError())
}

func TestTransactionReceiptStatusResult_HelperMethods(t *testing.T) {
	// 测试成功的情况
	successResult := TransactionReceiptStatusResult{
		Status: "1",
	}
	assert.True(t, successResult.IsSuccess())
	assert.False(t, successResult.HasError())

	// 测试失败的情况
	failResult := TransactionReceiptStatusResult{
		Status: "0",
	}
	assert.False(t, failResult.IsSuccess())
	assert.True(t, failResult.HasError())
}

// 测试请求参数验证
func TestTransactionAPI_RequestParameterValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// 验证必需的参数存在
		assert.NotEmpty(t, query.Get("module"))
		assert.NotEmpty(t, query.Get("action"))
		assert.NotEmpty(t, query.Get("txhash"))
		assert.Equal(t, "test-api-key", query.Get("apikey"))

		// 验证 chainid 参数
		chainid := query.Get("chainid")
		assert.NotEmpty(t, chainid)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockContractExecutionStatusResponse)
	}))
	defer server.Close()

	client := createTransactionClientWithMockServer(server)

	// 测试合约执行状态请求
	_, err := client.GetContractExecutionStatus(GetContractExecutionStatusOptions{
		TxHash:  "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
		ChainID: mo.Some(1),
	})

	assert.NoError(t, err)
}

// 基准测试
func BenchmarkTransactionClient_GetContractExecutionStatus(b *testing.B) {
	server := createTransactionMockServer(mockContractExecutionStatusResponse, 200)
	defer server.Close()

	client := createTransactionClientWithMockServer(server)
	options := GetContractExecutionStatusOptions{
		TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetContractExecutionStatus(options)
		require.NoError(b, err)
	}
}

// 集成测试（需要真实的 API 密钥）
func TestTransactionClient_Integration(t *testing.T) {
	t.Skip("跳过集成测试，需要真实的 API 密钥")

	apiKey := "YOUR_REAL_API_KEY" // 替换为真实的 API 密钥
	client := NewEtherscanClient(apiKey)

	// 测试合约执行状态
	result1, err := client.GetContractExecutionStatus(GetContractExecutionStatusOptions{
		TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result1)
	assert.Equal(t, "1", result1.Status)

	// 测试交易收据状态
	result2, err := client.GetTransactionReceiptStatus(GetTransactionReceiptStatusOptions{
		TxHash: "0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, "1", result2.Status)
}

// 测试并发安全性
func TestTransactionClient_ConcurrentRequests(t *testing.T) {
	server := createTransactionMockServer(mockContractExecutionStatusResponse, 200)
	defer server.Close()

	client := createTransactionClientWithMockServer(server)

	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			_, err := client.GetContractExecutionStatus(GetContractExecutionStatusOptions{
				TxHash: "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
			})
			results <- err
		}()
	}

	// 检查所有请求都成功
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
