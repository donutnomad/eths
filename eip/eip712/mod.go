package eip712

import (
	"github.com/donutnomad/eths/multiread"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var _pack = NewEIP712()

func GetEIP712Domain(contract common.Address, client bind.ContractCaller) (*Eip712DomainOutput, error) {
	return multiread.CALLD(client, contract, _pack.PackEip712Domain(), _pack.UnpackEip712Domain)
}

func GetDomainTypes(domain Eip712DomainOutput) ([]apitypes.Type, apitypes.TypedDataDomain) {
	var domainTypes = make([]apitypes.Type, 0, 5)
	var dataDomain = apitypes.TypedDataDomain{}
	var fields = domain.Fields[0]

	if fields&0x01 != 0 {
		dataDomain.Name = domain.Name
		domainTypes = append(domainTypes, apitypes.Type{Name: "name", Type: "string"})
	}
	if fields&0x02 != 0 {
		dataDomain.Version = domain.Version
		domainTypes = append(domainTypes, apitypes.Type{Name: "version", Type: "string"})
	}
	if fields&0x04 != 0 {
		dataDomain.ChainId = math.NewHexOrDecimal256(domain.ChainId.Int64())
		domainTypes = append(domainTypes, apitypes.Type{Name: "chainId", Type: "uint256"})
	}
	if fields&0x08 != 0 {
		dataDomain.VerifyingContract = domain.VerifyingContract.Hex()
		domainTypes = append(domainTypes, apitypes.Type{Name: "verifyingContract", Type: "address"})
	}
	if fields&0x10 != 0 {
		dataDomain.Salt = hexutil.Encode(domain.Salt[:])
		domainTypes = append(domainTypes, apitypes.Type{Name: "salt", Type: "bytes32"})
	}
	return domainTypes, dataDomain
}
