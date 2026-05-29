package ethtype

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

//go:generate bash ../internal/gencodec/run.sh -type Header -field-override headerMarshaling -out header_generated.go
//go:generate bash ../internal/gencodec/run.sh -type Block -field-override headerMarshaling -out block_generated.go
//go:generate bash ../internal/gencodec/run.sh -type LiteBlock -field-override headerMarshaling -out blocklite_generated.go
//go:generate bash ../internal/gencodec/run.sh -type TxReceipt -field-override receiptMarshaling -out receipt_generated.go
//go:generate bash ../internal/gencodec/run.sh -type Log -field-override logMarshaling -out log_generated.go
//go:generate bash ../internal/gencodec/run.sh -type Tx -field-override txMarshaling -out tx_generated.go
//go:generate bash ../internal/gencodec/run.sh -type Withdrawal -field-override withdrawalMarshaling -out withdrawal_generated.go
//go:generate bash ../internal/gencodec/run.sh -type TxDetail -field-override txMarshaling,receiptMarshaling -out txdetail_generated.go
//go:generate bash ../internal/gencodec/run.sh -type AccessTuple -out accesslist_generated.go
//go:generate bash ../internal/gencodec/run.sh -type SetCodeAuthorization -field-override authorizationMarshaling -out authorization_generated.go

type EHeader = ethTypes.Header
type EBlock = ethTypes.Block
type EReceipt = ethTypes.Receipt
type AccessList = ethTypes.AccessList
type ETransaction = ethTypes.Transaction
type ELog = ethTypes.Log
type DynamicFeeTx = ethTypes.DynamicFeeTx
type LegacyTx = ethTypes.LegacyTx
type AccessListTx = ethTypes.AccessListTx
type BlobTx = ethTypes.BlobTx
type SetCodeTx = ethTypes.SetCodeTx
type TxData = ethTypes.TxData
type BlobTxSidecar = ethTypes.BlobTxSidecar
type SetCodeAuthorization = ethTypes.SetCodeAuthorization

const ReceiptStatusSuccessful = ethTypes.ReceiptStatusSuccessful
const ReceiptStatusFailed = ethTypes.ReceiptStatusFailed

func NewTx(inner ethTypes.TxData) *ETransaction {
	return ethTypes.NewTx(inner)
}

type TypedDataDomain = apitypes.TypedDataDomain
type TypedDataType = apitypes.Type
