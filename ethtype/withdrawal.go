package ethtype

import (
	"github.com/donutnomad/eths/ecommon"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate go run github.com/fjl/gencodec -type Withdrawal -field-override withdrawalMarshaling -out withdrawal_generated.go

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
