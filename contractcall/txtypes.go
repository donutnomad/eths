package contractcall

import "github.com/donutnomad/eths/ethtype"

type TxType = ethtype.TxType

const (
	LegacyTxType     = ethtype.LegacyTxType
	AccessListTxType = ethtype.AccessListTxType
	DynamicFeeTxType = ethtype.DynamicFeeTxType
	BlobTxType       = ethtype.BlobTxType
	SetCodeTxType    = ethtype.SetCodeTxType
)
