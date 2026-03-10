package ethtype

import (
	"github.com/donutnomad/eths/ecommon"
	"github.com/donutnomad/eths/hexutil"
)

// Withdrawal represents a validator withdrawal from the consensus layer.
type Withdrawal struct {
	Index     uint64          `json:"index"`          // monotonically increasing identifier issued by consensus layer
	Validator uint64          `json:"validatorIndex"` // index of validator associated with withdrawal
	Address   ecommon.Address `json:"address"`        // target address for withdrawn ether
	Amount    uint64          `json:"amount"`         // value of withdrawal in Gwei
}

// field type overrides for gencodec
type withdrawalMarshaling struct {
	Index     hexutil.Uint64
	Validator hexutil.Uint64
	Amount    hexutil.Uint64
}
