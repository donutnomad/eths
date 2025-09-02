package contractcall

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/donutnomad/blockchain-alg/xsecp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// TxBuilder builds Ethereum transactions
type TxBuilder struct {
	ctx     context.Context
	chainId *big.Int

	nonce    *uint64
	from     common.Address
	to       *common.Address
	value    *big.Int
	data     []byte
	gasPrice *GasPrice
	gasLimit *big.Int

	checkContract bool
	err           error
}

// NewTxBuilder creates a new transaction builder
func NewTxBuilder(ctx context.Context, chainId *big.Int) *TxBuilder {
	return &TxBuilder{
		ctx:     ctx,
		chainId: chainId,
	}
}

func (b *TxBuilder) checkRequiredFields0() error {
	if b.from == (common.Address{}) {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "from is required")
	}
	return nil
}

func (b *TxBuilder) checkRequiredFields1() error {
	if err := b.checkRequiredFields0(); err != nil {
		return err
	}
	if b.gasPrice == nil {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "gas price is required")
	}
	if b.gasPrice.LegacyGas == nil && b.gasPrice.DynamicGas == nil {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "either legacy gas or dynamic gas must be set")
	}
	if b.nonce == nil {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "nonce is required")
	}
	if b.chainId == nil {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "chain id is required")
	}
	return nil
}

func (b *TxBuilder) checkRequiredFields2() error {
	if err := b.checkRequiredFields1(); err != nil {
		return err
	}
	if b.gasLimit == nil {
		return errors.Wrap(TxBuilderMissingRequiredFieldErr, "gas limit is required")
	}
	return nil
}

func (b *TxBuilder) setField(f func(b *TxBuilder)) *TxBuilder {
	if b.err != nil {
		return b
	}
	f(b)
	return b
}

func (b *TxBuilder) Error() error {
	return b.err
}

// SetFrom sets the sender address
func (b *TxBuilder) SetFrom(from common.Address) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.from = from
	})
}

func (b *TxBuilder) SetFromByKey(key ecdsa.PublicKey) *TxBuilder {
	from := common.Address(xsecp256k1.NewPublicKeyFromEcdsa(&key).Address())
	return b.SetFrom(from)
}

// SetTo sets the recipient address (nil for contract deployment)
func (b *TxBuilder) SetTo(to common.Address, checkContract bool) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.to = &to
		b.checkContract = checkContract
	})
}

// SetValue sets the amount of ETH to send
func (b *TxBuilder) SetValue(value *big.Int) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.value = value
	})
}

// SetData sets the transaction data
func (b *TxBuilder) SetData(data []byte) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.data = data
	})
}

// SetGasPrice sets the gas price configuration
func (b *TxBuilder) SetGasPrice(gasPrice *GasPrice) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.gasPrice = gasPrice
	})
}

// SetGasLimit sets the gas limit
func (b *TxBuilder) SetGasLimit(gasLimit uint64) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.gasLimit = new(big.Int).SetUint64(gasLimit)
	})
}

// SetNonce sets the nonce
func (b *TxBuilder) SetNonce(nonce uint64) *TxBuilder {
	return b.setField(func(b *TxBuilder) {
		b.nonce = &nonce
	})
}

// SetNonceBy gets nonce from chain
func (b *TxBuilder) SetNonceBy(transactor IGetNonce) *TxBuilder {
	if b.err != nil {
		return b
	}
	b.err = b.checkRequiredFields0()
	if b.err != nil {
		return b
	}
	nonce, err := transactor.GetNonce(b.ctx, b.from, true)
	if err != nil {
		b.err = errors.WithMessage(err, "failed to get nonce from chain")
		return b
	}
	b.nonce = &nonce
	return b
}

func (b *TxBuilder) SetGasPriceBy(gasPricer IGasPricer) *TxBuilder {
	if b.err != nil {
		return b
	}
	price, err := gasPricer.GetGasPrice(b.ctx, b.chainId)
	if err != nil {
		b.err = err
		return b
	}
	b.SetGasPrice(price)
	return b
}

// SetGasLimitBy estimates gas limit using the provided estimator
func (b *TxBuilder) SetGasLimitBy(estimator IEstimateGas) *TxBuilder {
	if b.err != nil {
		return b
	}
	b.err = b.checkRequiredFields1()
	if b.err != nil {
		return b
	}

	msg := ethereum.CallMsg{
		From:  b.from,
		To:    b.to,
		Data:  b.data,
		Value: b.value,
	}

	if b.gasPrice.LegacyGas != nil {
		msg.GasPrice = b.gasPrice.LegacyGas.GasPrice
	} else if b.gasPrice.DynamicGas != nil {
		msg.GasTipCap = b.gasPrice.DynamicGas.MaxPriorityFeePerGas
		msg.GasFeeCap = b.gasPrice.DynamicGas.MaxFeePerGas
	}

	gasLimit, err := estimator.EstimateGas(b.ctx, b.chainId, msg)
	if err != nil {
		b.err = &EstimateGasError{Err: err}
		return b
	}
	b.gasLimit = gasLimit
	return b
}

