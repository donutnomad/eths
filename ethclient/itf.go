package ethclient

//
//type IClient interface {
//	Close()
//	IsOverloaded() bool
//	OverloadRate() float64
//
//	ChainIDReader
//	BlockNumberReader
//	BlockReader
//	LiteBlockReader
//
//	GetBlock(ctx context.Context, blockNrOrHash ethtype.BlockNumberOrHash) (*ethtype.Block, error)
//
//	PeerCount(ctx context.Context) (uint64, error)
//	BlockReceipts(ctx context.Context, blockNrOrHash ethtype.BlockNumberOrHash) ([]*ethtype.TxReceipt, error)
//	HeaderByHash(ctx context.Context, hash ecommon.Hash) (*ethtype.Header, error)
//	HeaderByNumber(ctx context.Context, number *big.Int) (*ethtype.Header, error)
//	TransactionByHash(ctx context.Context, hash ecommon.Hash) (tx *ethtype.ETransaction, isPending bool, err error)
//	TransactionCount(ctx context.Context, blockHash ecommon.Hash) (uint, error)
//	TransactionInBlock(ctx context.Context, blockHash ecommon.Hash, index uint) (*ethtype.ETransaction, error)
//	TransactionReceipt(ctx context.Context, txHash ecommon.Hash) (*ethtype.Receipt, error)
//	SyncProgress(ctx context.Context) (*SyncProgress, error)
//	NetworkID(ctx context.Context) (*big.Int, error)
//	BalanceAt(ctx context.Context, account ecommon.Address, blockNumber *big.Int) (*big.Int, error)
//	BalanceAtHash(ctx context.Context, account ecommon.Address, blockHash ecommon.Hash) (*big.Int, error)
//	StorageAt(ctx context.Context, account ecommon.Address, key ecommon.Hash, blockNumber *big.Int) ([]byte, error)
//	StorageAtHash(ctx context.Context, account ecommon.Address, key ecommon.Hash, blockHash ecommon.Hash) ([]byte, error)
//	CodeAt(ctx context.Context, account ecommon.Address, blockNumber *big.Int) ([]byte, error)
//	CodeAtHash(ctx context.Context, account ecommon.Address, blockHash ecommon.Hash) ([]byte, error)
//	NonceAt(ctx context.Context, account ecommon.Address, blockNumber *big.Int) (uint64, error)
//	NonceAtHash(ctx context.Context, account ecommon.Address, blockHash ecommon.Hash) (uint64, error)
//	FilterLogs(ctx context.Context, q FilterQuery) ([]ethtype.Log, error)
//	SubscribeNewHead(ctx context.Context, ch chan<- *ethtype.Header) (Subscription, error)
//	SubscribeTransactionReceipts(ctx context.Context, q *TransactionReceiptsQuery, ch chan<- []*ethtype.EReceipt) (Subscription, error)
//	SubscribeFilterLogs(ctx context.Context, q FilterQuery, ch chan<- ethtype.ELog) (Subscription, error)
//	PendingBalanceAt(ctx context.Context, account ecommon.Address) (*big.Int, error)
//	PendingStorageAt(ctx context.Context, account ecommon.Address, key ecommon.Hash) ([]byte, error)
//	PendingCodeAt(ctx context.Context, account ecommon.Address) ([]byte, error)
//	PendingNonceAt(ctx context.Context, account ecommon.Address) (uint64, error)
//	PendingTransactionCount(ctx context.Context) (uint, error)
//	CallContract(ctx context.Context, msg CallMsg, blockNumber *big.Int) ([]byte, error)
//	CallContractAtHash(ctx context.Context, msg CallMsg, blockHash ecommon.Hash) ([]byte, error)
//	PendingCallContract(ctx context.Context, msg CallMsg) ([]byte, error)
//	SuggestGasPrice(ctx context.Context) (*big.Int, error)
//	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
//	BlobBaseFee(ctx context.Context) (*big.Int, error)
//	FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*FeeHistory, error)
//	EstimateGas(ctx context.Context, msg CallMsg) (uint64, error)
//	EstimateGasAtBlock(ctx context.Context, msg CallMsg, blockNumber *big.Int) (uint64, error)
//	EstimateGasAtBlockHash(ctx context.Context, msg CallMsg, blockHash ecommon.Hash) (uint64, error)
//	SendTransaction(ctx context.Context, tx *ethtype.ETransaction) error
//	SendTransactionSync(
//		ctx context.Context,
//		tx *ethtype.ETransaction,
//		timeout *time.Duration,
//	) (*ethtype.Receipt, error)
//	SendRawTransactionSync(
//		ctx context.Context,
//		rawTx []byte,
//		timeout *time.Duration,
//	) (*ethtype.Receipt, error)
//	SimulateV1(ctx context.Context, opts SimulateOptions, blockNrOrHash *rpc.BlockNumberOrHash) ([]SimulateBlockResult, error)
//}
