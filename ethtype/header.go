package ethtype

import (
	"math/big"

	"github.com/donutnomad/eths/ecommon"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate go run github.com/fjl/gencodec -type Header -field-override headerMarshaling -out header_generated.go

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	Hash       ecommon.Hash    `json:"hash"`
	ParentHash ecommon.Hash    `json:"parentHash"`
	UncleHash  ecommon.Hash    `json:"sha3Uncles"`
	Coinbase   ecommon.Address `json:"miner"`

	Root        ecommon.Hash `json:"stateRoot"`
	TxHash      ecommon.Hash `json:"transactionsRoot"`
	ReceiptHash ecommon.Hash `json:"receiptsRoot"`
	Bloom       Bloom        `json:"logsBloom"`
	Difficulty  *big.Int     `json:"difficulty"`
	Number      *big.Int     `json:"number"`
	GasLimit    uint64       `json:"gasLimit"`
	GasUsed     uint64       `json:"gasUsed"`
	Time        uint64       `json:"timestamp"`
	Extra       []byte       `json:"extraData"`
	MixDigest   common.Hash  `json:"mixHash"`
	Nonce       BlockNonce   `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`

	// WithdrawalsHash was added by EIP-4895 and is ignored in legacy headers.
	WithdrawalsHash *common.Hash `json:"withdrawalsRoot" rlp:"optional"`

	// BlobGasUsed was added by EIP-4844 and is ignored in legacy headers.
	BlobGasUsed *uint64 `json:"blobGasUsed" rlp:"optional"`

	// ExcessBlobGas was added by EIP-4844 and is ignored in legacy headers.
	ExcessBlobGas *uint64 `json:"excessBlobGas" rlp:"optional"`

	// ParentBeaconRoot was added by EIP-4788 and is ignored in legacy headers.
	ParentBeaconRoot *ecommon.Hash `json:"parentBeaconBlockRoot" rlp:"optional"`

	// RequestsHash was added by EIP-7685 and is ignored in legacy headers.
	RequestsHash *ecommon.Hash `json:"requestsHash" rlp:"optional"`

	TotalDifficulty *big.Int `json:"totalDifficulty"`
	Size            uint64   `json:"size"`
	// Transactions    []ecommon.Hash `json:"transactions"`
	// Uncles          []ecommon.Hash `json:"uncles"`

	// arbitrum
	L1Number  *big.Int      `json:"l1BlockNumber,omitempty"`
	SendCount *big.Int      `json:"sendCount,omitempty"`
	SendRoot  *ecommon.Hash `json:"sendRoot,omitempty"`
}

func (h Header) NumberU64() uint64 {
	return h.Number.Uint64()
}

// field type overrides for gencodec
type headerMarshaling struct {
	Difficulty      *hexutil.Big
	Number          *hexutil.Big
	GasLimit        hexutil.Uint64
	GasUsed         hexutil.Uint64
	Time            hexutil.Uint64
	Extra           hexutil.Bytes
	BaseFee         *hexutil.Big
	BlobGasUsed     *hexutil.Uint64
	ExcessBlobGas   *hexutil.Uint64
	TotalDifficulty *hexutil.Big
	Size            hexutil.Uint64

	L1Number  *hexutil.Big
	SendCount *hexutil.Big
}
