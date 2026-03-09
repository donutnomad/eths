package ethclient

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

// headerJSON is the shared JSON encoding for block headers.
type headerJSON struct {
	ParentHash       common.Hash      `json:"parentHash"`
	UncleHash        common.Hash      `json:"sha3Uncles"`
	Coinbase         common.Address   `json:"miner"`
	Root             common.Hash      `json:"stateRoot"`
	TxHash           common.Hash      `json:"transactionsRoot"`
	ReceiptHash      common.Hash      `json:"receiptsRoot"`
	Bloom            types.Bloom      `json:"logsBloom"`
	Difficulty       *hexutil.Big     `json:"difficulty"`
	Number           *hexutil.Big     `json:"number"`
	GasLimit         hexutil.Uint64   `json:"gasLimit"`
	GasUsed          hexutil.Uint64   `json:"gasUsed"`
	Time             hexutil.Uint64   `json:"timestamp"`
	Extra            hexutil.Bytes    `json:"extraData"`
	MixDigest        common.Hash      `json:"mixHash"`
	Nonce            types.BlockNonce `json:"nonce"`
	BaseFee          *hexutil.Big     `json:"baseFeePerGas,omitempty"`
	WithdrawalsHash  *common.Hash     `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed      *hexutil.Uint64  `json:"blobGasUsed,omitempty"`
	ExcessBlobGas    *hexutil.Uint64  `json:"excessBlobGas,omitempty"`
	ParentBeaconRoot *common.Hash     `json:"parentBeaconBlockRoot,omitempty"`
	RequestsHash     *common.Hash     `json:"requestsHash,omitempty"`
}

func newHeaderJSON(h *types.Header) headerJSON {
	return headerJSON{
		ParentHash:       h.ParentHash,
		UncleHash:        h.UncleHash,
		Coinbase:         h.Coinbase,
		Root:             h.Root,
		TxHash:           h.TxHash,
		ReceiptHash:      h.ReceiptHash,
		Bloom:            h.Bloom,
		Difficulty:       (*hexutil.Big)(h.Difficulty),
		Number:           (*hexutil.Big)(h.Number),
		GasLimit:         hexutil.Uint64(h.GasLimit),
		GasUsed:          hexutil.Uint64(h.GasUsed),
		Time:             hexutil.Uint64(h.Time),
		Extra:            h.Extra,
		MixDigest:        h.MixDigest,
		Nonce:            h.Nonce,
		BaseFee:          (*hexutil.Big)(h.BaseFee),
		WithdrawalsHash:  h.WithdrawalsHash,
		BlobGasUsed:      (*hexutil.Uint64)(h.BlobGasUsed),
		ExcessBlobGas:    (*hexutil.Uint64)(h.ExcessBlobGas),
		ParentBeaconRoot: h.ParentBeaconRoot,
		RequestsHash:     h.RequestsHash,
	}
}

func (hj *headerJSON) toHeader() *types.Header {
	h := &types.Header{
		ParentHash:       hj.ParentHash,
		UncleHash:        hj.UncleHash,
		Coinbase:         hj.Coinbase,
		Root:             hj.Root,
		TxHash:           hj.TxHash,
		ReceiptHash:      hj.ReceiptHash,
		Bloom:            hj.Bloom,
		Difficulty:       (*big.Int)(hj.Difficulty),
		Number:           (*big.Int)(hj.Number),
		GasLimit:         uint64(hj.GasLimit),
		GasUsed:          uint64(hj.GasUsed),
		Time:             uint64(hj.Time),
		Extra:            hj.Extra,
		MixDigest:        hj.MixDigest,
		Nonce:            hj.Nonce,
		WithdrawalsHash:  hj.WithdrawalsHash,
		ParentBeaconRoot: hj.ParentBeaconRoot,
		RequestsHash:     hj.RequestsHash,
	}
	if hj.BaseFee != nil {
		h.BaseFee = (*big.Int)(hj.BaseFee)
	}
	if hj.BlobGasUsed != nil {
		h.BlobGasUsed = (*uint64)(hj.BlobGasUsed)
	}
	if hj.ExcessBlobGas != nil {
		h.ExcessBlobGas = (*uint64)(hj.ExcessBlobGas)
	}
	return h
}

// blockHash returns the block hash, preferring the explicit hash over computing from header.
func blockHash(explicit *common.Hash, h *types.Header) common.Hash {
	if explicit != nil {
		return *explicit
	}
	return h.Hash()
}

// LiteBlock holds block header and transaction hashes without full tx data or uncles.
type LiteBlock struct {
	*types.Header
	BlockHash    *common.Hash        `json:"hash"`
	Transactions []common.Hash       `json:"transactions"`
	Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
}

// Hash returns the block hash.
func (b *LiteBlock) Hash() common.Hash {
	return blockHash(b.BlockHash, b.Header)
}

func (b *LiteBlock) MarshalJSON() ([]byte, error) {
	type enc struct {
		headerJSON
		Hash         common.Hash         `json:"hash"`
		Transactions []common.Hash       `json:"transactions"`
		Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		headerJSON:   newHeaderJSON(b.Header),
		Hash:         b.Hash(),
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
	})
}

func (b *LiteBlock) UnmarshalJSON(data []byte) error {
	type raw struct {
		headerJSON
		BlockHash    *common.Hash        `json:"hash"`
		Transactions []common.Hash       `json:"transactions"`
		Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	b.Header = r.headerJSON.toHeader()
	b.BlockHash = r.BlockHash
	b.Transactions = r.Transactions
	b.Withdrawals = r.Withdrawals
	return nil
}

// RichBlock holds block header with full transaction objects and withdrawals.
type RichBlock struct {
	*types.Header
	BlockHash    *common.Hash         `json:"hash"`
	Transactions []*types.Transaction `json:"transactions"`
	Withdrawals  []*types.Withdrawal  `json:"withdrawals,omitempty"`
}

// Hash returns the block hash.
func (b *RichBlock) Hash() common.Hash {
	return blockHash(b.BlockHash, b.Header)
}

// ToBlock converts the RichBlock to a types.Block.
func (b *RichBlock) ToBlock() *types.Block {
	return types.NewBlock(b.Header, &types.Body{
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
	}, nil, trie.NewStackTrie(nil))
}

func (b *RichBlock) MarshalJSON() ([]byte, error) {
	type enc struct {
		headerJSON
		Hash         common.Hash          `json:"hash"`
		Transactions []*types.Transaction `json:"transactions"`
		Withdrawals  []*types.Withdrawal  `json:"withdrawals,omitempty"`
	}
	return json.Marshal(&enc{
		headerJSON:   newHeaderJSON(b.Header),
		Hash:         b.Hash(),
		Transactions: b.Transactions,
		Withdrawals:  b.Withdrawals,
	})
}

func (b *RichBlock) UnmarshalJSON(data []byte) error {
	type raw struct {
		headerJSON
		BlockHash    *common.Hash         `json:"hash"`
		Transactions []*types.Transaction `json:"transactions"`
		Withdrawals  []*types.Withdrawal  `json:"withdrawals,omitempty"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	b.Header = r.headerJSON.toHeader()
	b.BlockHash = r.BlockHash
	b.Transactions = r.Transactions
	b.Withdrawals = r.Withdrawals
	return nil
}
