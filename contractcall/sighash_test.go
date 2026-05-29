package contractcall

import (
	"math/big"
	"testing"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethtype"
	"github.com/holiman/uint256"
)

// Fixtures generated once from go-ethereum v1.17.3 PragueSigner.Hash.
var expectedSigHashes = map[string]common.Hash{
	"Legacy (EIP-155)":             common.HexToHash("0x2beba12561f37f3a200202cc73876d088714f1d500b55be090db6cbe8e24fc81"),
	"AccessList (EIP-2930)":        common.HexToHash("0x2becc9bfdc52069fda74d7a3244f71c96daddbb4a603118dfad92d68d7e179d8"),
	"DynamicFee (EIP-1559)":        common.HexToHash("0x0960dd38498e3217c5cf2d1200fe47ebcae847c18e8bc8f402c0da34b358f5f2"),
	"Blob (EIP-4844)":              common.HexToHash("0x5e61636361106f0d150ed1c47ff644deeb8f1957ed4a164f6e99cef11ded186d"),
	"SetCode (EIP-7702)":           common.HexToHash("0x163524a1021c78bc229e76073d4cb7cf868be9c5d7ed11023d54a1ff6fbbf65f"),
	"Legacy Contract Creation":     common.HexToHash("0x42b39edde75a12ee7de4bbfd9407fe863ed09aeaedbdca8b704268e681e2d8e1"),
	"DynamicFee Contract Creation": common.HexToHash("0xd4e4f2ffd8a225f52eca1712a1457b3d3d023442d0e36be815002474ac20e13e"),
}

var expectedSigHashesByChainID = map[string]common.Hash{
	"1":        common.HexToHash("0x916aad7ad38d0296216a3f5678776d652defe97d536ac68da843062cc46159ef"),
	"5":        common.HexToHash("0xc27d8e9c1a8e8c51b9c7d4d37a6b24908031999d753c5151064970b233065155"),
	"11155111": common.HexToHash("0x49cb107caea1dcfb1578c51617ed70a9acc3391d2387e45648b8fb9bc88ad848"),
	"137":      common.HexToHash("0x62e08233cc7f9bda4d6be21fd8a944ff866aa81aaec57e97061dcd46497a5a77"),
	"42161":    common.HexToHash("0x55ebc27846c7d2d1d4eb198167a45efb50d9c15f296f7f2366673e421a0187ea"),
}

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
		setup  func() *txImpl
	}{
		{
			name:   "Legacy (EIP-155)",
			txType: LegacyTxType,
			setup: func() *txImpl {
				// 自定义实现
				impl := NewTxWith(LegacyTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGasPrice(gasPrice)
				impl.SetGas(gas)
				impl.SetTo(&to)
				impl.SetValue(value)
				impl.SetData(data)

				return impl
			},
		},
		{
			name:   "AccessList (EIP-2930)",
			txType: AccessListTxType,
			setup: func() *txImpl {
				accessList := ethtype.AccessList{
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

				return impl
			},
		},
		{
			name:   "DynamicFee (EIP-1559)",
			txType: DynamicFeeTxType,
			setup: func() *txImpl {
				accessList := ethtype.AccessList{
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

				return impl
			},
		},
		{
			name:   "Blob (EIP-4844)",
			txType: BlobTxType,
			setup: func() *txImpl {
				accessList := ethtype.AccessList{
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

				return impl
			},
		},
		{
			name:   "SetCode (EIP-7702)",
			txType: SetCodeTxType,
			setup: func() *txImpl {
				accessList := ethtype.AccessList{
					{
						Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
						StorageKeys: []common.Hash{
							common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
						},
					},
				}
				authList := []ethtype.SetCodeAuthorization{
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

				return impl
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := tt.setup()

			// 计算自定义实现的签名哈希
			customHash := impl.SigHash()

			expectedHash, ok := expectedSigHashes[tt.name]
			if !ok {
				t.Fatalf("%s: missing expected SigHash", tt.name)
			}

			// 对比两者
			if customHash != expectedHash {
				t.Errorf("%s: SigHash mismatch\nCustom:   %s\nExpected: %s",
					tt.name, common.Hash(customHash).Hex(), expectedHash.Hex())
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
		setup  func() *txImpl
	}{
		{
			name:   "Legacy Contract Creation",
			txType: LegacyTxType,
			setup: func() *txImpl {
				impl := NewTxWith(LegacyTxType, chainID).(*txImpl)
				impl.SetNonce(nonce)
				impl.SetGasPrice(gasPrice)
				impl.SetGas(gas)
				impl.SetTo(nil) // 合约创建
				impl.SetValue(value)
				impl.SetData(data)

				return impl
			},
		},
		{
			name:   "DynamicFee Contract Creation",
			txType: DynamicFeeTxType,
			setup: func() *txImpl {
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

				return impl
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := tt.setup()

			customHash := impl.SigHash()
			expectedHash, ok := expectedSigHashes[tt.name]
			if !ok {
				t.Fatalf("%s: missing expected SigHash", tt.name)
			}

			if customHash != expectedHash {
				t.Errorf("%s: SigHash mismatch\nCustom:   %s\nExpected: %s",
					tt.name, common.Hash(customHash).Hex(), expectedHash.Hex())
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

			customHash := impl.SigHash()
			expectedHash, ok := expectedSigHashesByChainID[chainID.String()]
			if !ok {
				t.Fatalf("ChainID %s: missing expected SigHash", chainID.String())
			}

			if customHash != expectedHash {
				t.Errorf("ChainID %s: SigHash mismatch\nCustom:   %s\nExpected: %s",
					chainID.String(), common.Hash(customHash).Hex(), expectedHash.Hex())
			} else {
				t.Logf("ChainID %s: SigHash matched ✓", chainID.String())
			}
		})
	}
}
