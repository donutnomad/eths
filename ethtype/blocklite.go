package ethtype

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// LiteBlock holds block header and transaction hashes without full tx data or uncles.
type LiteBlock struct {
	Header
	Transactions []common.Hash `json:"transactions"`
	Withdrawals  []*Withdrawal `json:"withdrawals,omitempty"`
}

func (b *LiteBlock) MarshalJSON() ([]byte, error) {
	type enc struct {
		Header       Header        `json:",inline"`
		Transactions []common.Hash `json:"transactions"`
		Withdrawals  []*Withdrawal `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		Header:       b.Header,
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
	})
}

func (b *LiteBlock) UnmarshalJSON(data []byte) error {
	type raw struct {
		Header       Header        `json:",inline"`
		Transactions []common.Hash `json:"transactions"`
		Withdrawals  []*Withdrawal `json:"withdrawals,omitempty"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	b.Header = r.Header
	b.Transactions = r.Transactions
	b.Withdrawals = r.Withdrawals
	return nil
}
