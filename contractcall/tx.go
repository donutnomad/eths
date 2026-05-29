package contractcall

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"unsafe"

	"slices"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethtype"
	"github.com/donutnomad/eths/hexutil"
	"github.com/donutnomad/eths/rlp"
	"github.com/holiman/uint256"
	"github.com/samber/lo"
)

var UNREACHABLE = "unreachable"
var ErrInvalidSig = errors.New("invalid transaction v, r, s values")

// Deprecated: Use Tx
type TxWrapper = Tx

type Tx = txImpl

type AccessListSetter interface {
	AccessList() ethtype.AccessList
	SetAccessList(accessList ethtype.AccessList) bool
}

type BaseTx interface {
	To() *common.Address
	SetTo(to *common.Address)
	Value() *big.Int
	SetValue(value *big.Int)
	Data() []byte
	SetData(data []byte)
	Nonce() uint64
	SetNonce(nonce uint64)
	Gas() uint64
	SetGas(gas uint64)
	ChainID() *big.Int
	SetChainID(chainID *big.Int)
	Signature() (v, r, s *big.Int)
	SetSignature(v, r, s *big.Int)
}

// //// OLDER //////////
// IEIP155 Legacy Tx (Post replay protection)
type IEIP155 interface {
	BaseTx
	GasPrice() *big.Int
	SetGasPrice(gasPrice *big.Int) bool
}

// IEIP2930 Access List Tx
type IEIP2930 interface {
	IEIP155
	AccessListSetter
}

// //// NEWER //////////

// IEIP1559 Dynamic Fee Tx
type IEIP1559 interface {
	BaseTx
	AccessListSetter
	MaxFeePerGas() *big.Int
	MaxPriorityFeePerGas() *big.Int
	SetMaxFeePerGas(maxFeePerGas *big.Int) bool
	SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) bool
}

// IEIP4844 Blob Tx
type IEIP4844 interface {
	IEIP1559
	MaxFeePerBlobGas() *big.Int
	BlobHashes() []common.Hash
	SetMaxFeePerBlobGas(maxFeePerBlobGas *big.Int) bool
	SetBlobHashes(blobHashes []common.Hash) bool
	Sidecar() *ethtype.BlobTxSidecar
	SetSidecar(sidecar *ethtype.BlobTxSidecar)
}

// IEIP7702 Set code
type IEIP7702 interface {
	IEIP1559
	AuthList() []ethtype.SetCodeAuthorization
	SetAuthList(authList []ethtype.SetCodeAuthorization) bool
}

type ITx interface {
	TxType() TxType
	IEIP155
	IEIP1559
	IEIP2930
	IEIP4844
	IEIP7702
	Sign(privateKey ISigner) error
	ToJSON() []byte
	ToTransaction() *ethtype.ETransaction
	Hash() common.Hash
	SigHash() common.Hash
	json.Marshaler
	json.Unmarshaler
	encoding.BinaryUnmarshaler
	encoding.BinaryMarshaler
}

var _ ITx = &txImpl{}

// txImpl [ethtype.DynamicFeeTx] [ethtype.BlobTx] [ethtype.AccessListTx] [ethtype.LegacyTx] [ethtype.SetCodeTx]
type txImpl struct {
	chainID *uint256.Int // safe, not nil
	nonce   uint64
	gas     uint64
	to      *common.Address
	value   *uint256.Int // safe, not nil
	data    []byte

	maxPriorityFeePerGas *uint256.Int       // notInclude(eip-155(legacy),eip-2930(access))
	maxFeePerGas         *uint256.Int       // notInclude(eip-155(legacy),eip-2930(access))
	gasPrice             *uint256.Int       // eip-155(legacy), eip-2930(access)
	accessList           ethtype.AccessList // notInclude(eip-155(legacy))
	maxFeePerBlobGas     *uint256.Int       // eip-4844(blob)
	blobHashes           []common.Hash      // eip-4844(blob)
	// A blob transaction can optionally contain blobs. This field must be set when BlobTx
	// is used to create a transaction for signing.
	sidecar  *ethtype.BlobTxSidecar         `rlp:"-"` // eip-4844(blob)
	authList []ethtype.SetCodeAuthorization // eip-7702(auth)

	v      [32]byte // safe, default 0， yParity: 0/1, EIP155: (0/1 + 35) + chainID * 2
	r      [32]byte // safe, default 0
	s      [32]byte // safe, default 0
	txType TxType
}

