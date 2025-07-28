package etherscan

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 模拟的 API 响应数据
var mockLogResponse = LogsResponse{
	Status:  "1",
	Message: "OK",
	Result: []LogEntry{
		func() LogEntry {
			topic0, _ := NewHashFromHex("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
			topic1, _ := NewHashFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
			txHash, _ := NewHashFromHex("0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9")
			return LogEntry{
				Address:          "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
				Topics:           []Hash{topic0, topic1},
				Data:             []byte{},     // 空的字节数组，对应 "0x"
				BlockNumber:      12878196,     // 解析后的 uint64 值
				TimeStamp:        1626984022,   // 解析后的 uint64 值
				GasPrice:         200000000000, // 解析后的 uint64 值
				GasUsed:          2388485,      // 解析后的 uint64 值
				LogIndex:         0,            // 解析后的 uint64 值
				TransactionHash:  txHash,
				TransactionIndex: 0, // 解析后的 uint64 值
			}
		}(),
	},
}

var mockErrorResponse = LogsResponse{
	Status:  "0",
	Message: "NOTOK",
	Result:  []LogEntry{},
}

// createMockServer 创建模拟服务器
func createMockServer(responseData interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		// 如果是 LogsResponse，返回原始的 JSON 格式
		if logsResp, ok := responseData.(LogsResponse); ok {
			// 构建原始 JSON 响应
			rawResponse := map[string]interface{}{
				"status":  logsResp.Status,
				"message": logsResp.Message,
				"result":  []map[string]interface{}{},
			}

			// 转换每个 LogEntry 为原始格式
			for _, entry := range logsResp.Result {
				rawEntry := map[string]interface{}{
					"address":          entry.Address,
					"topics":           []string{entry.Topics[0].String(), entry.Topics[1].String()},
					"data":             entry.Data,
					"blockNumber":      "0xc48174",
					"timeStamp":        "0x60f9ce56",
					"gasPrice":         "0x2e90edd000",
					"gasUsed":          "0x247205",
					"logIndex":         "0x",
					"transactionHash":  entry.TransactionHash.String(),
					"transactionIndex": "0x",
				}
				rawResponse["result"] = append(rawResponse["result"].([]map[string]interface{}), rawEntry)
			}

			json.NewEncoder(w).Encode(rawResponse)
		} else {
			json.NewEncoder(w).Encode(responseData)
		}
	}))
}

// createClientWithMockServer 创建使用模拟服务器的客户端
func createClientWithMockServer(server *httptest.Server) *EtherscanClient {
	client := NewEtherscanClient("test-api-key")
	client.baseURL = server.URL
	return client
}

func TestNewLogsClient(t *testing.T) {
	client := NewEtherscanClient("test-api-key")

	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.apiKey)
	assert.Equal(t, "https://api.etherscan.io/v2/api", client.baseURL)
	assert.NotNil(t, client.client)
}

