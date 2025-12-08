package contractcall

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

// TestHash_ConsistencyWithHash2 验证 Hash() 和 Hash2() 对所有交易类型返回一致的结果
func TestHash_ConsistencyWithHash2(t *testing.T) {
	// 生成测试私钥和签名器
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}

	// 创建一个签名函数
	signFunc := func(tx *txImpl) error {
		sigHash := tx.SigHash()
		sig, err := crypto.Sign(sigHash[:], privateKey)
		if err != nil {
			return err
		}

		// 设置签名值
		v := sig[64]
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:64])

		tx.SetSignature(
			computeVForEIP155(v, tx.ChainID(), tx.txType == LegacyTxType),
			r, s,
		)
		return nil
	}

	chainID := big.NewInt(1)
	nonce := uint64(42)
	gas := uint64(21000)
	gasPrice := big.NewInt(1000000000)
	maxFeePerGas := big.NewInt(2000000000)
	maxPriorityFeePerGas := big.NewInt(1000000000)
	value := big.NewInt(1000000000000000000)
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	data := []byte("test data")

	tests := []struct {
		name   string
		txType TxType
		setup  func() *txImpl
	}{
		{
			name:   "Legacy (EIP-155)",
			txType: LegacyTxType,
			setup: func() *txImpl {
				tx := NewTxWith(LegacyTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGasPrice(gasPrice)
				tx.SetGas(gas)
				tx.SetTo(&to)
				tx.SetValue(value)
				tx.SetData(data)
				return tx
			},
		},
		{
			name:   "Legacy Contract Creation",
			txType: LegacyTxType,
			setup: func() *txImpl {
				tx := NewTxWith(LegacyTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGasPrice(gasPrice)
				tx.SetGas(gas)
				tx.SetTo(nil) // 合约创建
				tx.SetValue(value)
				tx.SetData(data)
				return tx
			},
		},
		{
			name:   "AccessList (EIP-2930)",
			txType: AccessListTxType,
			setup: func() *txImpl {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}

				tx := NewTxWith(AccessListTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGasPrice(gasPrice)
				tx.SetGas(gas)
				tx.SetTo(&to)
				tx.SetValue(value)
				tx.SetData(data)
				tx.SetAccessList(accessList)
				return tx
			},
		},
		{
			name:   "AccessList Contract Creation",
			txType: AccessListTxType,
			setup: func() *txImpl {
				tx := NewTxWith(AccessListTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGasPrice(gasPrice)
				tx.SetGas(gas)
				tx.SetTo(nil) // 合约创建
				tx.SetValue(value)
				tx.SetData(data)
				return tx
			},
		},
		{
			name:   "DynamicFee (EIP-1559)",
			txType: DynamicFeeTxType,
			setup: func() *txImpl {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}

				tx := NewTxWith(DynamicFeeTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGas(gas)
				tx.SetTo(&to)
				tx.SetValue(value)
				tx.SetData(data)
				tx.SetMaxFeePerGas(maxFeePerGas)
				tx.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				tx.SetAccessList(accessList)
				return tx
			},
		},
		{
			name:   "DynamicFee Contract Creation",
			txType: DynamicFeeTxType,
			setup: func() *txImpl {
				tx := NewTxWith(DynamicFeeTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGas(gas)
				tx.SetTo(nil) // 合约创建
				tx.SetValue(value)
				tx.SetData(data)
				tx.SetMaxFeePerGas(maxFeePerGas)
				tx.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				return tx
			},
		},
		{
			name:   "Blob (EIP-4844)",
			txType: BlobTxType,
			setup: func() *txImpl {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}
				blobHashes := []common.Hash{
					common.HexToHash("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"),
				}
				maxFeePerBlobGas := big.NewInt(1000000)

				tx := NewTxWith(BlobTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGas(gas)
				tx.SetTo(&to)
				tx.SetValue(value)
				tx.SetData(data)
				tx.SetMaxFeePerGas(maxFeePerGas)
				tx.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				tx.SetAccessList(accessList)
				tx.SetBlobHashes(blobHashes)
				tx.SetMaxFeePerBlobGas(maxFeePerBlobGas)
				return tx
			},
		},
		{
			name:   "SetCode (EIP-7702)",
			txType: SetCodeTxType,
			setup: func() *txImpl {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}
				authList := []ethTypes.SetCodeAuthorization{
					{
						ChainID: *uint256.NewInt(1),
						Address: common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
						Nonce:   10,
						V:       27,
						R:       *uint256.NewInt(1),
						S:       *uint256.NewInt(2),
					},
				}

				tx := NewTxWith(SetCodeTxType, chainID).(*txImpl)
				tx.SetNonce(nonce)
				tx.SetGas(gas)
				tx.SetTo(&to)
				tx.SetValue(value)
				tx.SetData(data)
				tx.SetMaxFeePerGas(maxFeePerGas)
				tx.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				tx.SetAccessList(accessList)
				tx.SetAuthList(authList)
				return tx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建并签名交易
			tx := tt.setup()
			err := signFunc(tx)
			if err != nil {
				t.Fatalf("签名失败: %v", err)
			}

			// 验证签名非零
			_, r, s := tx.Signature()
			if r.Sign() == 0 || s.Sign() == 0 {
				t.Fatal("签名值为零")
			}

			// 计算 Hash() 和 Hash2()
			hash1 := tx.ToTransaction().Hash()
			hash2 := tx.Hash()

			// 验证两者一致
			if hash1 != hash2 {
				t.Errorf("%s: Hash 不一致\nHash():  %s\nHash2(): %s",
					tt.name, hash1.Hex(), hash2.Hex())
			} else {
				t.Logf("%s: Hash 一致 ✓ %s", tt.name, hash1.Hex())
			}

			// 额外验证：与 go-ethereum 库的哈希一致
			gethTx := tx.ToTransaction()
			gethHash := gethTx.Hash()

			if hash1 != gethHash {
				t.Errorf("%s: 与 go-ethereum 哈希不一致\nOurs: %s\nGeth: %s",
					tt.name, hash1.Hex(), gethHash.Hex())
			}
		})
	}
}

// TestHash_UnsignedTransaction 测试未签名交易的 Hash 和 Hash2
func TestHash_UnsignedTransaction(t *testing.T) {
	chainID := big.NewInt(1)
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	tests := []struct {
		name   string
		txType TxType
		setup  func() *txImpl
	}{
		{
			name:   "Unsigned Legacy",
			txType: LegacyTxType,
			setup: func() *txImpl {
				tx := NewTxWith(LegacyTxType, chainID).(*txImpl)
				tx.SetNonce(1)
				tx.SetGasPrice(big.NewInt(1000000000))
				tx.SetGas(21000)
				tx.SetTo(&to)
				tx.SetValue(big.NewInt(1000))
				return tx
			},
		},
		{
			name:   "Unsigned DynamicFee",
			txType: DynamicFeeTxType,
			setup: func() *txImpl {
				tx := NewTxWith(DynamicFeeTxType, chainID).(*txImpl)
				tx.SetNonce(1)
				tx.SetGas(21000)
				tx.SetTo(&to)
				tx.SetValue(big.NewInt(1000))
				tx.SetMaxFeePerGas(big.NewInt(2000000000))
				tx.SetMaxPriorityFeePerGas(big.NewInt(1000000000))
				return tx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.setup()

			// 未签名时，V、R、S 应该为 0
			v, r, s := tx.Signature()
			if v.Sign() != 0 || r.Sign() != 0 || s.Sign() != 0 {
				t.Log("警告: 未签名交易的签名值应该为 0")
			}

			hash1 := tx.ToTransaction().Hash()
			hash2 := tx.Hash()

			if hash1 != hash2 {
				t.Errorf("%s: 未签名交易的 Hash 不一致\nHash():  %s\nHash2(): %s",
					tt.name, hash1.Hex(), hash2.Hex())
			} else {
				t.Logf("%s: 未签名交易 Hash 一致 ✓ %s", tt.name, hash1.Hex())
			}
		})
	}
}

// TestHash_DifferentChainIDs 测试不同 chainID 的交易
func TestHash_DifferentChainIDs(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("生成私钥失败: %v", err)
	}

	// 创建签名函数
	signFunc := func(tx *txImpl) error {
		sigHash := tx.SigHash()
		sig, err := crypto.Sign(sigHash[:], privateKey)
		if err != nil {
			return err
		}

		v := sig[64]
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:64])

		tx.SetSignature(
			computeVForEIP155(v, tx.ChainID(), tx.txType == LegacyTxType),
			r, s,
		)
		return nil
	}

	chainIDs := []*big.Int{
		big.NewInt(1),        // Mainnet
		big.NewInt(5),        // Goerli
		big.NewInt(11155111), // Sepolia
		big.NewInt(137),      // Polygon
		big.NewInt(42161),    // Arbitrum
	}

	to := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	for _, chainID := range chainIDs {
		t.Run("ChainID_"+chainID.String(), func(t *testing.T) {
			// 测试 Legacy 交易
			tx := NewTxWith(LegacyTxType, chainID).(*txImpl)
			tx.SetNonce(1)
			tx.SetGasPrice(big.NewInt(1000000000))
			tx.SetGas(21000)
			tx.SetTo(&to)
			tx.SetValue(big.NewInt(1000))

			err := signFunc(tx)
			if err != nil {
				t.Fatalf("签名失败: %v", err)
			}

			hash1 := tx.ToTransaction().Hash()
			hash2 := tx.Hash()

			if hash1 != hash2 {
				t.Errorf("ChainID %s: Hash 不一致\nHash():  %s\nHash2(): %s",
					chainID.String(), hash1.Hex(), hash2.Hex())
			} else {
				t.Logf("ChainID %s: Hash 一致 ✓", chainID.String())
			}

			// 验证与 geth 一致
			gethHash := tx.ToTransaction().Hash()
			if hash1 != gethHash {
				t.Errorf("ChainID %s: 与 geth 不一致\nOurs: %s\nGeth: %s",
					chainID.String(), hash1.Hex(), gethHash.Hex())
			}
		})
	}
}
