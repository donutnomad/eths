package multiread

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"github.com/donutnomad/eths/contracts_pack"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Multicall3Call3 = contracts_pack.Multicall3Call3
type Multicall3Result = contracts_pack.Multicall3Result

// Address is the default Multicall3 address: https://www.multicall3.com/abi#ethers-js
var Address = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")

var (
	addressMu  sync.RWMutex
	addressMap = map[uint64]common.Address{}
)

// RegisterAddress registers a Multicall3 address for a specific chain ID.
func RegisterAddress(chainID uint64, addr common.Address) {
	addressMu.Lock()
	addressMap[chainID] = addr
	addressMu.Unlock()
}

// GetAddress returns the Multicall3 address for a specific chain ID.
// If no address is registered, it returns the default Address.
func GetAddress(chainID uint64) common.Address {
	addressMu.RLock()
	addr, found := addressMap[chainID]
	addressMu.RUnlock()
	if found {
		return addr
	}
	return Address
}

// getAddress resolves the Multicall3 address for the given client.
// If the client implements ethereum.ChainIDReader and the chain ID has a registered address, use it.
// Otherwise, fall back to the default Address.
func getAddress(client bind.ContractCaller) common.Address {
	if cr, ok := client.(ethereum.ChainIDReader); ok {
		if chainID, err := cr.ChainID(context.Background()); err == nil {
			return GetAddress(chainID.Uint64())
		}
	}
	return Address
}

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

func (t Func1[T]) prepareMultiCallArg() (Multicall3Call3, ReturnUnPackFunc[any]) {
	a1, a2, a3 := t()
	return Multicall3Call3{
			Target:       a1,
			AllowFailure: true,
			CallData:     a2,
		}, func(bytes []byte) (any, error) {
			return a3(bytes)
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

func One2(contractAddress common.Address, callData []byte) Func2 {
	return func() (common.Address, []byte) {
		return contractAddress, callData
	}
}

func Any[T any](contractAddress common.Address, callData []byte, returnUnpack ReturnUnPackFunc[T]) Func1[any] {
	return func() (common.Address, []byte, func([]byte) (any, error)) {
		return contractAddress, callData, func(bytes []byte) (any, error) {
			return returnUnpack(bytes)
		}
	}
}

func addrOrDefault(addr []common.Address) common.Address {
	if len(addr) > 0 {
		return addr[0]
	}
	return Address
}

func GetChainID(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetChainId(), multiCallPack.UnpackGetChainId)
}

func GetBaseFee(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetBasefee(), multiCallPack.UnpackGetBasefee)
}

func GetBlockNumber(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetBlockNumber(), multiCallPack.UnpackGetBlockNumber)
}

func GetCurrentBlockTimestamp(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetCurrentBlockTimestamp(), multiCallPack.UnpackGetCurrentBlockTimestamp)
}

func GetCurrentBlockGasLimit(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetCurrentBlockGasLimit(), multiCallPack.UnpackGetCurrentBlockGasLimit)
}

func GetCurrentBlockDifficulty(addr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(addr), multiCallPack.PackGetCurrentBlockDifficulty(), multiCallPack.UnpackGetCurrentBlockDifficulty)
}

func GetCurrentBlockCoinbase(addr ...common.Address) Func1[common.Address] {
	return One(addrOrDefault(addr), multiCallPack.PackGetCurrentBlockCoinbase(), multiCallPack.UnpackGetCurrentBlockCoinbase)
}

func GetEthBalance(ethAddr common.Address, multicallAddr ...common.Address) Func1[*big.Int] {
	return One(addrOrDefault(multicallAddr), multiCallPack.PackGetEthBalance(ethAddr), multiCallPack.UnpackGetEthBalance)
}

func GetBlockHash(blockNumber *big.Int, addr ...common.Address) Func1[common.Hash] {
	return One(addrOrDefault(addr), multiCallPack.PackGetBlockHash(blockNumber), func(bytes []byte) (common.Hash, error) {
		ret, err := multiCallPack.UnpackGetBlockHash(bytes)
		if err != nil {
			return common.Hash{}, err
		}
		return ret, nil
	})
}

