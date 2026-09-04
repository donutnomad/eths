package createx

import (
	"crypto/rand"
	"math/big"
	"slices"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/crypto"
	"github.com/samber/lo"
)

func fillRight(input []byte, value *big.Int) {
	bs := value.Bytes()
	copy(input[32-len(bs):], bs)
}

// computeCreate3Address 计算通过 CREATE3 模式部署的合约地址。
// CREATE3 地址不依赖 initCode，仅取决于 salt 和 deployer 地址。
//
// 计算逻辑（与 CreateX 合约一致）：
//  1. 用 CREATE2 + proxy bytecode hash 计算 proxy 地址
//  2. 用 CREATE(proxy, nonce=1) 计算最终合约地址
func computeCreate3Address(salt [32]byte, deployer common.Address) common.Address {
	// CreateX 使用的 proxy bytecode hash: keccak256(hex"67_36_3d_3d_37_36_3d_34_f0_3d_52_60_08_60_18_f3")
	proxyBytecodeHash := common.HexToHash("0x21c35dbe1b344a2488cf3321d6ce542f8e9f305544ff09e4993a62319a497c1f")
	proxyAddress := crypto.CreateAddress2(deployer, salt, proxyBytecodeHash.Bytes())
	return crypto.CreateAddress(proxyAddress, 1)
}

// genSaltZeroAddressRedeployProtection 对应 CreateX _guard 中 ZeroAddress + True 分支：
// guardedSalt = keccak256(chainID || salt)
func genSaltZeroAddressRedeployProtection(preSalt [32]byte, chainID *big.Int) (pre, post [32]byte) {
	var chainIDBytes32 [32]byte
	fillRight(chainIDBytes32[:], chainID)

	afterSalt := [32]byte(crypto.Keccak256(slices.Concat(chainIDBytes32[:], preSalt[:])))
	return preSalt, afterSalt
}

func buildSalt(sender common.Address, protection bool) [32]byte {
	var salt [32]byte
	copy(salt[0:20], sender[:])
	if protection {
		salt[20] = 0x1
	}
	lo.Must1(rand.Read(salt[21:]))
	return salt
}