func NewTxWith[
	ChainID *big.Int | *uint256.Int | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64,
](txType TxType, chainID ChainID) ITx {
	return &txImpl{
		chainID: newInt(bigIntOrIntToBigInt(chainID)),
		txType:  txType,
		value:   newIntBy(0),
	}
}

func NewTx[
	T *ethtype.LegacyTx | *ethtype.AccessListTx | *ethtype.DynamicFeeTx | *ethtype.BlobTx | *ethtype.SetCodeTx | *ethtype.ETransaction,
	ChainID *big.Int | *uint256.Int | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64,
](tx T, chainID ChainID) ITx {
	ret, err := newTxImpl(any(tx), bigIntOrIntToBigInt(chainID))
	if err != nil {
		panic(err)
	}
	return ret
}

// NewTxImpl
// Deprecated: Use NewTx
func NewTxImpl[
	T *ethtype.LegacyTx | *ethtype.AccessListTx | *ethtype.DynamicFeeTx | *ethtype.BlobTx | *ethtype.SetCodeTx | *ethtype.ETransaction,
	ChainID *big.Int | *uint256.Int | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64,
](tx T, chainID ChainID) ITx {
	return NewTx(tx, chainID)
}

func (t *txImpl) isModern() bool {
	return t.txType != LegacyTxType
}

func (t *txImpl) isLegacy() bool {
	return !t.isModern()
}

func (t *txImpl) Hash() common.Hash {
	var prefix = ifG(t.isModern(), byte(t.txType))
	return PrefixedRlpHash(prefix, t.BuildRlpFields(false))
}

func (t *txImpl) SigHash() common.Hash {
	var prefix = ifG(t.isModern(), byte(t.txType))
	return PrefixedRlpHash(prefix, t.BuildRlpFields(true))
}

func (t *txImpl) Sign(privateKey ISigner) error {
	sig, err := privateKey.Sign(t.SigHash().Bytes())
	if err != nil {
		return err
	}
	if sig == nil {
		return nil
	}

	t.SetSignature(
		computeVForEIP155(sig.V(), t.ChainID(), t.isLegacy()),
		sig.R(), sig.S(),
	)
	return nil
}

