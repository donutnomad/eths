package etherscan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Hash 表示一个 32 字节的哈希值
type Hash [32]byte

// NewHashFromHex 从十六进制字符串创建 Hash
func NewHashFromHex(hexStr string) (Hash, error) {
	var hash Hash

	// 移除 0x 前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// 检查长度，如果不足 64 个字符，在前面补零
	if len(hexStr) < 64 {
		hexStr = strings.Repeat("0", 64-len(hexStr)) + hexStr
	} else if len(hexStr) > 64 {
		return hash, fmt.Errorf("hex string too long: %d characters", len(hexStr))
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return hash, fmt.Errorf("invalid hex string: %w", err)
	}

	copy(hash[:], bytes)
	return hash, nil
}

// String 返回哈希的十六进制字符串表示（带 0x 前缀）
func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

// Hex 返回哈希的十六进制字符串表示（不带 0x 前缀）
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// Bytes 返回哈希的字节切片
func (h Hash) Bytes() []byte {
	return h[:]
}

// IsZero 检查哈希是否为零值
func (h Hash) IsZero() bool {
	return h == Hash{}
}

// MarshalJSON 实现 JSON 序列化
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 实现 JSON 反序列化
func (h *Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	hash, err := NewHashFromHex(hexStr)
	if err != nil {
		return err
	}

	*h = hash
	return nil
}
