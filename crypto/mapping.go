package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"hash"
	"math/big"
	"sync"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp256k1ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/common/math"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

// SignatureLength indicates the byte length required to carry a signature with recovery id.
const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

var hasherPool = sync.Pool{
	New: func() any {
		return sha3.NewLegacyKeccak256()
	},
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := hasherPool.Get().(hash.Hash)
	defer hasherPool.Put(d)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	ret := Keccak256(data...)
	copy(h[:], ret[:32])
	return h
}

var (
	secp256k1N     = secp256k1.S256().Params().N
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if homestead && s.Cmp(secp256k1halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", 32, len(hash))
	}
	if len(sig) != SignatureLength {
		return nil, fmt.Errorf("signature is required to be exactly %d bytes (%d)", SignatureLength, len(sig))
	}

	compactSig := make([]byte, SignatureLength)
	compactSig[0] = sig[64] + 27
	copy(compactSig[1:33], sig[:32])
	copy(compactSig[33:], sig[32:64])

	pub, _, err := secp256k1ecdsa.RecoverCompact(compactSig, hash)
	if err != nil {
		return nil, err
	}
	return pub.SerializeUncompressed(), nil
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
	if len(digestHash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", 32, len(digestHash))
	}
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer func() {
		clear(seckey)
	}()
	key := secp256k1.PrivKeyFromBytes(seckey)
	compactSig := secp256k1ecdsa.SignCompact(key, digestHash, false)

	sig = make([]byte, SignatureLength)
	copy(sig[:64], compactSig[1:])
	sig[64] = compactSig[0] - 27
	return sig, nil
}

// GenerateKey generates a new private key.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	key, err := secp256k1.GeneratePrivateKeyFromRand(rand.Reader)
	if err != nil {
		return nil, err
	}
	return key.ToECDSA(), nil
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
