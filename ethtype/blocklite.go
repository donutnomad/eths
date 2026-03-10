package ethtype

import (
	"slices"

	"github.com/donutnomad/eths/ecommon"
)

// LiteBlock holds block header and transaction hashes without full tx data or uncles.
type LiteBlock struct {
	Header
	Transactions []ecommon.Hash `json:"transactions"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
}

func (b *LiteBlock) Transaction(hash ecommon.Hash) *ecommon.Hash {
	idx := slices.Index(b.Transactions, hash)
	if idx < 0 {
		return nil
	}
	return &b.Transactions[idx]
}
