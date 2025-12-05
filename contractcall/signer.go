package contractcall

import (
	"github.com/donutnomad/blockchain-alg/xecdsa"
	"github.com/donutnomad/blockchain-alg/xsecp256k1"
	"github.com/ethereum/go-ethereum/common"
)

type ISigner interface {
	Address() common.Address
	Sign(msg []byte) (*xecdsa.RSVSignature, error)
}

type NoOpSigner struct {
	address common.Address
	signFn  func(data []byte)
}

func NewNoOpSigner(address common.Address, signFn func(data []byte)) *NoOpSigner {
	return &NoOpSigner{address: address, signFn: signFn}
}

func (s *NoOpSigner) Address() common.Address {
	return s.address
}

func (s *NoOpSigner) Sign(data []byte) (*xecdsa.RSVSignature, error) {
	if s.signFn != nil {
		s.signFn(data)
	}
	return nil, nil
}

type EcdsaPrivateKeySigner struct {
	privateKey *xecdsa.PrivateKey
}

func NewEcdsaPrivateKeySigner(privateKey *xecdsa.PrivateKey) *EcdsaPrivateKeySigner {
	return &EcdsaPrivateKeySigner{privateKey: privateKey}
}

func (s *EcdsaPrivateKeySigner) Address() common.Address {
	return common.Address(xsecp256k1.NewPublicKeyFromEcdsa(&s.privateKey.PublicKey).Address())
}

func (s *EcdsaPrivateKeySigner) Sign(msg []byte) (*xecdsa.RSVSignature, error) {
	sig, err := s.privateKey.Sign(msg)
	if err != nil {
		return nil, err
	}
	return sig.(*xecdsa.RSVSignature), nil
}
