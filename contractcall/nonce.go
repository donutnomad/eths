package contractcall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type INonceAt interface {
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

type ICodeAt interface {
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

type DefaultNonceManager struct {
	client INonceAt
}

func NewDefaultNonceManager(client INonceAt) *DefaultNonceManager {
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
