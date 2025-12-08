package contractcall

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/samber/lo"
)

// txJSON is the JSON representation of transactions.
type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID              *hexutil.Big                    `json:"chainId,omitempty"`
	Nonce                *hexutil.Uint64                 `json:"nonce"`
	To                   *common.Address                 `json:"to"`
	Gas                  *hexutil.Uint64                 `json:"gas"`
	GasPrice             *hexutil.Big                    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big                    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big                    `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *hexutil.Big                    `json:"maxFeePerBlobGas,omitempty"`
	Value                *hexutil.Big                    `json:"value"`
	Input                *hexutil.Bytes                  `json:"input"`
	AccessList           *ethTypes.AccessList            `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash                   `json:"blobVersionedHashes,omitempty"`
	AuthorizationList    []ethTypes.SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                    *hexutil.Big                    `json:"v"`
	R                    *hexutil.Big                    `json:"r"`
	S                    *hexutil.Big                    `json:"s"`
	YParity              *hexutil.Uint64                 `json:"yParity,omitempty"`

	// Blob transaction sidecar encoding:
	Blobs       []kzg4844.Blob       `json:"blobs,omitempty"`
	Commitments []kzg4844.Commitment `json:"commitments,omitempty"`
	Proofs      []kzg4844.Proof      `json:"proofs,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

func (t *txImpl) isProtected() bool {
	if t.txType != LegacyTxType {
		return false
	}
	v := t.getV()
	if v.BitLen() <= 8 {
		v_ := v.Uint64()
		return v_ != 27 && v_ != 28 && v_ != 1 && v_ != 0
	} else {
		// anything not 27 or 28 is considered protected
		return true
	}
}

// MarshalJSON marshals as JSON with a hash.
func (t *txImpl) MarshalJSON() ([]byte, error) {
	var enc = txJSON{
		Hash:    t.Hash(),
		Type:    hexutil.Uint64(t.TxType()),
		ChainID: (*hexutil.Big)(t.ChainID()),
		Nonce:   (*hexutil.Uint64)(lo.ToPtr(t.Nonce())),
		To:      t.To(),
		Gas:     (*hexutil.Uint64)(lo.ToPtr(t.Gas())),
		Value:   (*hexutil.Big)(t.Value()),
		Input:   (*hexutil.Bytes)(lo.ToPtr(t.Data())),
		V:       (*hexutil.Big)(t.getV().ToBig()),
		R:       (*hexutil.Big)(t.getR().ToBig()),
		S:       (*hexutil.Big)(t.getS().ToBig()),
	}
	if t.TxType().IsEIP1559Gas() {
		enc.MaxFeePerGas = (*hexutil.Big)(t.MaxFeePerGas())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(t.MaxPriorityFeePerGas())
	} else {
		enc.GasPrice = (*hexutil.Big)(t.GasPrice())
	}

	if t.isModern() {
		enc.AccessList = lo.ToPtr(t.AccessList())
		enc.YParity = (*hexutil.Uint64)(lo.ToPtr(t.getV().Uint64()))
	}

	switch t.TxType() {
	case BlobTxType:
		enc.MaxFeePerBlobGas = (*hexutil.Big)(t.MaxFeePerBlobGas())
		enc.BlobVersionedHashes = t.BlobHashes()
		if sidecar := t.Sidecar(); sidecar != nil {
			enc.Blobs = t.Sidecar().Blobs
			enc.Commitments = t.Sidecar().Commitments
			enc.Proofs = t.Sidecar().Proofs
		}
	case SetCodeTxType:
		enc.AuthorizationList = t.AuthList()
	}

	return json.Marshal(&enc)
}
