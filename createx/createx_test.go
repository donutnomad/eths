package createx

import (
	"encoding/hex"
	"github.com/samber/lo"
	"math/big"
	"testing"
)

func TestGenSalt(t *testing.T) {
	var testcases = [][2]string{
		{"00000000000000000000000000000000000000000180365a1680362939166689", "cb91198218b4e58c2d166ec31c1fe9a2ef13402ff0dcb63f669e30d577b5bac5"},
		{"000000000000000000000000000000000000000001bf20e95140bf8406847711", "a0cc9e70e3c27ae6c0522689f4895daac65906e95b1f831d22bd51d6c790052c"},
	}
	for _, tc := range testcases {
		a, b := genSaltZeroAddressRedeployProtection(mustDecode(tc[0]), big.NewInt(31337))
		if hex.EncodeToString(b[:]) != tc[1] {
			t.Fatal(a)
		}
	}
}

func mustDecode(input string) [32]byte {
	bs := lo.Must1(hex.DecodeString(input))
	return [32]byte(bs)
}
