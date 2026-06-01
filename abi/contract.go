package abi

import (
	"errors"

	"github.com/donutnomad/eths/common"
)

// Common errors shared by generated contract bindings.
var (
	ErrNoEventSignature        = errors.New("missing event signature")
	ErrEventSignatureMismatch  = errors.New("event signature mismatch")
	ErrMethodSignatureMismatch = errors.New("method signature mismatch")
	ErrErrorSignatureMismatch  = errors.New("error signature mismatch")
)

// ContractEvent is implemented by every generated contract event binding.
type ContractEvent interface {
	Name() string
	Topic0() common.Hash
}

// ContractError is implemented by every generated contract error binding.
type ContractError interface {
	Name() string
	ErrorID() common.Hash
}
