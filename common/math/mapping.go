package math

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

// HexOrDecimal256 marshals big.Int as hex or decimal.
type HexOrDecimal256 = math.HexOrDecimal256

// NewHexOrDecimal256 creates a new HexOrDecimal256
func NewHexOrDecimal256(x int64) *HexOrDecimal256 {
	return math.NewHexOrDecimal256(x)
}

// U256Bytes converts a big Int into a 256bit EVM number.
// This operation is destructive.
func U256Bytes(n *big.Int) []byte {
	return math.U256Bytes(n)
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}
