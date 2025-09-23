package ecommon

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

// Big json序列化和反序列化时，转换为字符串类型"xxxxxx"
// 支持gin binding, sql
type Big big.Int
type BigInt = Big

func NewInt(x int64) *Big {
	return (*Big)(big.NewInt(x))
}

func NewIntFromString(input string) (*Big, error) {
	v, ok := new(big.Int).SetString(input, 10)
	if !ok {
		return nil, fmt.Errorf("invalid bigInt %s", input)
	}
	return (*Big)(v), nil
}

func (x *Big) IsZero() bool {
	return x.ToInt().Sign() == 0
}

// AppendText implements the [encoding.TextAppender] interface.
func (x *Big) AppendText(buf []byte) (text []byte, err error) {
	return x.ToInt().AppendText(buf)
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (x *Big) MarshalText() (text []byte, err error) {
	return x.AppendText(nil)
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (x *Big) UnmarshalText(text []byte) error {
	return x.ToInt().UnmarshalText(text)
}

// MarshalJSON implements the [encoding/json.Marshaler] interface.
func (x *Big) MarshalJSON() ([]byte, error) {
	str, err := x.MarshalText()
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, len(str)+2)
	result = append(result, '"')
	result = append(result, str...)
	result = append(result, '"')
	return result, nil
}

// UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
func (x *Big) UnmarshalJSON(text []byte) error {
	// Ignore null, like in the main JSON package.
	if string(text) == "null" {
		return nil
	}
	if isString(text) {
		text = text[1 : len(text)-1]
	}
	return x.ToInt().UnmarshalJSON(text)
}

// ToInt converts b to a big.Int.
func (x *Big) ToInt() *big.Int {
	return (*big.Int)(x)
}

func (x *Big) BigInt() *big.Int {
	return (*big.Int)(x)
}

// String returns the hex encoding of b.
func (x *Big) String() string {
	return x.ToInt().Text(10)
}

// UnmarshalParam implement gin binding.BindUnmarshaler
func (x *Big) UnmarshalParam(param string) error {
	return x.UnmarshalText([]byte(param))
}

// Scan implements Scanner for database/sql.
func (x *Big) Scan(src any) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case string:
		return x.UnmarshalText([]byte(v))
	case []byte:
		return x.UnmarshalText(v)
	default:
		return fmt.Errorf("can't scan %T into Big", src)
	}
}

// Value implements valuer for database/sql.
func (x Big) Value() (driver.Value, error) {
	return x.MarshalText()
}

// SetInt64 sets x to y and returns x.
func (x *Big) SetInt64(y int64) *Big {
	x.ToInt().SetInt64(y)
	return x
}

// SetUint64 sets x to y and returns x.
func (x *Big) SetUint64(y uint64) *Big {
	x.ToInt().SetUint64(y)
	return x
}

// Set sets x to y and returns x.
func (x *Big) Set(y *Big) *Big {
	x.ToInt().Set(y.ToInt())
	return x
}

// SetBits sets x to abs and returns x.
func (x *Big) SetBits(abs []big.Word) *Big {
	x.ToInt().SetBits(abs)
	return x
}

// Neg sets x to -y and returns x.
func (x *Big) Neg(y *Big) *Big {
	x.ToInt().Neg(y.ToInt())
	return x
}

// Add sets x to the sum x+y and returns x.
func (x *Big) Add(y, z *Big) *Big {
	x.ToInt().Add(y.ToInt(), z.ToInt())
	return x
}

// Sub sets x to the difference x-y and returns x.
func (x *Big) Sub(y, z *Big) *Big {
	x.ToInt().Sub(y.ToInt(), z.ToInt())
	return x
}

// Mul sets x to the product x*y and returns x.
func (x *Big) Mul(y, z *Big) *Big {
	x.ToInt().Mul(y.ToInt(), z.ToInt())
	return x
}

// MulRange sets x to the product of all integers in the range [a, b] and returns x.
func (x *Big) MulRange(a, b int64) *Big {
	x.ToInt().MulRange(a, b)
	return x
}

