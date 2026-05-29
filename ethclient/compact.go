package ethclient

import (
	"context"
	"math/big"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/ethtype"
	"github.com/ethereum/go-ethereum"
)

var NotFound = ethereum.NotFound

type CallMsg = ethereum.CallMsg
type Subscription = ethereum.Subscription
type TransactionReceiptsQuery = ethereum.TransactionReceiptsQuery
type SyncProgress = ethereum.SyncProgress
type FeeHistory = ethereum.FeeHistory
type BlockOverrides = ethereum.BlockOverrides
type OverrideAccount = ethereum.OverrideAccount

// GasEstimator wraps EstimateGas, which tries to estimate the gas needed to execute a
// specific transaction based on the pending state. There is no guarantee that this is the
// true gas limit requirement as other transactions may be added or removed by miners, but
// it should provide a basis for setting a reasonable default.
type GasEstimator interface {
	EstimateGas(ctx context.Context, call CallMsg) (uint64, error)
}

// TransactionSender wraps transaction sending. The SendTransaction method injects a
// signed transaction into the pending transaction pool for execution. If the transaction
// was a contract creation, the TransactionReceipt method can be used to retrieve the
// contract address after the transaction has been mined.
//
// The transaction must be signed and have a valid nonce to be included. Consumers of the
// API can use package accounts to maintain local private keys and need can retrieve the
// next available nonce using PendingNonceAt.
type TransactionSender interface {
	SendTransaction(ctx context.Context, tx *ethtype.ETransaction) error
}

// BlockNumberReader provides access to the current block number.
type BlockNumberReader interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

// TransactionReader provides access to past transactions and their receipts.
// Implementations may impose arbitrary restrictions on the transactions and receipts that
// can be retrieved. Historic transactions may not be available.
//
// Avoid relying on this interface if possible. Contract logs (through the LogFilterer
// interface) are more reliable and usually safer in the presence of chain
// reorganisations.
//
// The returned error is NotFound if the requested item does not exist.
type TransactionReader interface {
	// TransactionByHash checks the pool of pending transactions in addition to the
	// blockchain. The isPending return value indicates whether the transaction has been
	// mined yet. Note that the transaction may not be part of the canonical chain even if
	// it's not pending.
	TransactionByHash(ctx context.Context, txHash common.Hash) (tx *ethtype.ETransaction, isPending bool, err error)
	// TransactionReceipt returns the receipt of a mined transaction. Note that the
	// transaction may not be included in the current canonical chain even if a receipt
	// exists.
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethtype.EReceipt, error)
	// SubscribeTransactionReceipts subscribes to notifications about transaction receipts.
	// The receipts are delivered in batches when transactions are included in blocks.
	// If q is nil or has empty TransactionHashes, all receipts from new blocks will be delivered.
	// Otherwise, only receipts for the specified transaction hashes will be delivered.
	SubscribeTransactionReceipts(ctx context.Context, q *TransactionReceiptsQuery, ch chan<- []*ethtype.EReceipt) (Subscription, error)
}

// GasPricer wraps the gas price oracle, which monitors the blockchain to determine the
// optimal gas price given current fee market conditions.
type GasPricer interface {
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

// GasPricer1559 provides access to the EIP-1559 gas price oracle.
type GasPricer1559 interface {
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

// ChainIDReader provides access to the chain ID.
type ChainIDReader interface {
	ChainID(ctx context.Context) (*big.Int, error)
}

// LogFilterer provides access to contract log events using a one-off query or continuous
// event subscription.
//
// Logs received through a streaming query subscription may have Removed set to true,
// indicating that the log was reverted due to a chain reorganisation.
type LogFilterer interface {
	FilterLogs(ctx context.Context, q FilterQuery) ([]ethtype.ELog, error)
	SubscribeFilterLogs(ctx context.Context, q FilterQuery, ch chan<- ethtype.ELog) (Subscription, error)
}

// A PendingStateReader provides access to the pending state, which is the result of all
// known executable transactions which have not yet been included in the blockchain. It is
// commonly used to display the result of ’unconfirmed’ actions (e.g. wallet value
// transfers) initiated by the user. The PendingNonceAt operation is a good way to
// retrieve the next available transaction nonce for a specific account.
type PendingStateReader interface {
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
}

// PendingContractCaller can be used to perform calls against the pending state.
type PendingContractCaller interface {
	PendingCallContract(ctx context.Context, call CallMsg) ([]byte, error)
}

// A ContractCaller provides contract calls, essentially transactions that are executed by
// the EVM but not mined into the blockchain. ContractCall is a low-level method to
// execute such calls. For applications which are structured around specific contracts,
// the abigen tool provides a nicer, properly typed way to perform calls.
type ContractCaller interface {
	CallContract(ctx context.Context, call CallMsg, blockNumber *big.Int) ([]byte, error)
}

// ChainStateReader wraps access to the state trie of the canonical blockchain. Note that
// implementations of the interface may be unable to return state values for old blocks.
// In many cases, using CallContract can be preferable to reading raw contract storage.
type ChainStateReader interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}

// ChainReader provides access to the blockchain. The methods in this interface access raw
// data from either the canonical chain (when requesting by block number) or any
// blockchain fork that was previously downloaded and processed by the node. The block
// number argument can be nil to select the latest canonical block. Reading block headers
// should be preferred over full blocks whenever possible.
//
// The returned error is NotFound if the requested item does not exist.
type ChainReader interface {
	BlockByHash(ctx context.Context, hash common.Hash) (*ethtype.EBlock, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethtype.EBlock, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*ethtype.EHeader, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethtype.EHeader, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*ethtype.ETransaction, error)

	// This method subscribes to notifications about changes of the head block of
	// the canonical chain.
	SubscribeNewHead(ctx context.Context, ch chan<- *ethtype.EHeader) (Subscription, error)
}
