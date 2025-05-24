package contractcall

import (
	"crypto/ecdsa"
	"github.com/donutnomad/blockchain-alg/xecdsa"
	"math/big"
)

type ISigner interface {
	PublicKey() ecdsa.PublicKey
	Sign(msg []byte) (*xecdsa.RSVSignature, error)
}

type NoOpSigner struct {
	key ecdsa.PublicKey
}

func (s *NoOpSigner) PublicKey() ecdsa.PublicKey {
	return s.key
}

func (s *NoOpSigner) Sign(_ []byte) (*xecdsa.RSVSignature, error) {
	v := byte(27)
	return xecdsa.NewSignature(big.NewInt(0), big.NewInt(0), &v).(*xecdsa.RSVSignature), nil
}

type EcdsaPrivateKeySigner struct {
	privateKey *xecdsa.PrivateKey
}

func NewEcdsaPrivateKeySigner(privateKey *xecdsa.PrivateKey) *EcdsaPrivateKeySigner {
	return &EcdsaPrivateKeySigner{privateKey: privateKey}
}

func (s *EcdsaPrivateKeySigner) PublicKey() ecdsa.PublicKey {
	return s.privateKey.PublicKey
}

func (s *EcdsaPrivateKeySigner) Sign(msg []byte) (*xecdsa.RSVSignature, error) {
	sig, err := s.privateKey.Sign(msg)
	if err != nil {
		return nil, err
	}
	return sig.(*xecdsa.RSVSignature), nil
}
