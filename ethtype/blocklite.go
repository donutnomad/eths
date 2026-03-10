package ethtype

import (
	"encoding/json"

	"github.com/donutnomad/eths/ecommon"
)

// LiteBlock holds block header and transaction hashes without full tx data or uncles.
type LiteBlock struct {
	Header
	Transactions []ecommon.Hash `json:"transactions"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
}

func (b *LiteBlock) Transaction(hash ecommon.Hash) *ecommon.Hash {
	for _, transactionHash := range b.Transactions {
		if transactionHash == hash {
			return &transactionHash
		}
	}
	return nil
}

func (b *LiteBlock) MarshalJSON() ([]byte, error) {
	type enc struct {
		Header       Header         `json:",inline"`
		Transactions []ecommon.Hash `json:"transactions"`
		Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		Header:       b.Header,
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
	})
}

func (b *LiteBlock) UnmarshalJSON(data []byte) error {
	type raw struct {
		Transactions []ecommon.Hash `json:"transactions"`
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
