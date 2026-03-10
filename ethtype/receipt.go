package ethtype

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

//go:generate go tool github.com/fjl/gencodec -type TxReceipt -field-override receiptMarshaling -out receipt_generated.go

type Receipt = TxReceipt
type Receipts = []*Receipt

type TxReceipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Type              uint8           `json:"type,omitempty"`
	PostState         []byte          `json:"root"`
	Status            uint64          `json:"status"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             ethTypes.Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs              []*ethTypes.Log `json:"logs"              gencodec:"required"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	TxHash            common.Hash     `json:"transactionHash" gencodec:"required"`
	ContractAddress   *common.Address `json:"contractAddress"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	GasUsed           uint64          `json:"gasUsed" gencodec:"required"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"` // required, but tag omitted for backwards compatibility
	BlobGasUsed       uint64          `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int        `json:"blobGasPrice,omitempty"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
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