func TestLogsClient_GetLogsByAddress(t *testing.T) {
	tests := []struct {
		name           string
		options        GetLogsByAddressOptions
		mockResponse   interface{}
		mockStatusCode int
		expectedError  string
	}{
		{
			name: "成功获取日志 - 仅地址参数",
			options: GetLogsByAddressOptions{
				Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "成功获取日志 - 包含所有可选参数",
			options: GetLogsByAddressOptions{
				Address:   "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
				FromBlock: mo.Some(uint64(12878196)),
				ToBlock:   mo.Some(uint64(12878196)),
				Page:      mo.Some(1),
				Offset:    mo.Some(1000),
				ChainID:   mo.Some(uint64(1)),
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "地址参数为空",
			options: GetLogsByAddressOptions{
				Address: "",
			},
			mockResponse:   nil,
			mockStatusCode: 200,
			expectedError:  "address parameter is required",
		},
		{
			name: "API 返回错误状态",
			options: GetLogsByAddressOptions{
				Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
			},
			mockResponse:   mockErrorResponse,
			mockStatusCode: 200,
			expectedError:  "API returned error: NOTOK",
		},
		{
			name: "HTTP 状态码错误",
			options: GetLogsByAddressOptions{
				Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 500,
			expectedError:  "API request failed with status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "address parameter is required" {
				// 对于参数验证错误，不需要启动服务器
				client := NewEtherscanClient("test-api-key")
				result, err := client.GetLogsByAddress(tt.options)

				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			server := createMockServer(tt.mockResponse, tt.mockStatusCode)
			defer server.Close()

			client := createClientWithMockServer(server)
			result, err := client.GetLogsByAddress(tt.options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Status)
				assert.Equal(t, "OK", result.Message)
				assert.Len(t, result.Result, 1)
			}
		})
	}
}

func TestLogsClient_GetLogsByTopics(t *testing.T) {
	tests := []struct {
		name           string
		options        GetLogsByTopicsOptions
		mockResponse   interface{}
		mockStatusCode int
		expectedError  string
	}{
		{
			name: "成功获取日志 - 基本参数",
			options: GetLogsByTopicsOptions{
				FromBlock: 12878196,
				ToBlock:   12879196,
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "成功获取日志 - 包含主题和操作符",
			options: GetLogsByTopicsOptions{
				FromBlock:  12878196,
				ToBlock:    12879196,
				Topic0:     mo.Some("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				Topic1:     mo.Some("0x0000000000000000000000000000000000000000000000000000000000000000"),
				Topic01Opr: mo.Some("and"),
				Page:       mo.Some(1),
				Offset:     mo.Some(1000),
				ChainID:    mo.Some(1),
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "API 返回错误状态",
			options: GetLogsByTopicsOptions{
				FromBlock: 12878196,
				ToBlock:   12879196,
			},
			mockResponse:   mockErrorResponse,
			mockStatusCode: 200,
			expectedError:  "API returned error: NOTOK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createMockServer(tt.mockResponse, tt.mockStatusCode)
			defer server.Close()

			client := createClientWithMockServer(server)
			result, err := client.GetLogsByTopics(tt.options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Status)
				assert.Equal(t, "OK", result.Message)
				assert.Len(t, result.Result, 1)
			}
		})
	}
}

func TestLogsClient_GetLogsByAddressAndTopics(t *testing.T) {
	tests := []struct {
		name           string
		options        GetLogsByAddressAndTopicsOptions
		mockResponse   interface{}
		mockStatusCode int
		expectedError  string
	}{
		{
			name: "成功获取日志 - 基本参数",
			options: GetLogsByAddressAndTopicsOptions{
				Address:   "0x59728544b08ab483533076417fbbb2fd0b17ce3a",
				FromBlock: 15073139,
				ToBlock:   15074139,
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "成功获取日志 - 包含所有参数",
			options: GetLogsByAddressAndTopicsOptions{
				Address:    "0x59728544b08ab483533076417fbbb2fd0b17ce3a",
				FromBlock:  15073139,
				ToBlock:    15074139,
				Topic0:     mo.Some("0x27c4f0403323142b599832f26acd21c74a9e5b809f2215726e244a4ac588cd7d"),
				Topic1:     mo.Some("0x00000000000000000000000023581767a106ae21c074b2276d25e5c3e136a68b"),
				Topic01Opr: mo.Some("and"),
				Page:       mo.Some(1),
				Offset:     mo.Some(1000),
				ChainID:    mo.Some(1),
			},
			mockResponse:   mockLogResponse,
			mockStatusCode: 200,
			expectedError:  "",
		},
		{
			name: "地址参数为空",
			options: GetLogsByAddressAndTopicsOptions{
				Address:   "",
				FromBlock: 15073139,
				ToBlock:   15074139,
			},
			mockResponse:   nil,
			mockStatusCode: 200,
			expectedError:  "address parameter is required",
		},
		{
			name: "API 返回错误状态",
			options: GetLogsByAddressAndTopicsOptions{
				Address:   "0x59728544b08ab483533076417fbbb2fd0b17ce3a",
				FromBlock: 15073139,
				ToBlock:   15074139,
			},
			mockResponse:   mockErrorResponse,
			mockStatusCode: 200,
			expectedError:  "API returned error: NOTOK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "address parameter is required" {
				client := NewEtherscanClient("test-api-key")
				result, err := client.GetLogsByAddressAndTopics(tt.options)

				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			server := createMockServer(tt.mockResponse, tt.mockStatusCode)
			defer server.Close()

			client := createClientWithMockServer(server)
			result, err := client.GetLogsByAddressAndTopics(tt.options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Status)
				assert.Equal(t, "OK", result.Message)
				assert.Len(t, result.Result, 1)
			}
		})
	}
}

func TestLogsClient_buildQueryParams(t *testing.T) {
	client := NewEtherscanClient("test-api-key")

	tests := []struct {
		name     string
		params   map[string]interface{}
		expected map[string]string
	}{
		{
			name: "基本参数",
			params: map[string]interface{}{
				"address":   "0x123",
				"fromBlock": 100,
			},
			expected: map[string]string{
				"module":    "logs",
				"action":    "getLogs",
				"apikey":    "test-api-key",
				"address":   "0x123",
				"fromBlock": "100",
			},
		},
		{
			name: "包含可选参数",
			params: map[string]interface{}{
				"address":   "0x123",
				"fromBlock": mo.Some(100),
				"toBlock":   mo.Some(200),
				"page":      mo.Some(1),
				"emptyStr":  "",
			},
			expected: map[string]string{
				"module":    "logs",
				"action":    "getLogs",
				"apikey":    "test-api-key",
				"address":   "0x123",
				"fromBlock": "100",
				"toBlock":   "200",
				"page":      "1",
			},
		},
		{
			name: "空的可选参数不应包含",
			params: map[string]interface{}{
				"address":     "0x123",
				"emptyOption": mo.None[string](),
				"emptyInt":    mo.None[int](),
			},
			expected: map[string]string{
				"module":  "logs",
				"action":  "getLogs",
				"apikey":  "test-api-key",
				"address": "0x123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.buildQueryParams(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// 测试请求参数验证
func TestRequestParameterValidation(t *testing.T) {
	// 创建一个验证请求参数的模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证查询参数
		query := r.URL.Query()

		// 验证必需的参数存在
		assert.Equal(t, "logs", query.Get("module"))
		assert.Equal(t, "getLogs", query.Get("action"))
		assert.Equal(t, "test-api-key", query.Get("apikey"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockLogResponse)
	}))
	defer server.Close()

	client := createClientWithMockServer(server)

	// 测试地址查询的参数
	_, err := client.GetLogsByAddress(GetLogsByAddressOptions{
		Address:   "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
		FromBlock: mo.Some(uint64(12878196)),
		ToBlock:   mo.Some(uint64(12878196)),
		Page:      mo.Some(1),
		Offset:    mo.Some(1000),
	})

	assert.NoError(t, err)
}

// 基准测试
func BenchmarkLogsClient_GetLogsByAddress(b *testing.B) {
	server := createMockServer(mockLogResponse, 200)
	defer server.Close()

	client := createClientWithMockServer(server)
	options := GetLogsByAddressOptions{
		Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetLogsByAddress(options)
		require.NoError(b, err)
	}
}

// 集成测试示例（需要真实的 API 密钥才能运行）
func TestLogsClient_Integration(t *testing.T) {
	t.Skip("跳过集成测试，需要真实的 API 密钥")

	apiKey := "api-key" // 替换为真实的 API 密钥
	client := NewEtherscanClient(apiKey)

	// 测试真实的 API 调用
	result, err := client.GetLogsByAddress(GetLogsByAddressOptions{
		Address:   "0x4d7aE0515784BB40E44551f68c3EFd81C00B085B",
		FromBlock: mo.Some(uint64(8205123)),
		ToBlock:   mo.Some(uint64(8216953)),
		Page:      mo.Some(1),
		Offset:    mo.Some(10),
		ChainID:   mo.Some(uint64(11155111)),
	})

	spew.Dump(result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Status)
}

// 测试 JSON 序列化和反序列化（包含数据）
func TestLogEntry_JSONSerializationWithData(t *testing.T) {
	// 包含实际数据的 JSON
	rawJSON := `{
		"address": "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
		"topics": [
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
		],
		"data": "0x000000000000000000000000000000000000000000000000000000000000001",
		"blockNumber": "0xc48174",
		"timeStamp": "0x60f9ce56",
		"gasPrice": "0x2e90edd000",
		"gasUsed": "0x247205",
		"logIndex": "0x0",
		"transactionHash": "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9",
		"transactionIndex": "0x0"
	}`

	var entry LogEntry
	err := json.Unmarshal([]byte(rawJSON), &entry)
	assert.NoError(t, err)

	// 验证 Data 字段解析
	expectedData := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	assert.Equal(t, expectedData, entry.Data)
	assert.Len(t, entry.Data, 32) // 32 字节的数据
}

// 测试 JSON 序列化和反序列化
func TestLogEntry_JSONSerialization(t *testing.T) {
	rawJSON := `{
		"address": "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
		"topics": [
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"0x0000000000000000000000000000000000000000000000000000000000000000"
		],
		"data": "0x",
		"blockNumber": "0xc48174",
		"timeStamp": "0x60f9ce56",
		"gasPrice": "0x2e90edd000",
		"gasUsed": "0x247205",
		"logIndex": "0x",
		"transactionHash": "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9",
		"transactionIndex": "0x"
	}`

	var entry LogEntry
	err := json.Unmarshal([]byte(rawJSON), &entry)
	assert.NoError(t, err)

	// 验证解析结果
	assert.Equal(t, "0xbd3531da5cf5857e7cfaa92426877b022e612cf8", entry.Address)
	assert.Len(t, entry.Topics, 2)

	expectedTopic0, _ := NewHashFromHex("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	expectedTopic1, _ := NewHashFromHex("0x0000000000000000000000000000000000000000000000000000000000000000")
	expectedTxHash, _ := NewHashFromHex("0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9")

	assert.Equal(t, expectedTopic0, entry.Topics[0])
	assert.Equal(t, expectedTopic1, entry.Topics[1])
	assert.Equal(t, []byte{}, entry.Data)                 // "0x" 解析为空字节数组
	assert.Equal(t, uint64(12878196), entry.BlockNumber)  // 0xc48174 = 12878196
	assert.Equal(t, uint64(1626984022), entry.TimeStamp)  // 0x60f9ce56 = 1626984022
	assert.Equal(t, uint64(200000000000), entry.GasPrice) // 0x2e90edd000 = 200000000000
	assert.Equal(t, uint64(2388485), entry.GasUsed)       // 0x247205 = 2388485
	assert.Equal(t, uint64(0), entry.LogIndex)            // 0x = 0
	assert.Equal(t, expectedTxHash, entry.TransactionHash)
	assert.Equal(t, uint64(0), entry.TransactionIndex) // 0x = 0
}

// 测试十六进制解析函数
func TestParseHexToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		hasError bool
	}{
		{"空字符串", "", 0, false},
		{"仅0x", "0x", 0, false},
		{"零值", "0x0", 0, false},
		{"小数值", "0xa", 10, false},
		{"大数值", "0xc48174", 12878196, false},
		{"最大值", "0x2e90edd000", 200000000000, false},
		{"无0x前缀", "c48174", 12878196, false},
		{"无效十六进制", "0xgg", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHexToUint64(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// 测试十六进制转字节数组函数
func TestParseHexToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		hasError bool
	}{
		{"空字符串", "", []byte{}, false},
		{"仅0x", "0x", []byte{}, false},
		{"零值", "0x0", []byte{0x00}, false},
		{"单字节", "0xff", []byte{0xff}, false},
		{"多字节", "0x1234", []byte{0x12, 0x34}, false},
		{"长数据", "0x000000000000000000000000000000000000000000000000000000000000001", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, false},
		{"无0x前缀", "1234", []byte{0x12, 0x34}, false},
		{"奇数长度（自动补零）", "0x123", []byte{0x01, 0x23}, false},
		{"无效十六进制", "0xgg", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHexToBytes(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// 测试 Hash 类型
func TestHash_Methods(t *testing.T) {
	// 测试从十六进制字符串创建 Hash
	hexStr := "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9"
	hash, err := NewHashFromHex(hexStr)
	assert.NoError(t, err)

	// 测试 String 方法
	assert.Equal(t, hexStr, hash.String())

	// 测试 Hex 方法（不带0x前缀）
	assert.Equal(t, "4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9", hash.Hex())

	// 测试 Bytes 方法
	bytes := hash.Bytes()
	assert.Len(t, bytes, 32)

	// 测试 IsZero 方法
	assert.False(t, hash.IsZero())

	var zeroHash Hash
	assert.True(t, zeroHash.IsZero())
}

// 测试 Hash JSON 序列化
func TestHash_JSONSerialization(t *testing.T) {
	originalHex := "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9"
	hash, err := NewHashFromHex(originalHex)
	assert.NoError(t, err)

	// 序列化
	data, err := json.Marshal(hash)
	assert.NoError(t, err)

	// 反序列化
	var deserializedHash Hash
	err = json.Unmarshal(data, &deserializedHash)
	assert.NoError(t, err)

	// 验证一致性
	assert.Equal(t, hash, deserializedHash)
	assert.Equal(t, originalHex, deserializedHash.String())
}

// 测试 NewHashFromHex 的各种情况
func TestNewHashFromHex(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"完整哈希", "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9", false},
		{"无0x前缀", "4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a9", false},
		{"短哈希（自动补零）", "0x123", false},
		{"零哈希", "0x0", false},
		{"太长的哈希", "0x4ffd22d986913d33927a392fe4319bcd2b62f3afe1c15a2c59f77fc2cc4c20a900", true},
		{"无效十六进制", "0xgg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHashFromHex(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// 测试并发安全性
func TestLogsClient_ConcurrentRequests(t *testing.T) {
	server := createMockServer(mockLogResponse, 200)
	defer server.Close()

	client := createClientWithMockServer(server)

	// 启动多个 goroutine 并发请求
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			_, err := client.GetLogsByAddress(GetLogsByAddressOptions{
				Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
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