func GetLastBlockHash(addr ...common.Address) Func1[common.Hash] {
	return One(addrOrDefault(addr), multiCallPack.PackGetLastBlockHash(), func(bytes []byte) (common.Hash, error) {
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

func CALLSlice[A1 any](
	client bind.ContractCaller,
	unpack func([]byte) (A1, error),
	inputs ...Func2,
) ([]*A1, error) {
	if len(inputs) == 0 {
		panic("[multiread] invalid inputs")
	}
	var args []Multicall3Call3
	var functions []ReturnUnPackFunc[any]
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

// CALL (direct) without multicall
func CALL[A1 any](client bind.ContractCaller, a1 Func1[A1]) (*A1, error) {
	target, calldata, unpack := a1()
	caller := bind.NewBoundContract(target, abi.ABI{}, client, nil, nil)
	response, err := caller.CallRaw(&bind.CallOpts{}, calldata)
	if err != nil {
		return nil, err
	}
	ret, err := unpack(response)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func CALLD[A1 any](client bind.ContractCaller, target common.Address, calldata []byte, unpack func([]byte) (A1, error)) (*A1, error) {
	caller := bind.NewBoundContract(target, abi.ABI{}, client, nil, nil)
	response, err := caller.CallRaw(&bind.CallOpts{}, calldata)
	if err != nil {
		return nil, err
	}
	ret, err := unpack(response)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func CALL1[A1 any](client bind.ContractCaller, a1 Func1[A1]) (r1 *A1, err error) {
	r1, _, _, _, _, _, _, _, _, _, err = callGeneric(
		client, a1, fnil, fnil, fnil, fnil, fnil, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL2[A1 any, A2 any](
	client bind.ContractCaller, a1 Func1[A1], a2 Func1[A2],
) (r1 *A1, r2 *A2, err error) {
	r1, r2, _, _, _, _, _, _, _, _, err = callGeneric(
		client, a1, a2, fnil, fnil, fnil, fnil, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL3[A1 any, A2 any, A3 any](
	client bind.ContractCaller, a1 Func1[A1], a2 Func1[A2], a3 Func1[A3],
) (r1 *A1, r2 *A2, r3 *A3, err error) {
	r1, r2, r3, _, _, _, _, _, _, _, err = callGeneric(
		client, a1, a2, a3, fnil, fnil, fnil, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL4[A1 any, A2 any, A3 any, A4 any](
	client bind.ContractCaller, a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, err error) {
	r1, r2, r3, r4, _, _, _, _, _, _, err = callGeneric(
		client, a1, a2, a3, a4, fnil, fnil, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL5[A1 any, A2 any, A3 any, A4 any, A5 any](
	client bind.ContractCaller, a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, err error) {
	r1, r2, r3, r4, r5, _, _, _, _, _, err = callGeneric(
		client, a1, a2, a3, a4, a5, fnil, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL6[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5], a6 Func1[A6],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, r6 *A6, err error) {
	r1, r2, r3, r4, r5, r6, _, _, _, _, err = callGeneric(
		client, a1, a2, a3, a4, a5, a6, fnil, fnil, fnil, fnil,
	)
	return
}

func CALL7[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5], a6 Func1[A6], a7 Func1[A7],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, r6 *A6, r7 *A7, err error) {
	r1, r2, r3, r4, r5, r6, r7, _, _, _, err = callGeneric(
		client, a1, a2, a3, a4, a5, a6, a7, fnil, fnil, fnil,
	)
	return
}

func CALL8[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5], a6 Func1[A6], a7 Func1[A7], a8 Func1[A8],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, r6 *A6, r7 *A7, r8 *A8, err error) {
	r1, r2, r3, r4, r5, r6, r7, r8, _, _, err = callGeneric(
		client, a1, a2, a3, a4, a5, a6, a7, a8, fnil, fnil,
	)
	return
}

func CALL9[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, A9 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5], a6 Func1[A6], a7 Func1[A7], a8 Func1[A8], a9 Func1[A9],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, r6 *A6, r7 *A7, r8 *A8, r9 *A9, err error) {
	r1, r2, r3, r4, r5, r6, r7, r8, r9, _, err = callGeneric(
		client, a1, a2, a3, a4, a5, a6, a7, a8, a9, fnil,
	)
	return
}

func CALL10[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, A9 any, A10 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5],
	a6 Func1[A6], a7 Func1[A7], a8 Func1[A8], a9 Func1[A9], a10 Func1[A10],
) (r1 *A1, r2 *A2, r3 *A3, r4 *A4, r5 *A5, r6 *A6, r7 *A7, r8 *A8, r9 *A9, r10 *A10, err error) {
	return callGeneric(client, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
}

func CALLN[Struct any](
	client bind.ContractCaller,
	slices ...Func1[any],
) (*Struct, error) {
	var args []Multicall3Call3
	var functions []ReturnUnPackFunc[any]
	for _, item := range slices {
		a, r := item.prepareMultiCallArg()
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
		return nil, errors.New("[multiread] generic type Struct must be a struct type")
	}
	if v.NumField() != len(results) {
		return nil, errors.New("[multiread] field count of struct does not match number of results")
	}
	for i := range v.NumField() {
		field := v.Field(i)
		if !field.CanSet() {
			return nil, errors.New("[multiread] cannot set field " + v.Type().Field(i).Name)
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
			return nil, fmt.Errorf("[multiread] cannot assign result type %s to struct pointer field: %s (%s)",
				rv.Type(), v.Type().Field(i).Name, fieldType)
		}

		if rv.Type().AssignableTo(fieldType) {
			field.Set(rv)
		} else if rv.Type().ConvertibleTo(fieldType) {
			field.Set(rv.Convert(fieldType))
		} else {
			return nil, fmt.Errorf("[multiread] cannot assign result type %s to struct field: %s (%s)", rv.Type(), v.Type().Field(i).Name, fieldType)
		}
	}

	return &out, nil
}

// CALLNE is an enhanced version of CALLN that can parse array fields in structs
// It can handle array fields of fixed size or slices and assign values based on type compatibility
//
//	type Example struct {
//		Bools1      [2]*bool
//		Bools2      [2]*bool
//		Addresses   [1]*common.Address
//		Addresses2  []*common.Address
//		BlockNumber *big.Int
//		BlockTime   *big.Int
//	}
func CALLNE[Struct any](
	client bind.ContractCaller,
	slices ...Func1[any],
) (*Struct, error) {
	results, err := prepareMultiCallArgs(client, slices)
	if err != nil {
		return nil, err
	}
	var out Struct
	v := reflect.ValueOf(&out).Elem()
	if v.Kind() != reflect.Struct {
		return nil, errors.New("[multiread] generic type Struct must be a struct type")
	}

	// Mark which results have been assigned
	assigned := make([]bool, len(results))

	// Step 1: Process non-array fields
	assignNonArrayFields(v, results, assigned)

	// Step 2: Process array and slice fields (by type matching)
	assignArrayFields(v, results, assigned)

	// Check unassigned results
	for i, wasAssigned := range assigned {
		if !wasAssigned && !lo.IsNil(results[i]) {
			return nil, fmt.Errorf("[multiread] result at index %d of type %v could not be assigned to any field in the struct",
				i, reflect.TypeOf(results[i]))
		}
	}

	return &out, nil
}

var fnil Func1[any] = nil

func callGeneric[A1 any, A2 any, A3 any, A4 any, A5 any, A6 any, A7 any, A8 any, A9 any, A10 any](
	client bind.ContractCaller,
	a1 Func1[A1], a2 Func1[A2], a3 Func1[A3], a4 Func1[A4], a5 Func1[A5],
	a6 Func1[A6], a7 Func1[A7], a8 Func1[A8], a9 Func1[A9], a10 Func1[A10],
) (rr1 *A1, rr2 *A2, rr3 *A3, rr4 *A4, rr5 *A5, rr6 *A6, rr7 *A7, rr8 *A8, rr9 *A9, rr10 *A10, _ error) {
	var args []Multicall3Call3
	var functions []ReturnUnPackFunc[any]
	var count int

	if a1 != nil {
		z1, r1 := a1.prepareMultiCallArg()
		args = append(args, z1)
		functions = append(functions, r1)
		count++
	}
	if a2 != nil {
		z2, r2 := a2.prepareMultiCallArg()
		args = append(args, z2)
		functions = append(functions, r2)
		count++
	}
	if a3 != nil {
		z3, r3 := a3.prepareMultiCallArg()
		args = append(args, z3)
		functions = append(functions, r3)
		count++
	}
	if a4 != nil {
		z4, r4 := a4.prepareMultiCallArg()
		args = append(args, z4)
		functions = append(functions, r4)
		count++
	}
	if a5 != nil {
		z5, r5 := a5.prepareMultiCallArg()
		args = append(args, z5)
		functions = append(functions, r5)
		count++
	}
	if a6 != nil {
		z6, r6 := a6.prepareMultiCallArg()
		args = append(args, z6)
		functions = append(functions, r6)
		count++
	}
	if a7 != nil {
		z7, r7 := a7.prepareMultiCallArg()
		args = append(args, z7)
		functions = append(functions, r7)
		count++
	}
	if a8 != nil {
		z8, r8 := a8.prepareMultiCallArg()
		args = append(args, z8)
		functions = append(functions, r8)
		count++
	}
	if a9 != nil {
		z9, r9 := a9.prepareMultiCallArg()
		args = append(args, z9)
		functions = append(functions, r9)
		count++
	}
	if a10 != nil {
		z10, r10 := a10.prepareMultiCallArg()
		args = append(args, z10)
		functions = append(functions, r10)
		count++
	}

	if count == 0 {
		panic("[multiread] no valid inputs")
	}

	var results = make([]any, count)
	err := callN(client, args, functions, results)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var resultIndex int
	if a1 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr1 = lo.ToPtr(results[resultIndex].(A1))
		}
		resultIndex++
	}
	if a2 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr2 = lo.ToPtr(results[resultIndex].(A2))
		}
		resultIndex++
	}
	if a3 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr3 = lo.ToPtr(results[resultIndex].(A3))
		}
		resultIndex++
	}
	if a4 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr4 = lo.ToPtr(results[resultIndex].(A4))
		}
		resultIndex++
	}
	if a5 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr5 = lo.ToPtr(results[resultIndex].(A5))
		}
		resultIndex++
	}
	if a6 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr6 = lo.ToPtr(results[resultIndex].(A6))
		}
		resultIndex++
	}
	if a7 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr7 = lo.ToPtr(results[resultIndex].(A7))
		}
		resultIndex++
	}
	if a8 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr8 = lo.ToPtr(results[resultIndex].(A8))
		}
		resultIndex++
	}
	if a9 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr9 = lo.ToPtr(results[resultIndex].(A9))
		}
		resultIndex++
	}
	if a10 != nil {
		if !lo.IsNil(results[resultIndex]) {
			rr10 = lo.ToPtr(results[resultIndex].(A10))
		}
	}

	return
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
	functions []ReturnUnPackFunc[any],
	returns []any,
) error {
	return callN1(nil, "aggregate3", client, args, functions, returns)
}

