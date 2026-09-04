package multicall

import (
	"math/big"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/contracts_pack"
	"github.com/donutnomad/eths/multiread"
)

func GetAddress[T *big.Int | int | uint | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](chainID T) common.Address {
	return multiread.GetAddress(chainID)
}

type Multicall3Call3 = multiread.Multicall3Call3
type Multicall3Result = multiread.Multicall3Result

func One(contractAddress common.Address, allowFailure bool, callData []byte) Multicall3Call3 {
	return Multicall3Call3{
		Target:       contractAddress,
		AllowFailure: allowFailure,
		CallData:     callData,
	}
}

func Pack3(calls ...Multicall3Call3) []byte {
	if len(calls) == 0 {
		panic("invalid parameter")
	}
	return contracts_pack.NewMulticall().PackAggregate3(calls)
}
