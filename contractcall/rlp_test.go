package contractcall

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/donutnomad/blockchain-alg/xecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

func TestRlp(t *testing.T) {
	x := xecdsa.Secp256k1.Curve().Params().N
	var a = new(big.Int).Set(x).Rsh(x, 2)
	// 28948022309329048855892746252171976963209391069768726095651290785379540373584
	fmt.Println(a.String())
	fmt.Println(PrefixedRlpHash(nil, a))
	b := new(uint256.Int)
	overflow := b.SetFromBig(a)
	fmt.Println(overflow)
	fmt.Println(b.String())
	fmt.Println(PrefixedRlpHash(nil, b))
}

func TestRlp2(t *testing.T) {
	var c = uint64(88889990)
	fmt.Println(PrefixedRlpHash(nil, c))
	b := new(uint256.Int)
	b.SetUint64(c)

	fmt.Println(PrefixedRlpHash(nil, new(big.Int).SetUint64(c)))
	fmt.Println(PrefixedRlpHash(nil, b))
}

func TestAddress(t *testing.T) {
	var a = common.HexToAddress("0x6168F1530d4F372d714F16DbcD51685233B59735")
	var b *common.Address = nil
	if PrefixedRlpHash(nil, a) != PrefixedRlpHash(nil, &a) {
		t.Error("not equal")
	}
	if PrefixedRlpHash(nil, b) != common.BytesToHash(crypto.Keccak256([]byte{128})) {
		t.Error("not equal 2")
	}
}
