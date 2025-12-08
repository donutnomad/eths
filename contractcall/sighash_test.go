package contractcall

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

// TestSigHash_MatchesGethImplementation 验证自定义 SigHash 实现与 go-ethereum 标准库一致
func TestSigHash_MatchesGethImplementation(t *testing.T) {
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
		setup  func() (*txImpl, *ethTypes.Transaction)
	}{
		{
			name:   "Legacy (EIP-155)",
			txType: LegacyTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
				// 自定义实现
				impl := NewTxWith(LegacyTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGasPrice(gasPrice)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)

				// 标准库实现
				gethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
					Nonce:    nonce,
					GasPrice: gasPrice,
					Gas:      gas,
					To:       &to,
					Value:    value,
					Data:     data,
				})

				return impl, gethTx
			},
		},
		{
			name:   "AccessList (EIP-2930)",
			txType: AccessListTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}

				// 自定义实现
				impl := NewTxWith(AccessListTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGasPrice(gasPrice)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)
				impl.SetAccessList(accessList)

				// 标准库实现
				gethTx := ethTypes.NewTx(&ethTypes.AccessListTx{
					ChainID:    chainID,
					Nonce:      nonce,
					GasPrice:   gasPrice,
					Gas:        gas,
					To:         &to,
					Value:      value,
					Data:       data,
					AccessList: accessList,
				})

				return impl, gethTx
			},
		},
		{
			name:   "DynamicFee (EIP-1559)",
			txType: DynamicFeeTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
				accessList := ethTypes.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}

				// 自定义实现
				impl := NewTxWith(DynamicFeeTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)
				impl.SetMaxFeePerGas(maxFeePerGas)
				impl.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				impl.SetAccessList(accessList)

				// 标准库实现
				gethTx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
					ChainID:    chainID,
					Nonce:      nonce,
					Gas:        gas,
					To:         &to,
					Value:      value,
					Data:       data,
					GasTipCap:  maxPriorityFeePerGas,
					GasFeeCap:  maxFeePerGas,
					AccessList: accessList,
				})

				return impl, gethTx
			},
		},
		{
			name:   "Blob (EIP-4844)",
			txType: BlobTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
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

				// 自定义实现
				impl := NewTxWith(BlobTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)
				impl.SetMaxFeePerGas(maxFeePerGas)
				impl.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				impl.SetAccessList(accessList)
				impl.SetBlobHashes(blobHashes)
				impl.SetMaxFeePerBlobGas(maxFeePerBlobGas)

				// 标准库实现
				gethTx := ethTypes.NewTx(&ethTypes.BlobTx{
					ChainID:    uint256.MustFromBig(chainID),
					Nonce:      nonce,
					Gas:        gas,
					To:         to,
					Value:      uint256.MustFromBig(value),
					Data:       data,
					GasTipCap:  uint256.MustFromBig(maxPriorityFeePerGas),
					GasFeeCap:  uint256.MustFromBig(maxFeePerGas),
					AccessList: accessList,
					BlobHashes: blobHashes,
					BlobFeeCap: uint256.MustFromBig(maxFeePerBlobGas),
				})

				return impl, gethTx
			},
		},
		{
			name:   "SetCode (EIP-7702)",
			txType: SetCodeTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
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

				// 自定义实现
				impl := NewTxWith(SetCodeTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)
				impl.SetMaxFeePerGas(maxFeePerGas)
				impl.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)
				impl.SetAccessList(accessList)
				impl.SetAuthList(authList)

				// 标准库实现
				gethTx := ethTypes.NewTx(&ethTypes.SetCodeTx{
					ChainID:    uint256.MustFromBig(chainID),
					Nonce:      nonce,
					Gas:        gas,
					To:         to,
					Value:      uint256.MustFromBig(value),
					Data:       data,
					GasTipCap:  uint256.MustFromBig(maxPriorityFeePerGas),
					GasFeeCap:  uint256.MustFromBig(maxFeePerGas),
					AccessList: accessList,
					AuthList:   authList,
				})

				return impl, gethTx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl, gethTx := tt.setup()

			// 计算自定义实现的签名哈希
			customHash := impl.SigHash()

			// 计算标准库的签名哈希
			signer := ethTypes.NewPragueSigner(chainID)
			gethHash := signer.Hash(gethTx)

			// 对比两者
			if customHash != gethHash {
				t.Errorf("%s: SigHash mismatch\nCustom: %s\nGeth:   %s",
					tt.name, common.Hash(customHash).Hex(), gethHash.Hex())
			} else {
				t.Logf("%s: SigHash matched ✓ %s", tt.name, common.Hash(customHash).Hex())
			}
		})
	}
}

