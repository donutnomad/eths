package contractcall

import (
	"fmt"
	"math/big"
	"unsafe"

	"slices"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/samber/lo"
)

var UNREACHABLE = "unreachable"

// Deprecated: Use Tx
type TxWrapper = Tx

type Tx = txImpl

// NewTxWrapperDynamic
// Deprecated: Use NewTxImpl
func NewTxWrapperDynamic(tx *ethTypes.DynamicFeeTx, chainID *big.Int) *Tx {
	return NewTxImpl(tx, chainID).(*Tx)
}

// NewTxWrapperLegacy
// Deprecated: Use NewTxImpl
func NewTxWrapperLegacy(tx *ethTypes.LegacyTx, chainID *big.Int) *Tx {
	return NewTxImpl(tx, chainID).(*Tx)
}

type AccessListSetter interface {
	AccessList() ethTypes.AccessList
	SetAccessList(accessList ethTypes.AccessList) bool
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
	BlobFeeCap() *big.Int
	BlobHashes() []common.Hash
	SetBlobFeeCap(blobFeeCap *big.Int) bool
	SetBlobHashes(blobHashes []common.Hash) bool
}

// IEIP7702 Set code
type IEIP7702 interface {
	IEIP1559
	AuthList() []ethTypes.SetCodeAuthorization
	SetAuthList(authList []ethTypes.SetCodeAuthorization) bool
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
	ToTransaction() *ethTypes.Transaction
}

var _ ITx = &txImpl{}

type txImpl struct {
	chainID *uint256.Int // safe, not nil
	nonce   uint64
	gas     uint64
	to      *common.Address
	value   *uint256.Int // safe, not nil
	data    []byte

	maxPriorityFeePerGas *uint256.Int                    // notInclude(eip-155(legacy),eip-2930(access))
	maxFeePerGas         *uint256.Int                    // notInclude(eip-155(legacy),eip-2930(access))
	gasPrice             *uint256.Int                    // eip-155(legacy), eip-2930(access)
	accessList           ethTypes.AccessList             // notInclude(eip-155(legacy))
	maxFeePerBlobGas     *uint256.Int                    // eip-4844(blob)
	blobHashes           []common.Hash                   // eip-4844(blob)
	authList             []ethTypes.SetCodeAuthorization // eip-7702(auth)

	v      [32]byte // safe, default 0
	r      [32]byte // safe, default 0
	s      [32]byte // safe, default 0
	txType TxType
}

func NewTxImplWith(txType TxType, chainID *big.Int) ITx {
	return &txImpl{
		chainID: newInt(chainID),
		txType:  txType,
		value:   newIntBy(0),
	}
}

func NewTxImpl[T *ethTypes.LegacyTx | *ethTypes.AccessListTx | *ethTypes.DynamicFeeTx | *ethTypes.BlobTx | *ethTypes.SetCodeTx | *ethTypes.Transaction](tx T, chainID *big.Int) ITx {
	ret, err := newTxImpl(any(tx), chainID)
	if err != nil {
		panic(err)
	}
	return ret
}

func (t *txImpl) SigHash() [32]byte {
	return ethTypes.NewPragueSigner(t.ChainID()).Hash(t.ToTransaction())
}

func (t *txImpl) Sign(privateKey ISigner) error {
	var sigHash = t.SigHash()
	if noOpSigner, ok := privateKey.(*NoOpSigner); ok {
		sig, err := noOpSigner.Sign(sigHash[:])
		if err != nil {
			return err
		}
		// We set the signature to v, 0x01, 0x01 to retain the chainID of unsigned LegacyTx transactions.
		if sig == nil && t.txType == LegacyTxType {
			v := computeVForEIP155(27, t.ChainID(), true)
			t.SetSignature(v, common.Big1, common.Big1)
		}
		return nil
	}
	sig, err := privateKey.Sign(sigHash[:])
	if err != nil {
		return err
	}
	v := computeVForEIP155(sig.V(), t.ChainID(), t.txType == LegacyTxType)
	t.SetSignature(v, sig.R(), sig.S())
	return nil
}

func computeVForEIP155(sigV byte, chainID *big.Int, isLegacyTx bool) *big.Int {
	// v: 0/1
	v := big.NewInt(int64(sigV) - 27)
	if isLegacyTx { // EIP155-Fork
		// (0/1 + 35)(35/36) + chainID * 2
		mulChainID := new(big.Int).Mul(chainID, big.NewInt(2))
		v.Add(v, big.NewInt(35))
		v.Add(v, mulChainID)
	}
	return v
}

func (t *txImpl) ToJSON() []byte {
	return lo.Must1(t.ToTransaction().MarshalJSON())
}

func (t *txImpl) ToTransaction() *ethTypes.Transaction {
	return ethTypes.NewTx(t.Export())
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
	if t.txType == LegacyTxType || t.txType == AccessListTxType {
		return t.gasPrice.ToBig()
	} else {
		return nil
	}
}

// SetGasPrice eip-155(legacy), eip-2930(access)
func (t *txImpl) SetGasPrice(gasPrice *big.Int) bool {
	if t.txType == LegacyTxType || t.txType == AccessListTxType {
		t.gasPrice = newInt(gasPrice)
		return true
	} else {
		return false
	}
}

// AccessList notInclude(eip-155(legacy))
func (t *txImpl) AccessList() ethTypes.AccessList {
	if t.txType != LegacyTxType {
		return t.accessList
	}
	return nil
}