func (t *txImpl) BuildRlpFields(forSignature bool) []any {
	return buildArgs(
		if_(t.isModern(), t.chainID),
		t.Nonce(),
		ifElse(t.txType.IsEIP1559Gas(), []any{t.maxPriorityFeePerGas, t.maxFeePerGas}, t.gasPrice),
		t.Gas(),
		t.To(),
		t.Value(),
		t.Data(),
		if_(t.isModern(), t.AccessList()),
		if_(t.txType == BlobTxType, t.maxFeePerBlobGas, t.blobHashes),
		if_(t.txType == SetCodeTxType, t.AuthList()),
		if_(forSignature && t.isLegacy(), t.chainID, uint(0), uint(0)),
		if_(!forSignature, t.getV(), t.getR(), t.getS()),
	)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (t *txImpl) MarshalBinary() ([]byte, error) {
	inner := t.Export()
	if t.isLegacy() {
		return rlp.EncodeToBytes(inner)
	}
	var err error

	var buf bytes.Buffer
	buf.WriteByte(byte(t.txType))

	switch t.txType {
	case AccessListTxType, DynamicFeeTxType, SetCodeTxType:
		err = rlp.Encode(&buf, inner)
	case BlobTxType:
		if t.sidecar == nil {
			err = rlp.Encode(&buf, t)
		} else {
			return t.ToTransaction().MarshalBinary()
		}
	default:
		panic(UNREACHABLE)
	}
	return buf.Bytes(), err
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP-2718 typed transactions.
func (t *txImpl) UnmarshalBinary(b []byte) error {
	tx := new(ethtype.ETransaction)
	err := tx.UnmarshalBinary(b)
	if err != nil {
		return err
	}
	ret, err := newTxImpl(any(tx), tx.ChainId())
	if err != nil {
		return err
	}
	*t = *ret
	return nil
}

func (t *txImpl) ToJSON() []byte {
	return lo.Must1(t.MarshalJSON())
}

func (t *txImpl) ToTransaction() *ethtype.ETransaction {
	return ethtype.NewTx(t.Export())
}

func (t *txImpl) Signature() (v, r, s *big.Int) {
	return t.getV().ToBig(), t.getR().ToBig(), t.getS().ToBig()
}

func (t *txImpl) SetSignature(v, r, s *big.Int) {
	t.v = bigToBytes32(v) // 0 or 1 (not 27/28)
	t.r = bigToBytes32(r)
	t.s = bigToBytes32(s)
}

func (t *txImpl) Nonce() uint64 {
	return t.nonce
}

func (t *txImpl) SetNonce(nonce uint64) {
	t.nonce = nonce
}

func (t *txImpl) To() *common.Address {
	return t.to
}

func (t *txImpl) SetTo(to *common.Address) {
	t.to = to
	if t.txType == BlobTxType || t.txType == SetCodeTxType {
		if t.to == nil {
			t.to = new(common.Address)
		}
	}
}

func (t *txImpl) Gas() uint64 {
	return t.gas
}

func (t *txImpl) SetGas(gas uint64) {
	t.gas = gas
}

func (t *txImpl) Value() *big.Int {
	return copyInt(t.value.ToBig())
}

func (t *txImpl) SetValue(value *big.Int) {
	t.value = newInt(value)
}

func (t *txImpl) Data() []byte {
	return t.data
}

func (t *txImpl) SetData(data []byte) {
	t.data = data
}

func (t *txImpl) TxType() TxType {
	return t.txType
}

func (t *txImpl) ChainID() *big.Int {
	return t.chainID.ToBig()
}

func (t *txImpl) SetChainID(chainID *big.Int) {
	t.chainID = newInt(chainID)
}

// GasPrice eip-155(legacy), eip-2930(access)
func (t *txImpl) GasPrice() *big.Int {
	if t.isLegacy() || t.txType == AccessListTxType {
		return t.gasPrice.ToBig()
	} else {
		return nil
	}
}

// SetGasPrice eip-155(legacy), eip-2930(access)
func (t *txImpl) SetGasPrice(gasPrice *big.Int) bool {
	if t.isLegacy() || t.txType == AccessListTxType {
		t.gasPrice = newInt(gasPrice)
		return true
	} else {
		return false
	}
}

// AccessList notInclude(eip-155(legacy))
func (t *txImpl) AccessList() ethtype.AccessList {
	if t.isModern() {
		return t.accessList
	}
	return nil
}

// SetAccessList notInclude(eip-155(legacy))
func (t *txImpl) SetAccessList(accessList ethtype.AccessList) bool {
	if t.isModern() {
		t.accessList = accessList
		return true
	}
	return false
}

// MaxFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) MaxFeePerGas() *big.Int {
	if t.txType.IsEIP1559Gas() {
		return t.maxFeePerGas.ToBig()
	} else {
		return nil
	}
}

// MaxPriorityFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) MaxPriorityFeePerGas() *big.Int {
	if t.txType.IsEIP1559Gas() {
		return t.maxPriorityFeePerGas.ToBig()
	} else {
		return nil
	}
}

// SetMaxFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) SetMaxFeePerGas(maxFeePerGas *big.Int) bool {
	if t.txType.IsEIP1559Gas() {
		t.maxFeePerGas = newInt(maxFeePerGas)
		return true
	} else {
		return false
	}
}

// SetMaxPriorityFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) bool {
	if t.txType.IsEIP1559Gas() {
		t.maxPriorityFeePerGas = newInt(maxPriorityFeePerGas)
		return true
	} else {
		return false
	}
}

// MaxFeePerBlobGas eip-4844(blob)
func (t *txImpl) MaxFeePerBlobGas() *big.Int {
	if t.txType == BlobTxType {
		return t.maxFeePerBlobGas.ToBig()
	} else {
		return nil
	}
}

// BlobHashes eip-4844(blob)
func (t *txImpl) BlobHashes() []common.Hash {
	if t.txType == BlobTxType {
		return t.blobHashes
	} else {
		return nil
	}
}

// Sidecar eip-4844(blob)
func (t *txImpl) Sidecar() *ethtype.BlobTxSidecar {
	return t.sidecar
}

// SetMaxFeePerBlobGas eip-4844(blob)
func (t *txImpl) SetMaxFeePerBlobGas(maxFeePerBlobGas *big.Int) bool {
	if t.txType == BlobTxType {
		t.maxFeePerBlobGas = newInt(maxFeePerBlobGas)
		return true
	}
	return false
}

// SetBlobHashes eip-4844(blob)
func (t *txImpl) SetBlobHashes(blobHashes []common.Hash) bool {
	if t.txType == BlobTxType {
		t.blobHashes = blobHashes
		return true
	}
	return false
}

// SetSidecar eip-4844(blob)
// A blob transaction can optionally contain blobs. This field must be set when BlobTx
// is used to create a transaction for signing.
func (t *txImpl) SetSidecar(sidecar *ethtype.BlobTxSidecar) {
	t.sidecar = sidecar
}

// AuthList eip-7702(auth)
func (t *txImpl) AuthList() []ethtype.SetCodeAuthorization {
	if t.txType == SetCodeTxType {
		return t.authList
	}
	return nil
}

// SetAuthList eip-7702(auth)
func (t *txImpl) SetAuthList(authList []ethtype.SetCodeAuthorization) bool {
	if t.txType == SetCodeTxType {
		t.authList = slices.Clone(authList)
		return true
	}
	return false
}

func (t *txImpl) getV() *uint256.Int {
	return new(uint256.Int).SetBytes(t.v[:])
}

func (t *txImpl) getR() *uint256.Int {
	return new(uint256.Int).SetBytes(t.r[:])
}

func (t *txImpl) getS() *uint256.Int {
	return new(uint256.Int).SetBytes(t.s[:])
}

var big8 = big.NewInt(8)

// Sender 获取交易发送方, 从签名中恢复，如果没有签名，会报错`ErrInvalidSig`
func (t *txImpl) Sender() (common.Address, error) {
	V, R, S := t.getV().ToBig(), t.getR().ToBig(), t.getS().ToBig()
	if t.isLegacy() {
		t.ToTransaction().Protected()
		if t.isProtected() {
			mulChainID := new(big.Int).Mul(t.ChainID(), big.NewInt(2))
			V = new(big.Int).Sub(V, mulChainID)
			V.Sub(V, big8)
		}
	} else {
		// 'modern' txs are defined to use 0 and 1 as their recovery
		// id, add 27 to become equivalent to unprotected Homestead signatures.
		V = new(big.Int).Add(V, big.NewInt(27))
	}
	return recoverPlain(t.SigHash(), R, S, V, true)
}

