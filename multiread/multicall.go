package multiread

import (
	"fmt"
	"github.com/donutnomad/eths/contracts_pack"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"math/big"
	"reflect"
	"sync"
)

type Multicall3Call3 = contracts_pack.Multicall3Call3
type Multicall3Result = contracts_pack.Multicall3Result

// Address Multicall3: https://www.multicall3.com/abi#ethers-js
var Address = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")

var multiCallPack = contracts_pack.NewMulticall()

type ReturnUnPackFunc[T any] = func([]byte) (T, error)
type Func1[T any] func() (common.Address, []byte, func([]byte) (T, error))
type Func2 func() (common.Address, []byte)

func (t Func1[T]) Downcast() Func2 {
	return func() (common.Address, []byte) {
		a, b, _ := t()
		return a, b
	}
}

func (t Func1[T]) Any() Func1[any] {
	return func() (common.Address, []byte, func([]byte) (any, error)) {
		a, b, c := t()
		return a, b, func(bytes []byte) (any, error) {
			return c(bytes)
		}
	}
}

func One[T any](contractAddress common.Address, callData []byte, returnUnpack ReturnUnPackFunc[T]) Func1[T] {
	return func() (common.Address, []byte, func([]byte) (T, error)) {
		return contractAddress, callData, returnUnpack
	}
}

func Any[T any](contractAddress common.Address, callData []byte, returnUnpack ReturnUnPackFunc[T]) Func1[any] {
	return func() (common.Address, []byte, func([]byte) (any, error)) {
		return contractAddress, callData, func(bytes []byte) (any, error) {
			return returnUnpack(bytes)
		}
	}
}

func One2(contractAddress common.Address, callData []byte) Func2 {
	return func() (common.Address, []byte) {
		return contractAddress, callData
	}
}

func GetChainID() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetChainId(), multiCallPack.UnpackGetChainId)
}

func GetBaseFee() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetBasefee(), multiCallPack.UnpackGetBasefee)
}

func GetBlockNumber() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetBlockNumber(), multiCallPack.UnpackGetBlockNumber)
}

func GetCurrentBlockTimestamp() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetCurrentBlockTimestamp(), multiCallPack.UnpackGetCurrentBlockTimestamp)
}

func GetCurrentBlockGasLimit() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetCurrentBlockGasLimit(), multiCallPack.UnpackGetCurrentBlockGasLimit)
}

func GetCurrentBlockDifficulty() Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetCurrentBlockDifficulty(), multiCallPack.UnpackGetCurrentBlockDifficulty)
}

func GetCurrentBlockCoinbase() Func1[common.Address] {
	return One(Address, multiCallPack.PackGetCurrentBlockCoinbase(), multiCallPack.UnpackGetCurrentBlockCoinbase)
}

func GetEthBalance(addr common.Address) Func1[*big.Int] {
	return One(Address, multiCallPack.PackGetEthBalance(addr), multiCallPack.UnpackGetEthBalance)
}

func GetBlockHash(blockNumber *big.Int) Func1[common.Hash] {
	return One(Address, multiCallPack.PackGetBlockHash(blockNumber), func(bytes []byte) (common.Hash, error) {
		ret, err := multiCallPack.UnpackGetBlockHash(bytes)
		if err != nil {
			return common.Hash{}, err
		}
		return ret, nil
	})
}

func GetLastBlockHash() Func1[common.Hash] {
	return One(Address, multiCallPack.PackGetLastBlockHash(), func(bytes []byte) (common.Hash, error) {
		ret, err := multiCallPack.UnpackGetLastBlockHash(bytes)
		if err != nil {
			return common.Hash{}, err
		}
		return ret, nil
	})
}

func AllSuccess(args ...any) bool {
	for _, item := range args {
		if lo.IsNil(item) {
			return false
		}
	}
	return true
}

func CALLAny[A1 any](
	client bind.ContractCaller,
	unpack func([]byte) (A1, error),
	inputs ...Func2,
) ([]*A1, error) {
	if len(inputs) == 0 {
		panic("invalid inputs")
	}
	var args []Multicall3Call3
	var functions []func([]byte) (any, error)
	for _, input := range inputs {
		a1, a2 := input()
		args = append(args, Multicall3Call3{
			Target:       a1,
			CallData:     a2,
			AllowFailure: true,
		})
		functions = append(functions, func(bytes []byte) (any, error) {
			return unpack(bytes)
		})
	}
	var results = make([]any, len(inputs))
	if err := callN(client, args, functions, results); err != nil {
		return nil, err
	}
	var result = make([]*A1, len(inputs))
	for i, item := range results {
		if !lo.IsNil(item) {
			result[i] = lo.ToPtr(item.(A1))
		}
	}
	return result, nil
}

