package createx

import (
	"crypto/rand"
	"github.com/donutnomad/eths/contracts_pack"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
	"math/big"
	"slices"
)

// Address hardhat createX's address
// Ethereum:
// https://etherscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// https://sepolia.etherscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// https://hoodi.etherscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// Arbitrum:
// https://arbiscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// https://sepolia.arbiscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// https://nova.arbiscan.io/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
// BSC:
// https://bscscan.com/address/0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed#code
var Address = common.HexToAddress("0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed")

// GenSaltZeroAddressRedeployProtection
// senderBytes == SenderBytes.ZeroAddress && redeployProtectionFlag == RedeployProtectionFlag.True
func GenSaltZeroAddressRedeployProtection(chainID *big.Int) (pre, post [32]byte) {
	var preSalt = buildSalt(common.Address{}, true)

	var chainIDBytes32 [32]byte
	fillRight(chainIDBytes32[:], chainID)

	afterSalt := [32]byte(crypto.Keccak256(slices.Concat(chainIDBytes32[:], preSalt[:])))
	return preSalt, afterSalt
}

func Create2(salt [32]byte, initCode []byte) []byte {
	return contracts_pack.NewCreatex().PackDeployCreate2(salt, initCode)
}

func Create2WithoutSalt(initCode []byte) []byte {
	return contracts_pack.NewCreatex().PackDeployCreate20(initCode)
}

func genSaltZeroAddressRedeployProtection(preSalt [32]byte, chainID *big.Int) (pre, post [32]byte) {
	var chainIDBytes32 [32]byte
	fillRight(chainIDBytes32[:], chainID)

	afterSalt := [32]byte(crypto.Keccak256(slices.Concat(chainIDBytes32[:], preSalt[:])))
	return preSalt, afterSalt
}

func fillRight(input []byte, value *big.Int) {
	bs := value.Bytes()
	copy(input[32-len(bs):], bs)
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
