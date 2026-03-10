package ethtype

import (
	"math/big"

	"github.com/donutnomad/eths/ecommon"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Receipt = TxReceipt
type Receipts = []*Receipt

type TxReceipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Type              uint8  `json:"type,omitempty"`
	PostState         []byte `json:"root"`
	Status            uint64 `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs              []*Log `json:"logs"              gencodec:"required"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	TxHash            ecommon.Hash     `json:"transactionHash" gencodec:"required"`
	ContractAddress   *ecommon.Address `json:"contractAddress"`
	From              ecommon.Address  `json:"from"`
	To                *ecommon.Address `json:"to"`
	GasUsed           uint64           `json:"gasUsed" gencodec:"required"`
	EffectiveGasPrice *big.Int         `json:"effectiveGasPrice"` // required, but tag omitted for backwards compatibility
	BlobGasUsed       uint64           `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int         `json:"blobGasPrice,omitempty"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        ecommon.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int     `json:"blockNumber,omitempty"`
	TransactionIndex uint         `json:"transactionIndex"`
}

type receiptMarshaling struct {
	Type              hexutil.Uint64
	PostState         hexutil.Bytes
	Status            hexutil.Uint64
	CumulativeGasUsed hexutil.Uint64
	GasUsed           hexutil.Uint64
	EffectiveGasPrice *hexutil.Big
	BlobGasUsed       hexutil.Uint64
	BlobGasPrice      *hexutil.Big
	BlockNumber       *hexutil.Big
	TransactionIndex  hexutil.Uint
}
