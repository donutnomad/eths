package ethtype

import (
	"slices"

	"github.com/donutnomad/eths/ecommon"
)

// Block holds block header with full transaction objects and withdrawals.
type Block struct {
	Header
	Transactions []*Transaction `json:"transactions"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
}

func (b *Block) Transaction(hash ecommon.Hash) *Transaction {
	idx := slices.IndexFunc(b.Transactions, func(transaction *Transaction) bool {
		return transaction.Hash == hash
	})
	if idx < 0 {
		return nil
	}
	return b.Transactions[idx]
}
