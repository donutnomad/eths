package eip712

import (
	"github.com/donutnomad/eths/abi/bind/v2"
	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/common/hexutil"
	"github.com/donutnomad/eths/common/math"
	"github.com/donutnomad/eths/ethtype"
	"github.com/donutnomad/eths/multiread"
)

var _pack = NewEIP712()

func GetEIP712Domain(contract common.Address, client bind.ContractCaller) (*Eip712DomainOutput, error) {
	return multiread.CALLD(client, contract, _pack.PackEip712Domain(), _pack.UnpackEip712Domain)
}

func GetDomainTypes(domain Eip712DomainOutput) ([]ethtype.TypedDataType, ethtype.TypedDataDomain) {
	var domainTypes = make([]ethtype.TypedDataType, 0, 5)
	var dataDomain = ethtype.TypedDataDomain{}
	var fields = domain.Fields[0]

	if fields&0x01 != 0 {
		dataDomain.Name = domain.Name
		domainTypes = append(domainTypes, ethtype.TypedDataType{Name: "name", Type: "string"})
	}
	if fields&0x02 != 0 {
		dataDomain.Version = domain.Version
		domainTypes = append(domainTypes, ethtype.TypedDataType{Name: "version", Type: "string"})
	}
	if fields&0x04 != 0 {
		dataDomain.ChainId = math.NewHexOrDecimal256(domain.ChainId.Int64())
		domainTypes = append(domainTypes, ethtype.TypedDataType{Name: "chainId", Type: "uint256"})
	}
	if fields&0x08 != 0 {
		dataDomain.VerifyingContract = domain.VerifyingContract.Hex()
		domainTypes = append(domainTypes, ethtype.TypedDataType{Name: "verifyingContract", Type: "address"})
	}
	if fields&0x10 != 0 {
		dataDomain.Salt = hexutil.Encode(domain.Salt[:])
		domainTypes = append(domainTypes, ethtype.TypedDataType{Name: "salt", Type: "bytes32"})
	}
	return domainTypes, dataDomain
}