func (t *txImpl) UnmarshalJSON(input []byte) error {
	var tx = new(ethtype.ETransaction)
	err := tx.UnmarshalJSON(input)
	if err != nil {
		return err
	}
	var chainID = tx.ChainId()

	v, r, s := tx.RawSignatureValues()
	if tx.Type() == byte(LegacyTxType) && (v.Sign() == 0 || r.Sign() == 0 || s.Sign() == 0) {
		var fields struct {
			ChainID string `json:"chainId"`
		}
		if err := json.Unmarshal(input, &fields); err == nil && fields.ChainID != "" {
			chainIDUint64, err := hexutil.DecodeUint64(fields.ChainID)
			if err == nil {
				chainID = new(big.Int).SetUint64(chainIDUint64)
			}
		}
	}
	ret, err := newTxImpl(tx, chainID)
	if err != nil {
		return err
	}
	*t = *ret
	return nil
}

func (t *txImpl) Export() ethtype.TxData {
	return txImplToTx(t)
}

func newTxImpl(tx any, chainID *big.Int) (*txImpl, error) {
	impl, err := newTxImplRaw(tx, chainID)
	if err != nil {
		return nil, err
	}
	return impl, nil
}

func txImplToTx(t *txImpl) ethtype.TxData {
	switch t.txType {
	case LegacyTxType:
		if t.gasPrice == nil {
			t.gasPrice = newIntBy(0)
		}
		return &ethtype.LegacyTx{
			Nonce:    t.nonce,
			GasPrice: t.gasPrice.ToBig(),
			Gas:      t.gas,
			To:       t.to,
			Value:    t.Value(),
			Data:     t.data,
			V:        t.getV().ToBig(),
			R:        t.getR().ToBig(),
			S:        t.getS().ToBig(),
		}
	case AccessListTxType:
		if t.gasPrice == nil {
			t.gasPrice = newIntBy(0)
		}
		return &ethtype.AccessListTx{
			ChainID:    t.ChainID(),
			Nonce:      t.nonce,
			GasPrice:   t.gasPrice.ToBig(),
			Gas:        t.gas,
			To:         t.to,
			Value:      t.Value(),
			Data:       t.data,
			AccessList: t.accessList,
			V:          t.getV().ToBig(),
			R:          t.getR().ToBig(),
			S:          t.getS().ToBig(),
		}
	case DynamicFeeTxType:
		if t.maxPriorityFeePerGas == nil {
			t.maxPriorityFeePerGas = newIntBy(0)
		}
		if t.maxFeePerGas == nil {
			t.maxFeePerGas = newIntBy(0)
		}
		return &ethtype.DynamicFeeTx{
			ChainID:    t.ChainID(),
			Nonce:      t.nonce,
			GasTipCap:  t.maxPriorityFeePerGas.ToBig(),
			GasFeeCap:  t.maxFeePerGas.ToBig(),
			Gas:        t.gas,
			To:         t.to,
			Value:      t.Value(),
			Data:       t.data,
			AccessList: t.accessList,
			V:          t.getV().ToBig(),
			R:          t.getR().ToBig(),
			S:          t.getS().ToBig(),
		}
	case BlobTxType:
		return &ethtype.BlobTx{
			ChainID:    t.chainID,
			Nonce:      t.nonce,
			GasTipCap:  t.maxPriorityFeePerGas,
			GasFeeCap:  t.maxFeePerGas,
			Gas:        t.gas,
			To:         *t.to,
			Value:      t.value,
			Data:       t.data,
			AccessList: t.accessList,
			BlobFeeCap: t.maxFeePerBlobGas,
			BlobHashes: t.blobHashes,
			Sidecar:    nil, //TODO: add sidecar
			V:          t.getV(),
			R:          t.getR(),
			S:          t.getS(),
		}
	case SetCodeTxType:
		return &ethtype.SetCodeTx{
			ChainID:    t.chainID,
			Nonce:      t.nonce,
			GasFeeCap:  t.maxFeePerGas,
			GasTipCap:  t.maxPriorityFeePerGas,
			Gas:        t.gas,
			Value:      t.value,
			Data:       t.data,
			To:         *t.to,
			AccessList: t.accessList,
			AuthList:   t.authList,
			V:          t.getV(),
			R:          t.getR(),
			S:          t.getS(),
		}
	default:
		panic(UNREACHABLE)
	}
}

