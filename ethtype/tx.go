package ethtype

import (
	"math/big"

	"github.com/donutnomad/eths/ecommon"
	"github.com/donutnomad/eths/hexutil"
)

type Transaction = Tx

// Object - A transaction object, or null when no transaction was found:
//
//blockHash: DATA, 32 Bytes - hash of the block where this transaction was in. null when its pending.
//blockNumber: QUANTITY - block number where this transaction was in. null when its pending.
//from: DATA, 20 Bytes - address of the sender.
//gas: QUANTITY - gas provided by the sender.
//gasPrice: QUANTITY - gas price provided by the sender in Wei.
//hash: DATA, 32 Bytes - hash of the transaction.
//input: DATA - the data send along with the transaction.
//nonce: QUANTITY - the number of transactions made by the sender prior to this one.
//to: DATA, 20 Bytes - address of the receiver. null when its a contract creation transaction.
//transactionIndex: QUANTITY - integer of the transactions index position in the block. null when its pending.
//value: QUANTITY - value transferred in Wei.
//v: QUANTITY - ECDSA recovery id
//r: QUANTITY - ECDSA signature r
//s: QUANTITY - ECDSA signature s

// Tx RPC: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
type Tx struct {
	Type uint8 `json:"type"`

	ChainID              *big.Int               `json:"chainId,omitempty"`
	Nonce                uint64                 `json:"nonce"`
	From                 ecommon.Address        `json:"from"`
	To                   *ecommon.Address       `json:"to"`
	Gas                  uint64                 `json:"gas"`
	GasPrice             *big.Int               `json:"gasPrice"`
	MaxPriorityFeePerGas *big.Int               `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *big.Int               `json:"maxFeePerGas"`
	MaxFeePerBlobGas     *big.Int               `json:"maxFeePerBlobGas,omitempty"`
	Value                *big.Int               `json:"value"`
	Input                []byte                 `json:"input"`
	AccessList           AccessList             `json:"accessList,omitempty"` // EIP-2930 access list
	BlobVersionedHashes  []ecommon.Hash         `json:"blobVersionedHashes,omitempty"`
	AuthorizationList    []SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                    *big.Int               `json:"v"`
	R                    *big.Int               `json:"r"`
	S                    *big.Int               `json:"s"`
	YParity              uint64                 `json:"yParity,omitempty"`

	// Only used for encoding:
	Hash ecommon.Hash `json:"hash"`
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
