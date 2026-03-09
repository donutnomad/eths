package ethtype

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

//go:generate go tool github.com/fjl/gencodec -type Tx -field-override txMarshaling -out gen_tx_json.go

type Tx struct {
	Type uint8 `json:"type"`

	ChainID              *big.Int                        `json:"chainId,omitempty"`
	Nonce                uint64                          `json:"nonce"`
	To                   *common.Address                 `json:"to"`
	Gas                  uint64                          `json:"gas"`
	GasPrice             *big.Int                        `json:"gasPrice"`
	MaxPriorityFeePerGas *big.Int                        `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *big.Int                        `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *big.Int                        `json:"maxFeePerBlobGas,omitempty"`
	Value                *big.Int                        `json:"value"`
	Input                []byte                          `json:"input"`
	AccessList           ethTypes.AccessList             `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash                   `json:"blobVersionedHashes,omitempty"`
	AuthorizationList    []ethTypes.SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                    *big.Int                        `json:"v"`
	R                    *big.Int                        `json:"r"`
	S                    *big.Int                        `json:"s"`
	YParity              uint64                          `json:"yParity,omitempty"`

	// Blob transaction sidecar encoding:
	Blobs       []kzg4844.Blob       `json:"blobs,omitempty"`
	Commitments []kzg4844.Commitment `json:"commitments,omitempty"`
	Proofs      []kzg4844.Proof      `json:"proofs,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

type txMarshaling struct {
	Type                 hexutil.Uint64
	ChainID              *hexutil.Big
	Nonce                hexutil.Uint64
	Gas                  hexutil.Uint64
	GasPrice             *hexutil.Big
	MaxPriorityFeePerGas *hexutil.Big
	MaxFeePerGas         *hexutil.Big
	MaxFeePerBlobGas     *hexutil.Big
	Value                *hexutil.Big
	Input                hexutil.Bytes
	V                    *hexutil.Big
	R                    *hexutil.Big
	S                    *hexutil.Big
	YParity              hexutil.Uint64
}
