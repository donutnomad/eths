package eip712

import (
	"github.com/donutnomad/eths/multiread"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

var _pack = NewEIP712()

func GetEIP712Domain(contract common.Address, client bind.ContractCaller) (*Eip712DomainOutput, error) {
	return multiread.CALLD(client, contract, _pack.PackEip712Domain(), _pack.UnpackEip712Domain)
}
