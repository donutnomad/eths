package multiread

import (
	"context"
	"math/big"
	"testing"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethclient"
	"github.com/donutnomad/eths/ethtype"
)

// mockContractCaller implements bind.ContractCaller for testing.
type mockContractCaller struct{}

func (m *mockContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}
func (m *mockContractCaller) CallContract(ctx context.Context, call ethclient.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}
func (m *mockContractCaller) HeaderByNumber(ctx context.Context, number *big.Int) (*ethtype.EHeader, error) {
	return nil, nil
}

// mockChainIDCaller implements both bind.ContractCaller and ethclient.ChainIDReader.
type mockChainIDCaller struct {
	mockContractCaller
	chainID *big.Int
}

func (m *mockChainIDCaller) ChainID(ctx context.Context) (*big.Int, error) {
	return m.chainID, nil
}

func TestRegisterAndGetAddress(t *testing.T) {
	customAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	var chainID uint64 = 31337

	RegisterAddress(chainID, customAddr)
	t.Cleanup(func() {
		addressMu.Lock()
		delete(addressMap, chainID)
		addressMu.Unlock()
	})

	got := GetAddress(chainID)
	if got != customAddr {
		t.Errorf("GetAddress(%d) = %s, want %s", chainID, got.Hex(), customAddr.Hex())
	}
}

func TestGetAddress_FallbackDefault(t *testing.T) {
	got := GetAddress(99999)
	if got != Address {
		t.Errorf("GetAddress(99999) = %s, want default %s", got.Hex(), Address.Hex())
	}
}

func TestGetAddress_WithChainIDReader(t *testing.T) {
	customAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	var chainID uint64 = 42161

	RegisterAddress(chainID, customAddr)
	t.Cleanup(func() {
		addressMu.Lock()
		delete(addressMap, chainID)
		addressMu.Unlock()
	})

	client := &mockChainIDCaller{chainID: new(big.Int).SetUint64(chainID)}
	got := getAddress(client)
	if got != customAddr {
		t.Errorf("getAddress() = %s, want %s", got.Hex(), customAddr.Hex())
	}
}

func TestGetAddress_WithoutChainIDReader(t *testing.T) {
	client := &mockContractCaller{}
	got := getAddress(client)
	if got != Address {
		t.Errorf("getAddress() = %s, want default %s", got.Hex(), Address.Hex())
	}
}

func TestGetAddress_UnregisteredChainID(t *testing.T) {
	client := &mockChainIDCaller{chainID: big.NewInt(88888)}
	got := getAddress(client)
	if got != Address {
		t.Errorf("getAddress() = %s, want default %s", got.Hex(), Address.Hex())
	}
}

func TestRegisterAddress_Override(t *testing.T) {
	var chainID uint64 = 10
	addr1 := common.HexToAddress("0x3333333333333333333333333333333333333333")
	addr2 := common.HexToAddress("0x4444444444444444444444444444444444444444")

	RegisterAddress(chainID, addr1)
	RegisterAddress(chainID, addr2)
	t.Cleanup(func() {
		addressMu.Lock()
		delete(addressMap, chainID)
		addressMu.Unlock()
	})

	got := GetAddress(chainID)
	if got != addr2 {
		t.Errorf("GetAddress(%d) = %s, want overridden %s", chainID, got.Hex(), addr2.Hex())
	}
}
