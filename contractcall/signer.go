package contractcall

import (
	"crypto/ecdsa"
	"github.com/donutnomad/blockchain-alg/xecdsa"
)

type ISigner interface {
	PublicKey() ecdsa.PublicKey
	Sign(msg []byte) (xecdsa.ISignature, error)
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

func (s *EcdsaPrivateKeySigner) Sign(msg []byte) (xecdsa.ISignature, error) {
	return s.privateKey.Sign(msg)
}