func callN1(
	opts *bind.CallOpts,
	method string,
	client bind.ContractCaller,
	args []Multicall3Call3,
	functions []ReturnUnPackFunc[any],
	returns []any,
) error {
	if len(args) != len(returns) || len(args) != len(functions) {
		panic("[multiread] invalid arguments")
	}
	multicallAddr := getAddress(client)

	var outputs []any
	caller := bind.NewBoundContract(multicallAddr, *getMultiABI(), client, nil, nil)
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

// Helper function: Prepare multi-call arguments
func prepareMultiCallArgs(client bind.ContractCaller, slices []Func1[any]) ([]any, error) {
	var args []Multicall3Call3
	var functions []ReturnUnPackFunc[any]
	for _, item := range slices {
		a, r := item.prepareMultiCallArg()
		args = append(args, a)
		functions = append(functions, r)
	}

	var results = make([]any, len(slices))
	err := callN(client, args, functions, results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Helper function: Check if result can be assigned to field
func isAssignableTo(result any, field reflect.Value) (reflect.Value, bool) {
	if lo.IsNil(result) {
		return reflect.Value{}, false
	}

	rv := reflect.ValueOf(result)
	fieldType := field.Type()

	// Directly assignable
	if rv.Type().AssignableTo(fieldType) {
		return rv, true
	}

	// Handle pointer type
	if fieldType.Kind() == reflect.Ptr {
		elemType := fieldType.Elem()
		if rv.Type().AssignableTo(elemType) {
			ptr := reflect.New(elemType)
			ptr.Elem().Set(rv)
			return ptr, true
		} else if rv.Type().ConvertibleTo(elemType) {
			ptr := reflect.New(elemType)
			ptr.Elem().Set(rv.Convert(elemType))
			return ptr, true
		}
	} else if rv.Type().ConvertibleTo(fieldType) {
		return rv.Convert(fieldType), true
	}

	return reflect.Value{}, false
}

// Helper function: Process non-array fields
func assignNonArrayFields(v reflect.Value, results []any, assigned []bool) {
	for i := range v.NumField() {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := field.Type()
		// Skip array fields
		if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			continue
		}

		// For regular fields, try to find matching results
		for j := range results {
			if assigned[j] {
				continue
			}

			if value, ok := isAssignableTo(results[j], field); ok {
				field.Set(value)
				assigned[j] = true
				break
			}
		}
	}
}

// Helper function: Process array and slice fields
func assignArrayFields(v reflect.Value, results []any, assigned []bool) {
	// Step 1: Collect all array and slice fields with their type information
	type ArrayFieldInfo struct {
		Field        reflect.Value
		FieldName    string
		ElementType  reflect.Type // Array element type
		IsFixedArray bool         // Whether it's a fixed-size array
		MaxSize      int          // Maximum capacity (fixed array)
	}

	var arrayFields []ArrayFieldInfo

	for i := range v.NumField() {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name
		fieldType := field.Type()

		if !field.CanSet() {
			continue
		}

		// Process array and slice fields
		if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			elemType := fieldType.Elem() // Get element type

			info := ArrayFieldInfo{
				Field:        field,
				FieldName:    fieldName,
				ElementType:  elemType,
				IsFixedArray: fieldType.Kind() == reflect.Array,
				MaxSize:      0,
			}

			if info.IsFixedArray {
				info.MaxSize = fieldType.Len()
			}

			arrayFields = append(arrayFields, info)
		}
	}

	// Step 2: Assign all unassigned results to appropriate array fields by type
	for _, info := range arrayFields {
		elemType := info.ElementType
		field := info.Field

		// Find matching results
		var matchedResults []any
		var matchedIndices []int

		for i := range results {
			if assigned[i] || lo.IsNil(results[i]) {
				continue
			}

			result := results[i]
			resultVal := reflect.ValueOf(result)
			resultType := resultVal.Type()

			// Check type compatibility
			canAssign := false

			// Direct match
			if resultType.AssignableTo(elemType) {
				canAssign = true
			} else if elemType.Kind() == reflect.Ptr {
				// Element is pointer type
				ptrElemType := elemType.Elem()
				if resultType.AssignableTo(ptrElemType) || resultType.ConvertibleTo(ptrElemType) {
					canAssign = true
				}
			} else if resultType.ConvertibleTo(elemType) {
				canAssign = true
			}

			if canAssign {
				matchedResults = append(matchedResults, result)
				matchedIndices = append(matchedIndices, i)

				// If fixed size array and already full, stop
				if info.IsFixedArray && len(matchedResults) >= info.MaxSize {
					break
				}
			}
		}

		// Assign matched results to array/slice
		if len(matchedResults) > 0 {
			if info.IsFixedArray {
				// Fixed size array
				arrayLen := info.MaxSize
				for i := 0; i < len(matchedResults) && i < arrayLen; i++ {
					result := matchedResults[i]
					resultVal := reflect.ValueOf(result)
					resultType := resultVal.Type()

					// Set array element
					setArrayElement(field.Index(i), elemType, resultVal, resultType)

					// Mark as assigned
					assigned[matchedIndices[i]] = true
				}
			} else {
				// Dynamic slice
				sliceValue := reflect.MakeSlice(field.Type(), len(matchedResults), len(matchedResults))
				for i := 0; i < len(matchedResults); i++ {
					result := matchedResults[i]
					resultVal := reflect.ValueOf(result)
					resultType := resultVal.Type()

					// Set slice element
					setArrayElement(sliceValue.Index(i), elemType, resultVal, resultType)

					// Mark as assigned
					assigned[matchedIndices[i]] = true
				}
				field.Set(sliceValue)
			}
		}
	}
}

// Helper function: Set array/slice element
func setArrayElement(dest reflect.Value, destType reflect.Type, src reflect.Value, srcType reflect.Type) {
	// Directly assignable
	if srcType.AssignableTo(destType) {
		dest.Set(src)
		return
	}

	// Destination is pointer type
	if destType.Kind() == reflect.Ptr {
		ptrElemType := destType.Elem()

		// Create new pointer
		ptr := reflect.New(ptrElemType)

		// Source can be assigned to pointer's element type
		if srcType.AssignableTo(ptrElemType) {
			ptr.Elem().Set(src)
		} else if srcType.ConvertibleTo(ptrElemType) {
			// Source can be converted to pointer's element type
			ptr.Elem().Set(src.Convert(ptrElemType))
		}

		dest.Set(ptr)
	} else if srcType.ConvertibleTo(destType) {
		// Source can be converted to destination type
		dest.Set(src.Convert(destType))
	}
}