// TestSigHash_ContractCreation 测试合约创建交易（to 为 nil）
func TestSigHash_ContractCreation(t *testing.T) {
	chainID := big.NewInt(1)
	nonce := uint64(42)
	gas := uint64(500000)
	gasPrice := big.NewInt(1000000000)
	value := big.NewInt(0)
	data := common.FromHex("0x608060405234801561001057600080fd5b50")

	tests := []struct {
		name   string
		txType TxType
		setup  func() (*txImpl, *ethTypes.Transaction)
	}{
		{
			name:   "Legacy Contract Creation",
			txType: LegacyTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
				impl := NewTxWith(LegacyTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGasPrice(gasPrice)
				impl.SetGas(gas)
				impl.SetTo(nil) // 合约创建
				impl.SetValue(value)
				impl.SetData(data)

				gethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
					Nonce:    nonce,
					GasPrice: gasPrice,
					Gas:      gas,
					To:       nil, // 合约创建
					Value:    value,
					Data:     data,
				})

				return impl, gethTx
			},
		},
		{
			name:   "DynamicFee Contract Creation",
			txType: DynamicFeeTxType,
			setup: func() (*txImpl, *ethTypes.Transaction) {
				maxFeePerGas := big.NewInt(2000000000)
				maxPriorityFeePerGas := big.NewInt(1000000000)

				impl := NewTxWith(DynamicFeeTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGas(gas)
				impl.SetTo(nil) // 合约创建
				impl.SetValue(value)
				impl.SetData(data)
				impl.SetMaxFeePerGas(maxFeePerGas)
				impl.SetMaxPriorityFeePerGas(maxPriorityFeePerGas)

				gethTx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
					ChainID:   chainID,
					Nonce:     nonce,
					Gas:       gas,
					To:        nil, // 合约创建
					Value:     value,
					Data:      data,
					GasTipCap: maxPriorityFeePerGas,
					GasFeeCap: maxFeePerGas,
				})

				return impl, gethTx
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl, gethTx := tt.setup()

			customHash := impl.SigHash()
			signer := ethTypes.NewPragueSigner(chainID)
			gethHash := signer.Hash(gethTx)

			if customHash != gethHash {
				t.Errorf("%s: SigHash mismatch\nCustom: %s\nGeth:   %s",
					tt.name, common.Hash(customHash).Hex(), gethHash.Hex())
			} else {
				t.Logf("%s: SigHash matched ✓ %s", tt.name, common.Hash(customHash).Hex())
			}
		})
	}
}

// TestSigHash_DifferentChainIDs 测试不同的 chainID
func TestSigHash_DifferentChainIDs(t *testing.T) {
	chainIDs := []*big.Int{
		big.NewInt(1),        // 主网
		big.NewInt(5),        // Goerli
		big.NewInt(11155111), // Sepolia
		big.NewInt(137),      // Polygon
		big.NewInt(42161),    // Arbitrum
	}

	nonce := uint64(42)
	gas := uint64(21000)
	gasPrice := big.NewInt(1000000000)
	value := big.NewInt(1000000000000000000)
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	data := []byte{}

	for _, chainID := range chainIDs {
		t.Run("ChainID_"+chainID.String(), func(t *testing.T) {
			// 自定义实现
			impl := NewTxWith(LegacyTxType, chainID).(*txImpl)
			impl.SetNonce(nonce)
			impl.SetGasPrice(gasPrice)
			impl.SetGas(gas)
			impl.SetTo(&to)
			impl.SetValue(value)
			impl.SetData(data)

			// 标准库实现
			gethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				Nonce:    nonce,
				GasPrice: gasPrice,
				Gas:      gas,
				To:       &to,
				Value:    value,
				Data:     data,
			})

			customHash := impl.SigHash()
			signer := ethTypes.NewPragueSigner(chainID)
			gethHash := signer.Hash(gethTx)

			if customHash != gethHash {
				t.Errorf("ChainID %s: SigHash mismatch\nCustom: %s\nGeth:   %s",
					chainID.String(), common.Hash(customHash).Hex(), gethHash.Hex())
			} else {
				t.Logf("ChainID %s: SigHash matched ✓", chainID.String())
			}
		})
	}
}