func newTxImplRaw(tx any, chainID *big.Int) (*txImpl, error) {
	switch v := tx.(type) {
	case *ethtype.LegacyTx:
		return &txImpl{
			chainID:              newInt(chainID),
			nonce:                v.Nonce,
			gas:                  v.Gas,
			to:                   v.To,
			value:                newInt(v.Value),
			data:                 v.Data,
			txType:               LegacyTxType,
			maxPriorityFeePerGas: nil,
			maxFeePerGas:         nil,
			gasPrice:             newInt(v.GasPrice),
			accessList:           ethtype.AccessList{},
			maxFeePerBlobGas:     nil,
			blobHashes:           nil,
			authList:             nil,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethtype.AccessListTx:
		chainIDVal := chainID
		if chainIDVal == nil {
			chainIDVal = v.ChainID
		}
		return &txImpl{
			chainID:              newInt(chainIDVal),
			nonce:                v.Nonce,
			gas:                  v.Gas,
			to:                   v.To,
			value:                newInt(v.Value),
			data:                 v.Data,
			txType:               AccessListTxType,
			maxPriorityFeePerGas: nil,
			maxFeePerGas:         nil,
			gasPrice:             newInt(v.GasPrice),
			accessList:           v.AccessList,
			maxFeePerBlobGas:     nil,
			blobHashes:           nil,
			authList:             nil,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethtype.DynamicFeeTx:
		chainIDVal := chainID
		if chainIDVal == nil {
			chainIDVal = v.ChainID
		}
		return &txImpl{
			chainID:              newInt(chainIDVal),
			nonce:                v.Nonce,
			gas:                  v.Gas,
			to:                   v.To,
			value:                newInt(v.Value),
			data:                 v.Data,
			txType:               DynamicFeeTxType,
			maxPriorityFeePerGas: newInt(v.GasTipCap),
			maxFeePerGas:         newInt(v.GasFeeCap),
			gasPrice:             nil,
			accessList:           v.AccessList,
			maxFeePerBlobGas:     nil,
			blobHashes:           nil,
			authList:             nil,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethtype.BlobTx:
		chainIDVal := chainID
		if chainIDVal == nil {
			chainIDVal = v.ChainID.ToBig()
		}
		toAddr := v.To
		return &txImpl{
			chainID:              newInt(chainIDVal),
			nonce:                v.Nonce,
			gas:                  v.Gas,
			to:                   &toAddr,
			value:                v.Value,
			data:                 v.Data,
			txType:               BlobTxType,
			maxPriorityFeePerGas: v.GasTipCap,
			maxFeePerGas:         v.GasFeeCap,
			gasPrice:             nil,
			accessList:           v.AccessList,
			maxFeePerBlobGas:     v.BlobFeeCap,
			blobHashes:           v.BlobHashes,
			authList:             nil,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethtype.SetCodeTx:
		chainIDVal := chainID
		if chainIDVal == nil {
			chainIDVal = v.ChainID.ToBig()
		}
		toAddr := v.To
		return &txImpl{
			chainID:              newInt(chainIDVal),
			nonce:                v.Nonce,
			gas:                  v.Gas,
			to:                   &toAddr,
			value:                v.Value,
			data:                 v.Data,
			txType:               SetCodeTxType,
			maxPriorityFeePerGas: v.GasTipCap,
			maxFeePerGas:         v.GasFeeCap,
			gasPrice:             nil,
			accessList:           v.AccessList,
			maxFeePerBlobGas:     nil,
			blobHashes:           nil,
			authList:             v.AuthList,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethtype.ETransaction:
		v1 := (*ethTransactionReflect)(unsafe.Pointer(v))
		return newTxImplRaw(v1.Inner, chainID)
	default:
		return nil, fmt.Errorf("unsupported tx type: %T", tx)
	}
}
