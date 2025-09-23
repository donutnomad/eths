package ethtype

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
)

var unprefixedHashT = reflect.TypeOf(UnprefixedHash{})

// UnprefixedHash represents the 32 byte Keccak256 hash of arbitrary data.
// 输出的是不带0x的，输入可以带有0x/不带0x，均可识别
type UnprefixedHash [HashLength]byte

// BytesToUnprefixedHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToUnprefixedHash(b []byte) UnprefixedHash {
	var h UnprefixedHash
	h.SetBytes(b)
	return h
}

// BigToUnprefixedHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToUnprefixedHash(b *big.Int) UnprefixedHash { return BytesToUnprefixedHash(b.Bytes()) }

// HexToUnprefixedHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToUnprefixedHash(s string) UnprefixedHash { return BytesToUnprefixedHash(FromHex(s)) }

// Cmp compares two hashes.
func (h UnprefixedHash) Cmp(other UnprefixedHash) int {
	return bytes.Compare(h[:], other[:])
}

// Bytes gets the byte representation of the underlying hash.
func (h UnprefixedHash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h UnprefixedHash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h UnprefixedHash) Hex() string {
	s := EncodeHexToString(h[:])
	return s[2:]
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h UnprefixedHash) TerminalString() string {
	return fmt.Sprintf("%x..%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h UnprefixedHash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter.
// Hash supports the %v, %s, %q, %x, %X and %d format verbs.
func (h UnprefixedHash) Format(s fmt.State, c rune) {
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
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *UnprefixedHash) UnmarshalJSON(input []byte) error {
	return UnmarshalFixedUnprefixedJSON(unprefixedHashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *UnprefixedHash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h UnprefixedHash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *UnprefixedHash) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return UnmarshalFixedUnprefixedText("UnprefixedHash", []byte(v), h[:])
	case []byte:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return UnmarshalFixedUnprefixedText("UnprefixedHash", v, h[:])
	default:
		return fmt.Errorf("can't scan %T into UnprefixedHash", src)
	}
}

// Value implements valuer for database/sql.
func (h UnprefixedHash) Value() (driver.Value, error) {
	return h.Hex(), nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (UnprefixedHash) ImplementsGraphQLType(name string) bool { return name == "Bytes32" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *UnprefixedHash) UnmarshalGraphQL(input any) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for UnprefixedHash", input)
	}
	return err
}

// UnmarshalParam implement gin binding.BindUnmarshaler
func (h *UnprefixedHash) UnmarshalParam(param string) error {
	return h.UnmarshalText([]byte(param))
}