// SetAccessList notInclude(eip-155(legacy))
func (t *txImpl) SetAccessList(accessList ethTypes.AccessList) bool {
	if t.txType != LegacyTxType {
		t.accessList = accessList
		return true
	}
	return false
}

// MaxFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) MaxFeePerGas() *big.Int {
	if t.txType != LegacyTxType && t.txType != AccessListTxType {
		return t.maxFeePerGas.ToBig()
	} else {
		return nil
	}
}

// MaxPriorityFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) MaxPriorityFeePerGas() *big.Int {
	if t.txType != LegacyTxType && t.txType != AccessListTxType {
		return t.maxPriorityFeePerGas.ToBig()
	} else {
		return nil
	}
}

// SetMaxFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) SetMaxFeePerGas(maxFeePerGas *big.Int) bool {
	if t.txType != LegacyTxType && t.txType != AccessListTxType {
		t.maxFeePerGas = newInt(maxFeePerGas)
		return true
	} else {
		return false
	}
}

// SetMaxPriorityFeePerGas notInclude(eip-155(legacy),eip-2930(access))
func (t *txImpl) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) bool {
	if t.txType != LegacyTxType && t.txType != AccessListTxType {
		t.maxPriorityFeePerGas = newInt(maxPriorityFeePerGas)
		return true
	} else {
		return false
	}
}

// BlobFeeCap eip-4844(blob)
func (t *txImpl) BlobFeeCap() *big.Int {
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

// SetBlobFeeCap eip-4844(blob)
func (t *txImpl) SetBlobFeeCap(blobFeeCap *big.Int) bool {
	if t.txType == BlobTxType {
		t.maxFeePerBlobGas = newInt(blobFeeCap)
		return true
	} else {
		return false
	}
}

// SetBlobHashes eip-4844(blob)
func (t *txImpl) SetBlobHashes(blobHashes []common.Hash) bool {
	if t.txType == BlobTxType {
		t.blobHashes = blobHashes
		return true
	} else {
		return false
	}
}

// AuthList eip-7702(auth)
func (t *txImpl) AuthList() []ethTypes.SetCodeAuthorization {
	if t.txType == SetCodeTxType {
		return t.authList
	} else {
		return nil
	}
}

// SetAuthList eip-7702(auth)
func (t *txImpl) SetAuthList(authList []ethTypes.SetCodeAuthorization) bool {
	if t.txType == SetCodeTxType {
		t.authList = slices.Clone(authList)
		return true
	} else {
		return false
	}
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

func (t *txImpl) UnmarshalJSON(input []byte) error {
	var tx = new(ethTypes.Transaction)
	err := tx.UnmarshalJSON(input)
	if err != nil {
		return err
	}
	ret, err := newTxImpl(tx, tx.ChainId())
	if err != nil {
		return err
	}
	*t = *ret
	return nil
}

func (t *txImpl) Export() ethTypes.TxData {
	return txImplToTx(t)
}

func txImplToTx(t *txImpl) ethTypes.TxData {
	switch t.txType {
	case LegacyTxType:
		if t.gasPrice == nil {
			t.gasPrice = newIntBy(0)
		}
		return &ethTypes.LegacyTx{
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
		return &ethTypes.AccessListTx{
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
		return &ethTypes.DynamicFeeTx{
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
		return &ethTypes.BlobTx{
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
			Sidecar:    nil,
			V:          t.getV(),
			R:          t.getR(),
			S:          t.getS(),
		}
	case SetCodeTxType:
		return &ethTypes.SetCodeTx{
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

func newTxImpl(tx any, chainID *big.Int) (*txImpl, error) {
	if v, ok := tx.(*ethTypes.Transaction); ok {
		chainID = v.ChainId()
	}
	impl, err := newTxImplRaw(tx, chainID)
	if err != nil {
		return nil, err
	}
	if v, ok := tx.(*ethTypes.Transaction); ok {
		// v > 0 && r == 1 && s == 1, This is a magic value we have set, which retains the chainID of unsigned LegacyTx transactions.
		_v, _r, _s := v.RawSignatureValues()
		if v.Type() == uint8(LegacyTxType) && _v.Cmp(common.Big0) > 0 && _r.Cmp(common.Big1) == 0 && _s.Cmp(common.Big1) == 0 {
			impl.v = [32]byte{}
			impl.r = [32]byte{}
			impl.s = [32]byte{}
		}
	}
	return impl, nil
}

func newTxImplRaw(tx any, chainID *big.Int) (*txImpl, error) {
	switch v := tx.(type) {
	case *ethTypes.LegacyTx:
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
			accessList:           ethTypes.AccessList{},
			maxFeePerBlobGas:     nil,
			blobHashes:           nil,
			authList:             nil,
			v:                    bigToBytes32(v.V),
			r:                    bigToBytes32(v.R),
			s:                    bigToBytes32(v.S),
		}, nil
	case *ethTypes.AccessListTx:
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
	case *ethTypes.DynamicFeeTx:
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
	case *ethTypes.BlobTx:
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
	case *ethTypes.SetCodeTx:
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
	case *ethTypes.Transaction:
		v1 := (*ethTransactionReflect)(unsafe.Pointer(v))
		return newTxImplRaw(v1.Inner, chainID)
	default:
		return nil, fmt.Errorf("unsupported tx type: %T", tx)
	}
}
