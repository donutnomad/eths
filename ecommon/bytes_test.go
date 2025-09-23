package ecommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// BytesRequest 用于测试 gin bind 的结构体
type BytesRequest struct {
	HexData string `json:"hex_data"`
}

func TestGinBindHexBytes(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedData []byte
	}{
		{
			name:         "有0x前缀的十六进制",
			input:        "0x48656c6c6f",
			expectedData: []byte("Hello"),
		},
		{
			name:         "无0x前缀的十六进制",
			input:        "576f726c64",
			expectedData: []byte("World"),
		},
		{
			name:         "奇数长度的十六进制",
			input:        "0x123",
			expectedData: []byte{0x01, 0x23},
		},
		{
			name:         "空字符串",
			input:        "",
			expectedData: []byte{},
		},
		{
			name:         "只有0x前缀",
			input:        "0x",
			expectedData: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeRequest[BytesRequest](map[string]string{
				"hex_data": tt.input,
			}, func(req BytesRequest) {
				data := FromHex(req.HexData)
				assert.Equal(t, tt.expectedData, data)
				assert.Equal(t, len(tt.expectedData), len(data))
			})
		})
	}
}

func TestGinBindEmptyField(t *testing.T) {
	makeRequest[BytesRequest](map[string]interface{}{}, func(req BytesRequest) {
		data := FromHex(req.HexData)
		assert.Equal(t, []byte{}, data)
		assert.Equal(t, 0, len(data))
	})
}

func TestGinBindInvalidJSON(t *testing.T) {
	assert.Panics(t, func() {
		makeRequest[BytesRequest]("invalid json", func(req BytesRequest) {
			t.Error("不应该调用callback，因为gin bind应该失败")
		})
	})
}
