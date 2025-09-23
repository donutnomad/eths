package ecommon

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
)

var (
	hashT = reflect.TypeOf(Hash{})
	// MaxHash represents the maximum possible hash value.
	MaxHash = HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
// 输出的是带0x的，输入只能是带有0x的，对于数据库（可以带0x或者不带0x，特殊处理了)
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Cmp compares two hashes.
func (h Hash) Cmp(other Hash) int {
	return bytes.Compare(h[:], other[:])
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return EncodeHexToString(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x..%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter.
// Hash supports the %v, %s, %q, %x, %X and %d format verbs.
func (h Hash) Format(s fmt.State, c rune) {
	hexb := make([]byte, 2+len(h)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], h[:])

	switch c {
	case 'x', 'X':
		if !s.Flag('#') {
			hexb = hexb[2:]
		}
		if c == 'X' {
			hexb = bytes.ToUpper(hexb)
		}
		fallthrough
	case 'v', 's':
		s.Write(hexb)
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(hexb)
		s.Write(q)
	case 'd':
		fmt.Fprint(s, ([len(h)]byte)(h))
	default:
		fmt.Fprintf(s, "%%!%c(hash=%x)", c, h)
	}
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return EncodeHexToBytes(h[:]), nil
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return UnmarshalFixedUnprefixedText("Hash", []byte(v), h[:])
	case []byte:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return UnmarshalFixedUnprefixedText("Hash", v, h[:])
	default:
		return fmt.Errorf("can't scan %T into Hash", src)
	}
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h.Hex(), nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (Hash) ImplementsGraphQLType(name string) bool { return name == "Bytes32" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input any) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Hash", input)
	}
	return err
}

// UnmarshalParam implement gin binding.BindUnmarshaler
func (h *Hash) UnmarshalParam(param string) error {
	return h.UnmarshalText([]byte(param))
}
