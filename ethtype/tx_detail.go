package ethtype

import (
	"encoding/json"
	"math/big"

	"github.com/donutnomad/eths/ecommon"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TxDetail combines transaction fields with receipt fields,
// embedding Tx and adding non-overlapping receipt fields.
type TxDetail struct {
	Tx

	// https://ethereum.org/developers/docs/apis/json-rpc/#eth_gettransactionreceipt

	// Receipt consensus fields
	PostState         []byte `json:"root"`   // pre Byzantium
	Status            uint64 `json:"status"` // after Byzantium
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	Bloom             Bloom  `json:"logsBloom"`
	Logs              []*Log `json:"logs"`

	// Receipt implementation fields
	ContractAddress   *ecommon.Address `json:"contractAddress"`
	From              ecommon.Address  `json:"from"`
	GasUsed           uint64           `json:"gasUsed"`
	EffectiveGasPrice *big.Int         `json:"effectiveGasPrice"`
	BlobGasUsed       uint64           `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int         `json:"blobGasPrice,omitempty"`

	// Receipt inclusion fields
	BlockHash        ecommon.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int     `json:"blockNumber,omitempty"`
	TransactionIndex uint         `json:"transactionIndex"`
}

// txDetailJSON is the JSON representation used for marshaling/unmarshaling.
type txDetailJSON struct {
	// Tx fields
	Type                 hexutil.Uint64         `json:"type"`
	ChainID              *hexutil.Big           `json:"chainId,omitempty"`
	Nonce                hexutil.Uint64         `json:"nonce"`
	To                   *ecommon.Address       `json:"to"`
	Gas                  hexutil.Uint64         `json:"gas"`
	GasPrice             *hexutil.Big           `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big           `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big           `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *hexutil.Big           `json:"maxFeePerBlobGas,omitempty"`
	Value                *hexutil.Big           `json:"value"`
	Input                hexutil.Bytes          `json:"input"`
	AccessList           AccessList             `json:"accessList,omitempty"`
	BlobVersionedHashes  []ecommon.Hash         `json:"blobVersionedHashes,omitempty"`
	AuthorizationList    []SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                    *hexutil.Big           `json:"v"`
	R                    *hexutil.Big           `json:"r"`
	S                    *hexutil.Big           `json:"s"`
	YParity              hexutil.Uint64         `json:"yParity,omitempty"`
	Hash                 ecommon.Hash           `json:"hash"`

	// https://ethereum.org/developers/docs/apis/json-rpc/#eth_gettransactionreceipt

	// Receipt fields
	PostState hexutil.Bytes  `json:"root"`   // pre Byzantium
	Status    hexutil.Uint64 `json:"status"` // after Byzantium

	CumulativeGasUsed hexutil.Uint64   `json:"cumulativeGasUsed"`
	Bloom             Bloom            `json:"logsBloom"`
	Logs              []*Log           `json:"logs"`
	ContractAddress   *ecommon.Address `json:"contractAddress"`
	From              ecommon.Address  `json:"from"`
	GasUsed           hexutil.Uint64   `json:"gasUsed"`
	EffectiveGasPrice *hexutil.Big     `json:"effectiveGasPrice"`
	BlobGasUsed       hexutil.Uint64   `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *hexutil.Big     `json:"blobGasPrice,omitempty"`

	BlockHash        ecommon.Hash `json:"blockHash,omitempty"`
	BlockNumber      *hexutil.Big `json:"blockNumber,omitempty"`
	TransactionIndex hexutil.Uint `json:"transactionIndex"`
}

func (t TxDetail) MarshalJSON() ([]byte, error) {
	enc := txDetailJSON{
		// Tx fields
		Type:                 hexutil.Uint64(t.Tx.Type),
		ChainID:              (*hexutil.Big)(t.Tx.ChainID),
		Nonce:                hexutil.Uint64(t.Tx.Nonce),
		To:                   t.Tx.To,
		Gas:                  hexutil.Uint64(t.Tx.Gas),
		GasPrice:             (*hexutil.Big)(t.Tx.GasPrice),
		MaxPriorityFeePerGas: (*hexutil.Big)(t.Tx.MaxPriorityFeePerGas),
		MaxFeePerGas:         (*hexutil.Big)(t.Tx.MaxFeePerGas),
		MaxFeePerBlobGas:     (*hexutil.Big)(t.Tx.MaxFeePerBlobGas),
		Value:                (*hexutil.Big)(t.Tx.Value),
		Input:                t.Tx.Input,
		//AccessList:           t.Tx.AccessList,
		BlobVersionedHashes: t.Tx.BlobVersionedHashes,
		AuthorizationList:   t.Tx.AuthorizationList,
		V:                   (*hexutil.Big)(t.Tx.V),
		R:                   (*hexutil.Big)(t.Tx.R),
		S:                   (*hexutil.Big)(t.Tx.S),
		YParity:             hexutil.Uint64(t.Tx.YParity),
		Hash:                t.Tx.Hash,

		// Receipt fields
		PostState:         t.PostState,
		Status:            hexutil.Uint64(t.Status),
		CumulativeGasUsed: hexutil.Uint64(t.CumulativeGasUsed),
		Bloom:             t.Bloom,
		Logs:              t.Logs,
		ContractAddress:   t.ContractAddress,
		From:              t.From,
		GasUsed:           hexutil.Uint64(t.GasUsed),
		EffectiveGasPrice: (*hexutil.Big)(t.EffectiveGasPrice),
		BlobGasUsed:       hexutil.Uint64(t.BlobGasUsed),
		BlobGasPrice:      (*hexutil.Big)(t.BlobGasPrice),
		BlockHash:         t.BlockHash,
		BlockNumber:       (*hexutil.Big)(t.BlockNumber),
		TransactionIndex:  hexutil.Uint(t.TransactionIndex),
	}
	return json.Marshal(&enc)
}

func (t *TxDetail) UnmarshalJSON(input []byte) error {
	var dec txDetailJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Tx fields
	t.Tx.Type = uint8(dec.Type)
	t.Tx.ChainID = (*big.Int)(dec.ChainID)
	t.Tx.Nonce = uint64(dec.Nonce)
	t.Tx.To = dec.To
	t.Tx.Gas = uint64(dec.Gas)
	t.Tx.GasPrice = (*big.Int)(dec.GasPrice)
	t.Tx.MaxPriorityFeePerGas = (*big.Int)(dec.MaxPriorityFeePerGas)
	t.Tx.MaxFeePerGas = (*big.Int)(dec.MaxFeePerGas)
	t.Tx.MaxFeePerBlobGas = (*big.Int)(dec.MaxFeePerBlobGas)
	t.Tx.Value = (*big.Int)(dec.Value)
	t.Tx.Input = dec.Input
	//t.Tx.AccessList = dec.AccessList
	t.Tx.BlobVersionedHashes = dec.BlobVersionedHashes
	t.Tx.AuthorizationList = dec.AuthorizationList
	t.Tx.V = (*big.Int)(dec.V)
	t.Tx.R = (*big.Int)(dec.R)
	t.Tx.S = (*big.Int)(dec.S)
	t.Tx.YParity = uint64(dec.YParity)
	t.Tx.Hash = dec.Hash

	// Receipt fields
	t.PostState = dec.PostState
	t.Status = uint64(dec.Status)
	t.CumulativeGasUsed = uint64(dec.CumulativeGasUsed)
	t.Bloom = dec.Bloom
	t.Logs = dec.Logs
	t.ContractAddress = dec.ContractAddress
	t.From = dec.From
	t.GasUsed = uint64(dec.GasUsed)
	t.EffectiveGasPrice = (*big.Int)(dec.EffectiveGasPrice)
	t.BlobGasUsed = uint64(dec.BlobGasUsed)
	t.BlobGasPrice = (*big.Int)(dec.BlobGasPrice)
	t.BlockHash = dec.BlockHash
	t.BlockNumber = (*big.Int)(dec.BlockNumber)
	t.TransactionIndex = uint(dec.TransactionIndex)
	return nil
}
