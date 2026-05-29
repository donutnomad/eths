// Package createx 提供 CreateX 工厂合约的 Go 封装，支持 CREATE2/CREATE3 部署及确定性 salt 派生。
//
// # CreateX 合约地址（所有支持的链上相同）：0xba5Ed099633D3B313e4D5F7bdc1305d3c28ba5Ed
//
// # Salt 结构
//
// CreateX 合约的 salt 是 32 字节，内部结构如下：
//
//	[0:20]   sender 地址（用于权限控制）
//	[20]     跨链保护标志（0x01=开启, 0x00=关闭）
//	[21:32]  11 字节熵（随机或派生）
//
// 合约的 _guard 函数会根据前 21 字节的值选择不同的哈希策略，生成最终的 guardedSalt：
//
//	| sender      | protection | guard 逻辑                                          |
//	|-------------|------------|-----------------------------------------------------|
//	| MsgSender   | True       | keccak256(abi.encode(sender, chainid, salt))        |
//	| MsgSender   | False      | keccak256(sender || salt)                           |
//	| ZeroAddress  | True       | keccak256(chainid || salt)                          |
//	| ZeroAddress  | False      | keccak256(salt)                                     |
//
// # 两种使用模式
//
// ## 1. 随机模式（旧）
//
// 使用 [GenSaltZeroAddressRedeployProtection] 生成随机 salt，每次调用产生不同的值。
// 适合不需要预知或重现部署地址的场景。
//
//	pre, post := createx.GenSaltZeroAddressRedeployProtection(chainID)
//	calldata := createx.Create2(pre, initCode)
//	addr := crypto.CreateAddress2(createx.Address, post, crypto.Keccak256Hash(initCode).Bytes())
//
// ## 2. 确定性模式（新）
//
// 使用 [Deployer] 结构体，通过 masterSalt + path 字符串确定性派生 salt。
// 同样的 masterSalt + path 永远产生相同的 salt 和部署地址，无需额外存储。
//
// masterSalt 的生成建议：使用部署钱包的公钥哈希，这样只需要钱包就能复现所有 salt：
//
//	masterSalt := keccak256(publicKey)
//
// path 是任意字符串，用于区分不同的合约部署：
//
//	"project-a/token"
//	"project-a/whitelist"
//	"project-b/vault"
//
// 派生公式：salt[21:32] = keccak256(masterSalt || path)[:11]
//
// ### 构造函数
//
// 提供三个构造函数，对应 CreateX _guard 的不同模式：
//
//	// 零地址 + 跨链保护：任何人可部署，不同链不同地址
//	d := createx.NewDeployerZeroAddress(masterSalt, chainID)
//
//	// 零地址 + 无跨链保护：任何人可部署，所有链相同地址
//	d := createx.NewDeployerZeroAddressNoProtection(masterSalt)
//
//	// 指定发送者 + 跨链保护：仅 sender 可部署，不同链不同地址
//	d := createx.NewDeployerSender(myAddress, masterSalt, chainID)
//
// ### 使用示例
//
//	masterSalt := [32]byte(crypto.Keccak256(crypto.CompressPubkey(publicKey)))
//	d := createx.NewDeployerZeroAddress(masterSalt, chainID)
//
//	// 一步获取 calldata 和预计算地址
//	calldata1, whitelistAddr := d.Create2("project-a/whitelist", initCode1)
//	calldata2, tokenAddr := d.Create2("project-a/token", initCode2)
//
//	// CREATE3：地址不依赖 initCode
//	calldata3, vaultAddr := d.Create3("project-a/vault", initCode3)
//
//	// 只计算地址，不生成 calldata
//	addr := d.ComputeCreate3Address("project-a/vault")
//
// # CREATE2 vs CREATE3
//
// CREATE2 的部署地址取决于 salt + initCode，initCode 变化则地址变化。
//
// CREATE3 的部署地址仅取决于 salt，与 initCode 无关。
// 适合需要在部署前确定地址、且 initCode 可能变化的场景（如构造参数含其他合约地址）。
package createx
