package contractcall

type TxType int

const (
	LegacyTxType     TxType = 0x00
	AccessListTxType TxType = 0x01
	DynamicFeeTxType TxType = 0x02
	BlobTxType       TxType = 0x03
	SetCodeTxType    TxType = 0x04
)
