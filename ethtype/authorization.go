package ethtype

import (
	"github.com/donutnomad/eths/ecommon"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

//go:generate go run github.com/fjl/gencodec -type SetCodeAuthorization -field-override authorizationMarshaling -out authorization_generated.go

// SetCodeAuthorization is an authorization from an account to deploy code at its address.
type SetCodeAuthorization struct {
	ChainID uint256.Int     `json:"chainId"`
	Address ecommon.Address `json:"address"`
	Nonce   uint64          `json:"nonce"`
	V       uint8           `json:"yParity"`
	R       uint256.Int     `json:"r"`
	S       uint256.Int     `json:"s"`
}

// field type overrides for gencodec
type authorizationMarshaling struct {
	ChainID hexutil.U256
	Nonce   hexutil.Uint64
	V       hexutil.Uint64
	R       hexutil.U256
	S       hexutil.U256
}
