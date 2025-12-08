package contractcall

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// hasherPool holds LegacyKeccak256 buffer for rlpHash.
var hasherPool = sync.Pool{
	New: func() any { return crypto.NewKeccakState() },
}

// PrefixedRlpHash writes the prefix into the hasher before rlp-encoding x.
// It's used for typed transactions.
func PrefixedRlpHash(prefix []byte, x any) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	if len(prefix) > 0 {
		sha.Write(prefix)
	}
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}

//   rlp:"nil" 标签用于指针类型字段，允许空值（empty value）被解码为 nil 指针。
//   // ✅ 正确 - 用于指针类型
//  type Transaction struct {
//      To *common.Address `rlp:"nil"`
//  }
//
//  // ❌ 错误 - 不能用于非指针类型
//  type Invalid struct {
//      X []byte `rlp:"nil"`  // 编译时会报错
//  }
//   2. 两种空值类型
//
//  RLP 支持两种空值，rlp:"nil" 会根据字段类型自动选择：
//
//  | 字段类型                                     | 空值类型 | RLP 编码 |
//  |------------------------------------------|------|--------|
//  | *uint, *string, *bool, *[N]byte(common.Address), *[]byte | 空字符串 | 0x80(128)   |
//  | *struct, 其他复杂类型                          | 空列表  | 0xC0(192)   |
// 相关标签
//
//  - rlp:"nil" - 自动选择空值类型（字符串或列表）
//  - rlp:"nilString" - 强制使用空字符串 (0x80)
//  - rlp:"nilList" - 强制使用空列表 (0xC0)
