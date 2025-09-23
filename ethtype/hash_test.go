package ethtype

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

// HashRequest 用于测试 gin bind 的结构体
type HashRequest struct {
	HashValue Hash `json:"hash_value"`
}

// HashCompareRequest 哈希比较请求
type HashCompareRequest struct {
	Hash1 Hash `json:"hash1" form:"hash1"`
	Hash2 Hash `json:"hash2" form:"hash2"`
}

func TestMarshalHash(t *testing.T) {
	var req = HashRequest{
		HashValue: MaxHash,
	}
	marshal, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(marshal))
}

func TestGinBindHash(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedHex string
	}{
		{
			name:        "有0x前缀的32字节哈希",
			input:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expectedHex: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:        "另一个32字节哈希",
			input:       "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
			expectedHex: "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
		},
		{
			name:        "全零哈希",
			input:       "0x0000000000000000000000000000000000000000000000000000000000000000",
			expectedHex: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:        "最大哈希值",
			input:       "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			expectedHex: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeRequest(map[string]string{
				"hash_value": tt.input,
			}, func(req HashRequest) {
				assert.Equal(t, tt.expectedHex, req.HashValue.Hex())
				assert.Equal(t, HashLength, len(req.HashValue.Bytes()))

				expectedHash := HexToHash(tt.input)
				assert.Equal(t, expectedHash, req.HashValue)
			})
		})
	}
}

func TestGinBindHashInvalidFormat(t *testing.T) {
	invalidTests := []struct {
		name  string
		input string
	}{
		{
			name:  "无0x前缀的哈希",
			input: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:  "短哈希",
			input: "0x123",
		},
		{
			name:  "过长的哈希（超过32字节）",
			input: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef00",
		},
		{
			name:  "包含非十六进制字符",
			input: "0x123456789gabcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		},
		{
			name:  "空字符串",
			input: "",
		},
		{
			name:  "只有0x前缀",
			input: "0x",
		},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				makeRequest(map[string]string{
					"hash_value": tt.input,
				}, func(req HashRequest) {
					t.Error("不应该调用callback，因为gin bind应该失败")
				})
			})
		})
	}
}

func TestGinBindHashMissingField(t *testing.T) {
	makeRequest(map[string]interface{}{}, func(req HashRequest) {
		expectedZeroHash := Hash{}
		assert.Equal(t, expectedZeroHash, req.HashValue)
		assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", req.HashValue.Hex())
	})
}

func TestGinBindHashInvalidJSON(t *testing.T) {
	assert.Panics(t, func() {
		makeRequest("invalid json", func(req HashRequest) {
			t.Error("不应该调用callback，因为gin bind应该失败")
		})
	})
}

func TestGinBindHashComparison(t *testing.T) {
	makeRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req HashCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

func TestForm(t *testing.T) {
	makeFormRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req HashCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

func TestFormFailed(t *testing.T) {
	assert.Panics(t, func() {
		makeFormRequest(map[string]string{
			"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefed",
			"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		}, func(req HashCompareRequest) {

		})
	})
}

// 正常
func TestQuery(t *testing.T) {
	makeGetRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req HashCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

// 超过32个字节，报错
func TestQuery2(t *testing.T) {
	assert.Panics(t, func() {
		makeGetRequest(map[string]string{
			"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefed",
			"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		}, func(req HashCompareRequest) {
			spew.Dump(req)
			//cmp := req.Hash1.Cmp(req.Hash2)
			//assert.Equal(t, 0, cmp)
			//assert.True(t, cmp == 0)
			//assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
			//assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
		})
	})
}
