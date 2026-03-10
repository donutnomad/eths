package ethtype

// TxDetail combines transaction fields with receipt fields,
// embedding Tx and adding non-overlapping receipt fields.
type TxDetail struct {
	Tx
	Receipt
}
