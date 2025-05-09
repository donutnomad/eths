package contractcall

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/shopspring/decimal"
)

func covertGweiToWei(gwei string) *big.Int {
	fromString, err := decimal.NewFromString(gwei)
	if err != nil {
		return new(big.Int).SetUint64(0)
	}
	fromString = fromString.Mul(decimal.New(1, 9))
	return fromString.BigInt()
}

type ILogger interface {
	Printf(string, ...interface{})
}

func logCallMsg(logger ILogger, msg *ethereum.CallMsg) {
	if logger == nil || msg == nil {
		return
	}
	logger.Printf(fmt.Sprintf("estimateGasLimit: gasPrice:%s, gasTipCap:%s, gasFeeCap:%s\n",
		bigIntToString(msg.GasPrice),
		bigIntToString(msg.GasTipCap),
		bigIntToString(msg.GasFeeCap)),
	)
	logger.Printf(fmt.Sprintf("from: %v, to:%v, value:%v, data:%x\n", msg.From, msg.To, msg.Value, msg.Data))
}

func bigIntToString(input *big.Int) string {
	if input == nil {
		return "<nil>"
	}
	return input.String()
}

type KnownError struct {
	Name      string
	Signature string
	Arguments []string
}
