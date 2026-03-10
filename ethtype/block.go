package ethtype

import (
	"encoding/json"

	"github.com/donutnomad/eths/ecommon"
)

// Block holds block header with full transaction objects and withdrawals.
type Block struct {
	Header
	Transactions []*Transaction `json:"transactions"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
}

func (b *Block) Transaction(hash ecommon.Hash) *Transaction {
	for _, transaction := range b.Transactions {
		if transaction.Hash == hash {
			return transaction
		}
	}
	return nil
}

func (b Block) MarshalJSON() ([]byte, error) {
	type headerNoMethods Header
	type enc struct {
		headerNoMethods
		Transactions []*Transaction `json:"transactions"`
		Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		headerNoMethods: headerNoMethods(b.Header),
		Transactions:    b.Transactions,
		Withdrawals:     b.Withdrawals,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	type raw struct {
		Transactions []*Transaction `json:"transactions"`
		Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
	}
	var h Header
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	b.Header = h
	b.Transactions = r.Transactions
	b.Withdrawals = r.Withdrawals
	return nil
}
