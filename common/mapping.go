package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Address = common.Address
type MixedcaseAddress = common.MixedcaseAddress
type Hash = common.Hash

var Big1 = common.Big1
var Big0 = common.Big0
var Big32 = common.Big32

const HashLength = common.HashLength
const AddressLength = common.AddressLength

var MaxHash = common.MaxHash
var MaxAddress = common.MaxAddress

func HexToAddress(s string) Address { return common.HexToAddress(s) }
func IsHexAddress(s string) bool {
	return common.IsHexAddress(s)
}

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return common.Bytes2Hex(d)
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	return common.BytesToAddress(b)
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	return common.BytesToHash(b)
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	return common.FromHex(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	return common.Hex2Bytes(str)
}

// RightPadBytes zero-pads slice to the right up to length l.
func RightPadBytes(slice []byte, l int) []byte {
	return common.RightPadBytes(slice, l)
}

// LeftPadBytes zero-pads slice to the left up to length l.
func LeftPadBytes(slice []byte, l int) []byte {
	return common.LeftPadBytes(slice, l)
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return common.BigToHash(b) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return common.HexToHash(s) }