func (b *TxBuilder) BalanceCheck(checker IBalanceChecker) *TxBuilder {
	if b.err != nil {
		return b
	}
	if err := b.checkRequiredFields0(); err != nil {
		b.err = err
		return b
	}
	if checker == nil {
		return b
	}
	err := checker.CheckBalance(b.ctx, b.from, b.data, b.to, b.gasPrice)
	if err != nil {
		b.err = err
		return b
	}
	return b
}

func (b *TxBuilder) Check(transactor ICodeAt, gasPriceValidator IGasPriceValidator, opt ...SendTxOption) *TxBuilder {
	var nonStrict = false
	for _, o := range opt {
		if !nonStrict && !o.Strict {
			nonStrict = true
		}
	}

	if b.to != nil {
		if b.checkContract || !nonStrict {
			// Gas estimation cannot succeed without code for method invocations.
			code, err := transactor.PendingCodeAt(b.ctx, *b.to)
			if err != nil {
				b.err = err
				return b
			}
			if b.checkContract && len(code) == 0 {
				b.err = bind.ErrNoCode
				return b
			}
			if !nonStrict && len(code) > 0 && len(b.data) == 0 {
				b.err = ErrContractCallEmptyData
				return b
			}
		}
	}

	if gasPriceValidator != nil {
		if err := gasPriceValidator.ValidateGasPrice(b.ctx, b.chainId, b.gasPrice); err != nil {
			b.err = err
			return b
		}
	}
	return b
}

// Build builds and returns the transaction
func (b *TxBuilder) Build() (*TxWrapper, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := b.checkRequiredFields2(); err != nil {
		b.err = err
		return nil, err
	}

	// Create transaction based on gas price type
	if b.gasPrice.DynamicGas != nil {
		// EIP-1559 transaction
		tx := &ethTypes.DynamicFeeTx{
			To:        b.to,
			Nonce:     *b.nonce,
			GasFeeCap: b.gasPrice.DynamicGas.MaxFeePerGas,
			GasTipCap: b.gasPrice.DynamicGas.MaxPriorityFeePerGas,
			Gas:       b.gasLimit.Uint64(),
			Value:     b.value,
			Data:      b.data,
		}
		return NewTxWrapperDynamic(tx, b.chainId), nil
	} else if b.gasPrice.LegacyGas != nil {
		// Legacy transaction
		tx := &ethTypes.LegacyTx{
			To:       b.to,
			Nonce:    *b.nonce,
			GasPrice: b.gasPrice.LegacyGas.GasPrice,
			Gas:      b.gasLimit.Uint64(),
			Value:    b.value,
			Data:     b.data,
		}
		return NewTxWrapperLegacy(tx, b.chainId), nil
	} else {
		panic(UNREACHABLE)
	}
}

type TxWrapper struct {
	dynamicFeeTx *ethTypes.DynamicFeeTx
	legacyTx     *ethTypes.LegacyTx
	chainID      *big.Int
}

func NewTxWrapperDynamic(tx *ethTypes.DynamicFeeTx, chainID *big.Int) *TxWrapper {
	tx.ChainID = chainID
	return &TxWrapper{
		dynamicFeeTx: tx,
		chainID:      chainID,
	}
}

func NewTxWrapperLegacy(tx *ethTypes.LegacyTx, chainID *big.Int) *TxWrapper {
	return &TxWrapper{
		legacyTx: tx,
		chainID:  chainID,
	}
}

func (w *TxWrapper) SetSignatureValues(v, r, s *big.Int) *TxWrapper {
	if w.dynamicFeeTx != nil {
		w.dynamicFeeTx.V, w.dynamicFeeTx.R, w.dynamicFeeTx.S = v, r, s
	} else if w.legacyTx != nil {
		w.legacyTx.V, w.legacyTx.R, w.legacyTx.S = v, r, s
	}
	return w
}

func (w *TxWrapper) Sign(privateKey ISigner) (*TxWrapper, error) {
	txHashForSign := ethTypes.NewLondonSigner(w.chainID).Hash(w.ToTransaction()).Bytes() // not txHash
	if noOpSigner, ok := privateKey.(*NoOpSigner); ok {
		_, _ = noOpSigner.Sign(txHashForSign)
		return w, nil
	}
	sig, err := privateKey.Sign(txHashForSign)
	if err != nil {
		return nil, err
	}
	w.SetSignatureValues(big.NewInt(int64(sig.V()-27)), sig.R(), sig.S())
	return w, nil
}

func (w *TxWrapper) ToJSON() []byte {
	return lo.Must1(w.ToTransaction().MarshalJSON())
}

func (w *TxWrapper) ToTransaction() *ethTypes.Transaction {
	if w.dynamicFeeTx != nil {
		return ethTypes.NewTx(w.dynamicFeeTx)
	} else if w.legacyTx != nil {
		return ethTypes.NewTx(w.legacyTx)
	}
	panic(UNREACHABLE)
}

var UNREACHABLE = "unreachable"
