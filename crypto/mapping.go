package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/donutnomad/eths/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/keccak"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"
)

const SignatureLength = crypto.SignatureLength

type KeccakState = crypto.KeccakState

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return keccak.NewLegacyKeccak256().(KeccakState)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	return crypto.Keccak256(data...)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	return crypto.Keccak256Hash(data...)
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	return crypto.ValidateSignatureValues(v, r, s, homestead)
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversary. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	return crypto.Sign(digestHash, prv)
}

// GenerateKey generates a new private key.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b common.Address, nonce uint64) common.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return common.BytesToAddress(Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2(b common.Address, salt [32]byte, inithash []byte) common.Address {
	return common.BytesToAddress(Keccak256([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}
