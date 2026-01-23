package contractcall

import (
	"math/big"
	"testing"

	"github.com/donutnomad/eths/contracts_pack"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestParseERC20Errors(t *testing.T) {
	erc20ABI := lo.Must1(contracts_pack.ERC20MetaData.ParseABI())
	knownABIs := []*abi.ABI{erc20ABI}

	tests := []struct {
		name     string
		errData  string
		wantName string
		wantSig  string
		wantArgs []string
	}{
		{
			// ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
			// selector: 0xe450d38c
			name:     "ERC20InsufficientBalance",
			errData:  encodeERC20Error(erc20ABI, "ERC20InsufficientBalance", common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"), big100, big1000),
			wantName: "ERC20InsufficientBalance",
			wantSig:  "ERC20InsufficientBalance(address,uint256,uint256)",
			wantArgs: []string{"0x1234567890AbcdEF1234567890aBcdef12345678", "100", "1000"},
		},
		{
			// ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
			// selector: 0xfb8f41b2
			name:     "ERC20InsufficientAllowance",
			errData:  encodeERC20Error(erc20ABI, "ERC20InsufficientAllowance", common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"), big50, big200),
			wantName: "ERC20InsufficientAllowance",
			wantSig:  "ERC20InsufficientAllowance(address,uint256,uint256)",
			wantArgs: []string{"0xabCDEF1234567890ABcDEF1234567890aBCDeF12", "50", "200"},
		},
		{
			// ERC20InvalidApprover(address approver)
			// selector: 0xe602df05
			name:     "ERC20InvalidApprover",
			errData:  encodeERC20Error(erc20ABI, "ERC20InvalidApprover", common.HexToAddress("0x0000000000000000000000000000000000000000")),
			wantName: "ERC20InvalidApprover",
			wantSig:  "ERC20InvalidApprover(address)",
			wantArgs: []string{"0x0000000000000000000000000000000000000000"},
		},
		{
			// ERC20InvalidReceiver(address receiver)
			// selector: 0xec442f05
			name:     "ERC20InvalidReceiver",
			errData:  encodeERC20Error(erc20ABI, "ERC20InvalidReceiver", common.HexToAddress("0x0000000000000000000000000000000000000000")),
			wantName: "ERC20InvalidReceiver",
			wantSig:  "ERC20InvalidReceiver(address)",
			wantArgs: []string{"0x0000000000000000000000000000000000000000"},
		},
		{
			// ERC20InvalidSender(address sender)
			// selector: 0x96c6fd1e
			name:     "ERC20InvalidSender",
			errData:  encodeERC20Error(erc20ABI, "ERC20InvalidSender", common.HexToAddress("0x0000000000000000000000000000000000000000")),
			wantName: "ERC20InvalidSender",
			wantSig:  "ERC20InvalidSender(address)",
			wantArgs: []string{"0x0000000000000000000000000000000000000000"},
		},
		{
			// ERC20InvalidSpender(address spender)
			// selector: 0x94280d62
			name:     "ERC20InvalidSpender",
			errData:  encodeERC20Error(erc20ABI, "ERC20InvalidSpender", common.HexToAddress("0x0000000000000000000000000000000000000000")),
			wantName: "ERC20InvalidSpender",
			wantSig:  "ERC20InvalidSpender(address)",
			wantArgs: []string{"0x0000000000000000000000000000000000000000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseContractError(knownABIs, &EvmError{
				ErrData: tt.errData,
			})

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantName, result.Name)
			assert.Equal(t, tt.wantSig, result.Signature)
			assert.Equal(t, tt.wantArgs, result.Arguments)

			t.Logf("Name: %s", result.Name)
			t.Logf("Signature: %s", result.Signature)
			t.Logf("Definition: %s", result.Definition)
			t.Logf("Formatted: %s", result.Formatted)
		})
	}
}

func TestParseStandardErrors(t *testing.T) {
	tests := []struct {
		name     string
		errData  string
		wantName string
		wantArgs []string
	}{
		{
			// Error(string) - standard revert with message "Not enough balance"
			// selector: 0x08c379a0
			// The data is: selector(4) + offset(32) + length(32) + string_data(32) = 100 bytes
			name:     "Error(string)",
			errData:  "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001254657374206572726f72206d6573736167650000000000000000000000000000",
			wantName: "Error",
			wantArgs: []string{"Test error message"},
		},
		{
			// Panic(uint256) - arithmetic overflow (code 0x11)
			// selector: 0x4e487b71
			name:     "Panic(uint256) - overflow",
			errData:  "0x4e487b710000000000000000000000000000000000000000000000000000000000000011",
			wantName: "Panic",
			wantArgs: []string{"17"}, // 0x11 = 17 (arithmetic overflow)
		},
		{
			// Panic(uint256) - division by zero (code 0x12)
			name:     "Panic(uint256) - division by zero",
			errData:  "0x4e487b710000000000000000000000000000000000000000000000000000000000000012",
			wantName: "Panic",
			wantArgs: []string{"18"}, // 0x12 = 18 (division by zero)
		},
		{
			// Panic(uint256) - assert failure (code 0x01)
			name:     "Panic(uint256) - assert failure",
			errData:  "0x4e487b710000000000000000000000000000000000000000000000000000000000000001",
			wantName: "Panic",
			wantArgs: []string{"1"}, // 0x01 = 1 (assert failure)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseContractError(nil, &EvmError{
				ErrData: tt.errData,
			})

			if !assert.NotNil(t, result) {
				return
			}
			assert.Equal(t, tt.wantName, result.Name)
			assert.Equal(t, tt.wantArgs, result.Arguments)

			t.Logf("Name: %s", result.Name)
			t.Logf("Definition: %s", result.Definition)
			t.Logf("Formatted: %s", result.Formatted)
		})
	}
}

// Helper variables for test data
var (
	big50   = big.NewInt(50)
	big100  = big.NewInt(100)
	big200  = big.NewInt(200)
	big1000 = big.NewInt(1000)
)

// encodeERC20Error encodes an ERC20 error with the given name and arguments
func encodeERC20Error(abiDef *abi.ABI, errName string, args ...any) string {
	errDef, ok := abiDef.Errors[errName]
	if !ok {
		panic("error not found: " + errName)
	}
	packed, err := errDef.Inputs.Pack(args...)
	if err != nil {
		panic(err)
	}
	return "0x" + common.Bytes2Hex(append(errDef.ID[:4], packed...))
}