func CALL1[A1 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
) (*A1, error) {
	a, r := prepareMultiCallArg(a1)
	var results = make([]any, 1)
	err := callN(client, defSlice(a), defSlice(r), results)
	if err != nil {
		return nil, err
	}
	if !lo.IsNil(results[0]) {
		return lo.ToPtr(results[0].(A1)), nil
	}
	return nil, nil
}

func CALL2[A1 any, A2 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
) (rr1 *A1, rr2 *A2, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	var results = make([]any, 2)
	err := callN(client, defSlice(z1, z2), defSlice(r1, r2), results)
	if err != nil {
		return nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	return
}

func CALL3[A1 any, A2 any, A3 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	var results = make([]any, 3)
	err := callN(client, defSlice(z1, z2, z3), defSlice(r1, r2, r3), results)
	if err != nil {
		return nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	return
}

func CALL4[A1 any, A2 any, A3 any, A4 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	var results = make([]any, 4)
	err := callN(client, defSlice(z1, z2, z3, z4), defSlice(r1, r2, r3, r4), results)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	return
}

func CALL5[A1 any, A2 any, A3 any, A4 any, A5 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	var results = make([]any, 5)
	err := callN(client, defSlice(z1, z2, z3, z4, z5), defSlice(r1, r2, r3, r4, r5), results)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	return
}

func CALL6[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
	a6 func() (common.Address, []byte, func([]byte) (A6, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	z6, r6 := prepareMultiCallArg(a6)
	var results = make([]any, 6)
	err := callN(client, defSlice(z1, z2, z3, z4, z5, z6), defSlice(r1, r2, r3, r4, r5, r6), results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	if !lo.IsNil(results[5]) {
		rr6 = lo.ToPtr(results[5].(A6))
	}
	return
}

func CALL7[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
	a6 func() (common.Address, []byte, func([]byte) (A6, error)),
	a7 func() (common.Address, []byte, func([]byte) (A7, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, rr7 *A7, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	z6, r6 := prepareMultiCallArg(a6)
	z7, r7 := prepareMultiCallArg(a7)
	var results = make([]any, 7)
	err := callN(client, defSlice(z1, z2, z3, z4, z5, z6, z7), defSlice(r1, r2, r3, r4, r5, r6, r7), results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	if !lo.IsNil(results[5]) {
		rr6 = lo.ToPtr(results[5].(A6))
	}
	if !lo.IsNil(results[6]) {
		rr7 = lo.ToPtr(results[6].(A7))
	}
	return
}

func CALL8[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
	a6 func() (common.Address, []byte, func([]byte) (A6, error)),
	a7 func() (common.Address, []byte, func([]byte) (A7, error)),
	a8 func() (common.Address, []byte, func([]byte) (A8, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, rr7 *A7, rr8 *A8, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	z6, r6 := prepareMultiCallArg(a6)
	z7, r7 := prepareMultiCallArg(a7)
	z8, r8 := prepareMultiCallArg(a8)
	var results = make([]any, 8)
	err := callN(client, defSlice(z1, z2, z3, z4, z5, z6, z7, z8), defSlice(r1, r2, r3, r4, r5, r6, r7, r8), results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	if !lo.IsNil(results[5]) {
		rr6 = lo.ToPtr(results[5].(A6))
	}
	if !lo.IsNil(results[6]) {
		rr7 = lo.ToPtr(results[6].(A7))
	}
	if !lo.IsNil(results[7]) {
		rr8 = lo.ToPtr(results[7].(A8))
	}
	return
}

func CALL9[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, A9 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
	a6 func() (common.Address, []byte, func([]byte) (A6, error)),
	a7 func() (common.Address, []byte, func([]byte) (A7, error)),
	a8 func() (common.Address, []byte, func([]byte) (A8, error)),
	a9 func() (common.Address, []byte, func([]byte) (A9, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, rr7 *A7, rr8 *A8, rr9 *A9, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	z6, r6 := prepareMultiCallArg(a6)
	z7, r7 := prepareMultiCallArg(a7)
	z8, r8 := prepareMultiCallArg(a8)
	z9, r9 := prepareMultiCallArg(a9)
	var results = make([]any, 9)
	err := callN(client, defSlice(z1, z2, z3, z4, z5, z6, z7, z8, z9), defSlice(r1, r2, r3, r4, r5, r6, r7, r8, r9), results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	if !lo.IsNil(results[5]) {
		rr6 = lo.ToPtr(results[5].(A6))
	}
	if !lo.IsNil(results[6]) {
		rr7 = lo.ToPtr(results[6].(A7))
	}
	if !lo.IsNil(results[7]) {
		rr8 = lo.ToPtr(results[7].(A8))
	}
	if !lo.IsNil(results[8]) {
		rr9 = lo.ToPtr(results[8].(A9))
	}
	return
}

func CALL10[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, A9 any, A10 any](
	client bind.ContractCaller,
	a1 func() (common.Address, []byte, func([]byte) (A1, error)),
	a2 func() (common.Address, []byte, func([]byte) (A2, error)),
	a3 func() (common.Address, []byte, func([]byte) (A3, error)),
	a4 func() (common.Address, []byte, func([]byte) (A4, error)),
	a5 func() (common.Address, []byte, func([]byte) (A5, error)),
	a6 func() (common.Address, []byte, func([]byte) (A6, error)),
	a7 func() (common.Address, []byte, func([]byte) (A7, error)),
	a8 func() (common.Address, []byte, func([]byte) (A8, error)),
	a9 func() (common.Address, []byte, func([]byte) (A9, error)),
	a10 func() (common.Address, []byte, func([]byte) (A10, error)),
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, rr7 *A7, rr8 *A8, rr9 *A9, rr10 *A10, _ error) {
	z1, r1 := prepareMultiCallArg(a1)
	z2, r2 := prepareMultiCallArg(a2)
	z3, r3 := prepareMultiCallArg(a3)
	z4, r4 := prepareMultiCallArg(a4)
	z5, r5 := prepareMultiCallArg(a5)
	z6, r6 := prepareMultiCallArg(a6)
	z7, r7 := prepareMultiCallArg(a7)
	z8, r8 := prepareMultiCallArg(a8)
	z9, r9 := prepareMultiCallArg(a9)
	z10, r10 := prepareMultiCallArg(a10)
	var results = make([]any, 10)
	err := callN(client, defSlice(z1, z2, z3, z4, z5, z6, z7, z8, z9, z10), defSlice(r1, r2, r3, r4, r5, r6, r7, r8, r9, r10), results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	if !lo.IsNil(results[0]) {
		rr1 = lo.ToPtr(results[0].(A1))
	}
	if !lo.IsNil(results[1]) {
		rr2 = lo.ToPtr(results[1].(A2))
	}
	if !lo.IsNil(results[2]) {
		rr3 = lo.ToPtr(results[2].(A3))
	}
	if !lo.IsNil(results[3]) {
		rr4 = lo.ToPtr(results[3].(A4))
	}
	if !lo.IsNil(results[4]) {
		rr5 = lo.ToPtr(results[4].(A5))
	}
	if !lo.IsNil(results[5]) {
		rr6 = lo.ToPtr(results[5].(A6))
	}
	if !lo.IsNil(results[6]) {
		rr7 = lo.ToPtr(results[6].(A7))
	}
	if !lo.IsNil(results[7]) {
		rr8 = lo.ToPtr(results[7].(A8))
	}
	if !lo.IsNil(results[8]) {
		rr9 = lo.ToPtr(results[8].(A9))
	}
	if !lo.IsNil(results[9]) {
		rr10 = lo.ToPtr(results[9].(A10))
	}
	return
}

func CALLN[Struct any](
	client bind.ContractCaller,
	slices ...func() (common.Address, []byte, func([]byte) (any, error)),
) (*Struct, error) {
	var args []Multicall3Call3
	var functions []func([]byte) (any, error)
	for _, item := range slices {
		a, r := prepareMultiCallArg(item)
		args = append(args, a)
		functions = append(functions, r)
	}

	var results = make([]any, len(slices))
	err := callN(client, args, functions, results)
	if err != nil {
		return nil, err
	}

	var out Struct
	v := reflect.ValueOf(&out).Elem()
	if v.Kind() != reflect.Struct {
		return nil, errors.New("generic type Struct must be a struct type")
	}
	if v.NumField() != len(results) {
		return nil, errors.New("field count of struct does not match number of results")
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			return nil, errors.New("cannot set field " + v.Type().Field(i).Name)
		}
		result := results[i]
		rv := reflect.ValueOf(result)
		if !rv.IsValid() { // result == nil
			field.Set(reflect.Zero(field.Type()))
			continue
		}

		fieldType := field.Type()

		// if the field is a pointer type
		if fieldType.Kind() == reflect.Ptr {
			// 1) Result type if it can be assigned directly to a pointer field (usually of the same pointer type), assign it directly
			if rv.Type().AssignableTo(fieldType) {
				field.Set(rv)
				continue
			}
			// 2) The result is non-pointer, try to assign to the type pointed to by the pointer
			elemType := fieldType.Elem()
			if rv.Type().AssignableTo(elemType) {
				ptr := reflect.New(elemType)
				ptr.Elem().Set(rv)
				field.Set(ptr)
				continue
			} else if rv.Type().ConvertibleTo(elemType) {
				ptr := reflect.New(elemType)
				ptr.Elem().Set(rv.Convert(elemType))
				field.Set(ptr)
				continue
			}
			// 3) Other cases, e.g. if the pointer does not match, throw an error
			return nil, fmt.Errorf("cannot assign result type %s to struct pointer field: %s (%s)",
				rv.Type(), v.Type().Field(i).Name, fieldType)
		}

		if rv.Type().AssignableTo(fieldType) {
			field.Set(rv)
		} else if rv.Type().ConvertibleTo(fieldType) {
			field.Set(rv.Convert(fieldType))
		} else {
			return nil, fmt.Errorf("cannot assign result type %s to struct field: %s (%s)", rv.Type(), v.Type().Field(i).Name, fieldType)
		}
	}

	return &out, nil
}

var (
	_cacheABI *abi.ABI
	_once     sync.Once
)

func getMultiABI() *abi.ABI {
	_once.Do(func() {
		var err error
		_cacheABI, err = contracts_pack.MulticallMetaData.ParseABI()
		if err != nil {
			panic(err)
		}
	})
	return _cacheABI
}

func callN(
	client bind.ContractCaller,
	args []Multicall3Call3,
	functions []func([]byte) (any, error),
	returns []any,
) error {
	return callN1(nil, "aggregate3", client, args, functions, returns)
}

func callN1(
	opts *bind.CallOpts,
	method string,
	client bind.ContractCaller,
	args []Multicall3Call3,
	functions []func([]byte) (any, error),
	returns []any,
) error {
	if len(args) != len(returns) || len(args) != len(functions) {
		panic("invalid arguments")
	}

	var outputs []any
	caller := bind.NewBoundContract(Address, *getMultiABI(), client, nil, nil)
	if err := caller.Call(opts, &outputs, method, args); err != nil {
		return err
	}

	out0 := *abi.ConvertType(outputs[0], new([]Multicall3Result)).(*[]Multicall3Result)
	for idx := range returns {
		ele, err := functions[idx](out0[idx].ReturnData)
		if out0[idx].Success && len(out0[idx].ReturnData) > 0 {
			if err != nil {
				return err
			}
			returns[idx] = ele
		} else {
			returns[idx] = makeNilPtr(reflect.TypeOf(ele))
		}
	}
	return nil
}

func makeNilPtr(typ reflect.Type) interface{} {
	if typ.Kind() == reflect.Ptr {
		return reflect.Zero(typ).Interface()
	}
	ptrType := reflect.PointerTo(typ)
	return reflect.Zero(ptrType).Interface()
}

func prepareMultiCallArg[T any](input Func1[T]) (Multicall3Call3, func([]byte) (any, error)) {
	a1, a2, a3 := input()
	return Multicall3Call3{
			Target:       a1,
			AllowFailure: true,
			CallData:     a2,
		}, func(bytes []byte) (any, error) {
			return a3(bytes)
		}
}

func defSlice[T any](items ...T) []T {
	return items
}
