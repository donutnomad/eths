// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethtype

import (
	"hash"

	"github.com/donutnomad/eths/ecommon"
	"golang.org/x/crypto/sha3"
)

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h ecommon.Hash) {
	d := hasherPool.Get().(KeccakState)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	hasherPool.Put(d)
	return h
}

var (
	// EmptyRootHash is the known root hash of an empty merkle trie.
	EmptyRootHash = ecommon.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyUncleHash is the known hash of the empty uncle set.
	EmptyUncleHash = ecommon.HexToHash("1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")

	// EmptyCodeHash is the known hash of the empty EVM bytecode.
	EmptyCodeHash = Keccak256Hash(nil) // c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470

	// EmptyTxsHash is the known hash of the empty transaction set.
	EmptyTxsHash = ecommon.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyReceiptsHash is the known hash of the empty receipt set.
	EmptyReceiptsHash = ecommon.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyWithdrawalsHash is the known hash of the empty withdrawal set.
	EmptyWithdrawalsHash = ecommon.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyRequestsHash is the known hash of an empty request set, sha256("").
	EmptyRequestsHash = ecommon.HexToHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	// EmptyVerkleHash is the known hash of an empty verkle trie.
	EmptyVerkleHash = ecommon.Hash{}

	// EmptyBinaryHash is the known hash of an empty binary trie.
	EmptyBinaryHash = ecommon.Hash{}
)
