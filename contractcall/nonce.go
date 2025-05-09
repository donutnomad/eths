package contractcall

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DefaultNonceManager struct {
	client *ethclient.Client
}

func NewDefaultNonceManager(client *ethclient.Client) *DefaultNonceManager {
	return &DefaultNonceManager{client: client}
}

func (n *DefaultNonceManager) GetNonce(ctx context.Context, account common.Address, isPending bool) (uint64, error) {
	if isPending {
		return n.client.PendingNonceAt(ctx, account)
	}
	return n.client.NonceAt(ctx, account, nil)
}

type StaticNonceManager struct {
	nonce uint64
}

func NewStaticNonceManager(nonce uint64) *StaticNonceManager {
	return &StaticNonceManager{nonce: nonce}
}

func (n *StaticNonceManager) GetNonce(ctx context.Context, account common.Address, isPending bool) (uint64, error) {
	return n.nonce, nil
}
