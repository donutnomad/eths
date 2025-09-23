package ecommon

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
)

var (
	bytesT = reflect.TypeOf(Bytes(nil))
)

// Bytes marshals/unmarshals as a JSON string with 0x prefix.
// The empty slice marshals as "0x".
type Bytes []byte

// Cmp compares two byte arrays.
func (b Bytes) Cmp(other Bytes) int {
	return bytes.Compare(b[:], other[:])
}

// Bytes gets the byte representation of the underlying hash.
func (b Bytes) Bytes() []byte { return b[:] }

// Big converts a hash to a big integer.
func (b Bytes) Big() *big.Int { return new(big.Int).SetBytes(b[:]) }

// Hex converts a hash to a hex string.
func (b Bytes) Hex() string { return EncodeHexToString(b[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (b Bytes) TerminalString() string {
	if len(b) <= 6 {
		return fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("%x..%x", b[:3], b[len(b)-3:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (b Bytes) String() string {
	return b.Hex()
}

// UnmarshalText parses a hash in hex syntax.
func (b *Bytes) UnmarshalText(input []byte) error {
	return b.unmarshalTextWith(input, true)
}

// UnmarshalJSON parses a hash in hex syntax.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(bytesT)
	}
	return wrapTypeError(b.UnmarshalText(input[1:len(input)-1]), bytesT)
}

// MarshalText returns the hex representation of b.
func (b Bytes) MarshalText() ([]byte, error) {
	return EncodeHexToBytes(b[:]), nil
}

// SetBytes sets the bytes to the value of input.
// If input is larger than the current capacity, the slice will be grown.
// If input is smaller, only the necessary bytes will be copied.
func (b *Bytes) SetBytes(input []byte) {
	if input == nil {
		*b = nil
		return
	}
	// 确保 *b 有足够的容量
	if len(*b) < len(input) {
		*b = make([]byte, len(input))
	} else {
		*b = (*b)[:len(input)]
	}
	copy(*b, input)
}

// Scan implements Scanner for database/sql.
func (b *Bytes) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return b.unmarshalTextWith([]byte(v), false)
	case []byte:
		// 支持不带0x前缀的解析, 对数据库的宽容处理
		return b.unmarshalTextWith(v, false)
	default:
		return fmt.Errorf("can't scan %T into Bytes", src)
	}
}

// Value implements valuer for database/sql.
func (b Bytes) Value() (driver.Value, error) {
	return b.Hex(), nil
}

// ImplementsGraphQLType returns true if Bytes implements the specified GraphQL type.
func (b Bytes) ImplementsGraphQLType(name string) bool { return name == "Bytes" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (b *Bytes) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		data, err := DecodeHexFromString(input)
		if err != nil {
			return err
		}
		*b = data
	default:
		err = fmt.Errorf("unexpected type %T for Bytes", input)
	}
	return err
}

// UnmarshalParam implement gin binding.BindUnmarshaler
func (b *Bytes) UnmarshalParam(param string) error {
	return b.UnmarshalText([]byte(param))
}

func (b *Bytes) unmarshalTextWith(input []byte, wantPrefix bool) error {
	raw, err := checkText(input, wantPrefix)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		err = mapError(err)
	} else {
		*b = dec
	}
	return err
}