// Div sets x to the quotient x/y for y != 0 and returns x.
func (x *Big) Div(y, z *Big) *Big {
	x.ToInt().Div(y.ToInt(), z.ToInt())
	return x
}

// Mod sets x to the modulus x%y for y != 0 and returns x.
func (x *Big) Mod(y, z *Big) *Big {
	x.ToInt().Mod(y.ToInt(), z.ToInt())
	return x
}

// DivMod sets x to the quotient x/y and m to the modulus x%y and returns the pair (x, m).
func (x *Big) DivMod(y, z, m *Big) (*Big, *Big) {
	x.ToInt().DivMod(y.ToInt(), z.ToInt(), m.ToInt())
	return x, m
}

// Cmp compares x and y and returns -1, 0, or +1.
func (x *Big) Cmp(y *Big) (r int) {
	return x.ToInt().Cmp(y.ToInt())
}

// CmpAbs compares the absolute values of x and y and returns -1, 0, or +1.
func (x *Big) CmpAbs(y *Big) int {
	return x.ToInt().CmpAbs(y.ToInt())
}

// SetString sets x to the value of s and returns x and a boolean indicating success.
func (x *Big) SetString(s string, base int) (*Big, bool) {
	_, ok := x.ToInt().SetString(s, base)
	return x, ok
}

// SetBytes sets x to the value of buf and returns x.
func (x *Big) SetBytes(buf []byte) *Big {
	x.ToInt().SetBytes(buf)
	return x
}

// Exp sets x = y**z mod |m| and returns x.
func (x *Big) Exp(y, z, m *Big) *Big {
	x.ToInt().Exp(y.ToInt(), z.ToInt(), m.ToInt())
	return x
}

// ModInverse sets x to the multiplicative inverse of g in the ring ℤ/nℤ and returns x.
func (x *Big) ModInverse(g, n *Big) *Big {
	x.ToInt().ModInverse(g.ToInt(), n.ToInt())
	return x
}

// ModSqrt sets x to a square root of a mod p if such a square root exists, and returns x.
func (x *Big) ModSqrt(a, p *Big) *Big {
	x.ToInt().ModSqrt(a.ToInt(), p.ToInt())
	return x
}

// Lsh sets x = y << n and returns x.
func (x *Big) Lsh(y *Big, n uint) *Big {
	x.ToInt().Lsh(y.ToInt(), n)
	return x
}

// Rsh sets x = y >> n and returns x.
func (x *Big) Rsh(y *Big, n uint) *Big {
	x.ToInt().Rsh(y.ToInt(), n)
	return x
}

// SetBit sets x to y, with y's i'th bit set to b (0 or 1), and returns x.
func (x *Big) SetBit(y *Big, i int, b uint) *Big {
	x.ToInt().SetBit(y.ToInt(), i, b)
	return x
}

// And sets x = x & y and returns x.
func (x *Big) And(y, z *Big) *Big {
	x.ToInt().And(y.ToInt(), z.ToInt())
	return x
}

// AndNot sets x = x &^ y and returns x.
func (x *Big) AndNot(y, z *Big) *Big {
	x.ToInt().AndNot(y.ToInt(), z.ToInt())
	return x
}

// Or sets x = x | y and returns x.
func (x *Big) Or(y, z *Big) *Big {
	x.ToInt().Or(y.ToInt(), z.ToInt())
	return x
}

// Xor sets x = x ^ y and returns x.
func (x *Big) Xor(y, z *Big) *Big {
	x.ToInt().Xor(y.ToInt(), z.ToInt())
	return x
}

// Not sets x = ^y and returns x.
func (x *Big) Not(y *Big) *Big {
	x.ToInt().Not(y.ToInt())
	return x
}

// Sqrt sets x to ⌊√y⌋ and returns x.
func (x *Big) Sqrt(y *Big) *Big {
	x.ToInt().Sqrt(y.ToInt())
	return x
}
