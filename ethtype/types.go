package ethtype

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	hexutil2 "github.com/donutnomad/eths/common/hexutil"
	"github.com/donutnomad/eths/ecommon"
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil2.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil2.UnmarshalFixedText("BlockNonce", input, n[:])
}

type BlockNumber int64

const (
	EarliestBlockNumber  = BlockNumber(-5)
	SafeBlockNumber      = BlockNumber(-4)
	FinalizedBlockNumber = BlockNumber(-3)
	LatestBlockNumber    = BlockNumber(-2)
	PendingBlockNumber   = BlockNumber(-1)
)

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "safe", "finalized", "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*bn = EarliestBlockNumber
		return nil
	case "latest":
		*bn = LatestBlockNumber
		return nil
	case "pending":
		*bn = PendingBlockNumber
		return nil
	case "finalized":
		*bn = FinalizedBlockNumber
		return nil
	case "safe":
		*bn = SafeBlockNumber
		return nil
	}

	blckNum, err := hexutil2.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return errors.New("block number larger than int64")
	}
	*bn = BlockNumber(blckNum)
	return nil
}

// Int64 returns the block number as int64.
func (bn BlockNumber) Int64() int64 {
	return (int64)(bn)
}

// MarshalText implements encoding.TextMarshaler. It marshals:
// - "safe", "finalized", "latest", "earliest" or "pending" as strings
// - other numbers as hex
func (bn BlockNumber) MarshalText() ([]byte, error) {
	return []byte(bn.String()), nil
}

func (bn BlockNumber) String() string {
	switch bn {
	case EarliestBlockNumber:
		return "earliest"
	case LatestBlockNumber:
		return "latest"
	case PendingBlockNumber:
		return "pending"
	case FinalizedBlockNumber:
		return "finalized"
	case SafeBlockNumber:
		return "safe"
	default:
		if bn < 0 {
			return fmt.Sprintf("<invalid %d>", bn)
		}
		return hexutil2.Uint64(bn).String()
	}
}

type BlockID = BlockNumberOrHash

type BlockNumberOrHash struct {
	BlockNumber      *BlockNumber  `json:"blockNumber,omitempty"`
	BlockHash        *ecommon.Hash `json:"blockHash,omitempty"`
	RequireCanonical bool          `json:"requireCanonical,omitempty"`
}

func BlockNumberOrHashWithNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](blockNr T) BlockNumberOrHash {
	return *BlockNum(blockNr)
}

func BlockNum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](blockNr T) *BlockNumberOrHash {
	b := BlockNumber(blockNr)
	return &BlockNumberOrHash{
		BlockNumber:      &b,
		BlockHash:        nil,
		RequireCanonical: false,
	}
}

func BlockHash(hash ecommon.Hash, canonical ...bool) *BlockNumberOrHash {
	return &BlockNumberOrHash{
		BlockNumber:      nil,
		BlockHash:        &hash,
		RequireCanonical: len(canonical) > 0 && canonical[0],
	}
}

func blockNumberOrHashWithTag(bn BlockNumber) *BlockNumberOrHash {
	return &BlockNumberOrHash{
		BlockNumber:      &bn,
		BlockHash:        nil,
		RequireCanonical: false,
	}
}

// EarliestBlock returns a BlockNumberOrHash selecting the genesis (earliest) block.
func EarliestBlock() *BlockNumberOrHash {
	return blockNumberOrHashWithTag(EarliestBlockNumber)
}

// SafeBlock returns a BlockNumberOrHash selecting the latest safe head block.
func SafeBlock() *BlockNumberOrHash {
	return blockNumberOrHashWithTag(SafeBlockNumber)
}

// FinalizedBlock returns a BlockNumberOrHash selecting the latest finalized block.
func FinalizedBlock() *BlockNumberOrHash {
	return blockNumberOrHashWithTag(FinalizedBlockNumber)
}

// LatestBlock returns a BlockNumberOrHash selecting the most recently included block.
func LatestBlock() *BlockNumberOrHash {
	return blockNumberOrHashWithTag(LatestBlockNumber)
}

// PendingBlock returns a BlockNumberOrHash selecting the pending block.
func PendingBlock() *BlockNumberOrHash {
	return blockNumberOrHashWithTag(PendingBlockNumber)
}

func BlockNumberOrHashWithHash(hash ecommon.Hash, canonical bool) BlockNumberOrHash {
	return *BlockHash(hash, canonical)
}

func (bnh *BlockNumberOrHash) UnmarshalJSON(data []byte) error {
	type erased BlockNumberOrHash
	e := erased{}
	err := json.Unmarshal(data, &e)
	if err == nil {
		if e.BlockNumber != nil && e.BlockHash != nil {
			return errors.New("cannot specify both BlockHash and BlockNumber, choose one or the other")
		}
		bnh.BlockNumber = e.BlockNumber
		bnh.BlockHash = e.BlockHash
		bnh.RequireCanonical = e.RequireCanonical
		return nil
	}
	var input string
	err = json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	switch input {
	case "earliest":
		bn := EarliestBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "latest":
		bn := LatestBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "pending":
		bn := PendingBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "finalized":
		bn := FinalizedBlockNumber
		bnh.BlockNumber = &bn
		return nil
	case "safe":
		bn := SafeBlockNumber
		bnh.BlockNumber = &bn
		return nil
	default:
		if len(input) == 66 {
			hash := ecommon.Hash{}
			err := hash.UnmarshalText([]byte(input))
			if err != nil {
				return err
			}
			bnh.BlockHash = &hash
			return nil
		} else {
			blckNum, err := hexutil2.DecodeUint64(input)
			if err != nil {
				return err
			}
			if blckNum > math.MaxInt64 {
				return errors.New("blocknumber too high")
			}
			bn := BlockNumber(blckNum)
			bnh.BlockNumber = &bn
			return nil
		}
	}
}

func (bnh *BlockNumberOrHash) Number() (BlockNumber, bool) {
	if bnh.BlockNumber != nil {
		return *bnh.BlockNumber, true
	}
	return BlockNumber(0), false
}

func (bnh *BlockNumberOrHash) String() string {
	if bnh.BlockNumber != nil {
		return bnh.BlockNumber.String()
	}
	if bnh.BlockHash != nil {
		return bnh.BlockHash.String()
	}
	return "nil"
}

func (bnh *BlockNumberOrHash) Hash() (ecommon.Hash, bool) {
	if bnh.BlockHash != nil {
		return *bnh.BlockHash, true
	}
	return ecommon.Hash{}, false
}
