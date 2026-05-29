package createx

import (
	"crypto/rand"
	"math/big"
	"slices"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/contracts_pack"
	"github.com/donutnomad/eths/crypto"
	"github.com/samber/lo"
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

func Create3(salt [32]byte, initCode []byte) []byte {
	return contracts_pack.NewCreatex().PackDeployCreate30(salt, initCode)
}

func Create3WithoutSalt(initCode []byte) []byte {
	return contracts_pack.NewCreatex().PackDeployCreate3(initCode)
}

// ComputeCreate3Address 计算通过 CREATE3 模式部署的合约地址。
// CREATE3 地址不依赖 initCode，仅取决于 salt 和 deployer 地址。
//
// 计算逻辑（与 CreateX 合约一致）：
//  1. 用 CREATE2 + proxy bytecode hash 计算 proxy 地址
//  2. 用 CREATE(proxy, nonce=1) 计算最终合约地址
func ComputeCreate3Address(salt [32]byte, deployer common.Address) common.Address {
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

// genSaltMsgSenderRedeployProtection 对应 CreateX _guard 中 MsgSender + True 分支：
// guardedSalt = keccak256(abi.encode(msg.sender, block.chainid, salt))
func genSaltMsgSenderRedeployProtection(preSalt [32]byte, sender common.Address, chainID *big.Int) (pre, post [32]byte) {
	var senderBytes32 [32]byte
	copy(senderBytes32[12:], sender[:])

	var chainIDBytes32 [32]byte
	fillRight(chainIDBytes32[:], chainID)

	afterSalt := [32]byte(crypto.Keccak256(slices.Concat(senderBytes32[:], chainIDBytes32[:], preSalt[:])))
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

// DeriveSalt 通过 masterSalt 和 path 确定性生成兼容 CreateX 安全结构的 salt。
//
//	salt[0:20]  = sender
//	salt[20]    = protection ? 0x01 : 0x00
//	salt[21:32] = keccak256(masterSalt || path)[:11]
func DeriveSalt(sender common.Address, protection bool, masterSalt [32]byte, path string) [32]byte {
	var salt [32]byte
	copy(salt[0:20], sender[:])
	if protection {
		salt[20] = 0x01
	}
	entropy := crypto.Keccak256(slices.Concat(masterSalt[:], []byte(path)))
	copy(salt[21:32], entropy[:11])
	return salt
}

// Deployer 封装 CreateX 部署操作，支持确定性 salt 派生。
type Deployer struct {
	sender     common.Address
	protection bool
	masterSalt [32]byte
	chainID    *big.Int
}

// NewDeployerZeroAddress 创建零地址 + 跨链保护模式的 Deployer。
// 任何人都可以用此 salt 部署，但不同链上会产生不同地址。
func NewDeployerZeroAddress(masterSalt [32]byte, chainID *big.Int) *Deployer {
	return &Deployer{
		sender:     common.Address{},
		protection: true,
		masterSalt: masterSalt,
		chainID:    chainID,
	}
}

// NewDeployerZeroAddressNoProtection 创建零地址 + 无跨链保护模式的 Deployer。
// 任何人都可以用此 salt 部署，且同一 salt 在所有链上产生相同地址。
func NewDeployerZeroAddressNoProtection(masterSalt [32]byte) *Deployer {
	return &Deployer{
		sender:     common.Address{},
		protection: false,
		masterSalt: masterSalt,
	}
}

// NewDeployerSender 创建指定发送者 + 跨链保护模式的 Deployer。
// 只有指定的 sender 才能用此 salt 部署，且不同链上产生不同地址。
func NewDeployerSender(sender common.Address, masterSalt [32]byte, chainID *big.Int) *Deployer {
	return &Deployer{
		sender:     sender,
		protection: true,
		masterSalt: masterSalt,
		chainID:    chainID,
	}
}

// Salt 返回指定 path 的原始 salt（pre）和经 _guard 处理后的 salt（post）。
func (d *Deployer) Salt(path string) (pre, post [32]byte) {
	preSalt := DeriveSalt(d.sender, d.protection, d.masterSalt, path)
	if d.sender == (common.Address{}) && d.protection {
		// ZeroAddress + True: guardedSalt = keccak256(chainID || salt)
		return genSaltZeroAddressRedeployProtection(preSalt, d.chainID)
	}
	if d.sender == (common.Address{}) && !d.protection {
		// ZeroAddress + False: guardedSalt = keccak256(salt)
		afterSalt := [32]byte(crypto.Keccak256(preSalt[:]))
		return preSalt, afterSalt
	}
	// MsgSender + True: guardedSalt = keccak256(abi.encode(sender, chainID, salt))
	return genSaltMsgSenderRedeployProtection(preSalt, d.sender, d.chainID)
}

// Create2 返回 deployCreate2 的 calldata 和预计算的部署地址。
func (d *Deployer) Create2(path string, initCode []byte) (calldata []byte, addr common.Address) {
	pre, post := d.Salt(path)
	calldata = contracts_pack.NewCreatex().PackDeployCreate2(pre, initCode)
	addr = crypto.CreateAddress2(Address, post, crypto.Keccak256Hash(initCode).Bytes())
	return
}

// Create3 返回 deployCreate3 的 calldata 和预计算的部署地址。
func (d *Deployer) Create3(path string, initCode []byte) (calldata []byte, addr common.Address) {
	pre, post := d.Salt(path)
	calldata = contracts_pack.NewCreatex().PackDeployCreate30(pre, initCode)
	addr = ComputeCreate3Address(post, Address)
	return
}

// ComputeCreate2Address 返回指定 path 和 initCode 的 Create2 预计算部署地址。
func (d *Deployer) ComputeCreate2Address(path string, initCode []byte) common.Address {
	_, post := d.Salt(path)
	return crypto.CreateAddress2(Address, post, crypto.Keccak256Hash(initCode).Bytes())
}

// ComputeCreate3Address 返回指定 path 的 Create3 预计算部署地址。
// CREATE3 地址不依赖 initCode。
func (d *Deployer) ComputeCreate3Address(path string) common.Address {
	_, post := d.Salt(path)
	return ComputeCreate3Address(post, Address)
}
