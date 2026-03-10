package ethtype

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
