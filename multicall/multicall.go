package multicall

import (
	"github.com/donutnomad/eths/contracts_pack"
	"github.com/ethereum/go-ethereum/common"
)

var Address = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")

type Multicall3Call3 = contracts_pack.Multicall3Call3
type Multicall3Result = contracts_pack.Multicall3Result

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
