package ethtype

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
)

// Block holds block header with full transaction objects and withdrawals.
type Block struct {
	Header
	Transactions []*Transaction `json:"transactions"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
}

func (b *Block) MarshalJSON() ([]byte, error) {
	type enc struct {
		Header       Header         `json:",inline"`
		Transactions []*Transaction `json:"transactions"`
		Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		Header:       b.Header,
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
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

	a, _ := b.MarshalJSON()
	spew.Dump(string(a))
	return nil
}
