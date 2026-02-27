package ethclient

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

const sepoliaRPC = "https://ethereum-sepolia-rpc.publicnode.com"

var (
	// knownSepoliaTx is a confirmed Sepolia ERC-20 Transfer tx in block 0x9ddc3a.
	knownSepoliaTx   = common.HexToHash("0xf542851546efc4a89a06c8b394a607afa74e8500688d09d9d50be7698adefd5f")
	knownBlockHash   = common.HexToHash("0x9962d2a3f02043ff6dce9b757eb4782cd10dcde92d66bfbe9377d37ea4bf6aef")
	knownBlockNumber = big.NewInt(0x9ddc3a)
	knownTxSender    = common.HexToAddress("0x690c39adabdea83322bf8e90626cd40eeb456a95")
	knownTxTo        = common.HexToAddress("0x3429519ee7cdbb13b49161f1eb6e1b026939113a") // ERC-20 contract
	sepoliaChainID   = big.NewInt(11155111)
	transferTopic    = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	fakeHash         = common.HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
)

func dialSepolia(t *testing.T) *Client {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := rpc.DialContext(ctx, sepoliaRPC)
	if err != nil {
		t.Skipf("cannot connect to Sepolia RPC: %v", err)
	}
	return NewClient(c)
}

func sepoliaCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 15*time.Second)
}

// --- ChainID / NetworkID ---

func TestChainID_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	id, err := ec.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if id.Cmp(sepoliaChainID) != 0 {
		t.Fatalf("ChainID = %v, want %v", id, sepoliaChainID)
	}
}

func TestNetworkID_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	id, err := ec.NetworkID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if id.Cmp(sepoliaChainID) != 0 {
		t.Fatalf("NetworkID = %v, want %v", id, sepoliaChainID)
	}
}

// --- BlockNumber ---

func TestBlockNumber_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	num, err := ec.BlockNumber(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if num < knownBlockNumber.Uint64() {
		t.Fatalf("BlockNumber = %d, expected >= %d", num, knownBlockNumber.Uint64())
	}
}

// --- SyncProgress ---

func TestSyncProgress_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	// Public node should be synced, so progress is nil.
	progress, err := ec.SyncProgress(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if progress != nil {
		t.Fatalf("expected nil SyncProgress on synced node, got %+v", progress)
	}
}

// --- HeaderByNumber / HeaderByHash ---

func TestHeaderByNumber_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	head, err := ec.HeaderByNumber(ctx, knownBlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if head.Number.Cmp(knownBlockNumber) != 0 {
		t.Fatalf("HeaderByNumber number = %v, want %v", head.Number, knownBlockNumber)
	}
	if head.Hash() != knownBlockHash {
		t.Fatalf("HeaderByNumber hash = %s, want %s", head.Hash(), knownBlockHash)
	}
}

func TestHeaderByHash_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	head, err := ec.HeaderByHash(ctx, knownBlockHash)
	if err != nil {
		t.Fatal(err)
	}
	if head.Number.Cmp(knownBlockNumber) != 0 {
		t.Fatalf("HeaderByHash number = %v, want %v", head.Number, knownBlockNumber)
	}
}

func TestHeaderByNumber_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, err := ec.HeaderByNumber(ctx, big.NewInt(0x7fffffffffffffff))
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

// --- LiteBlockByNumber / LiteBlockByHash ---

func TestLiteBlockByNumber_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	block, err := ec.LiteBlockByNumber(ctx, knownBlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if block.Hash() != knownBlockHash {
		t.Fatalf("hash = %s, want %s", block.Hash(), knownBlockHash)
	}
	if block.Number.Cmp(knownBlockNumber) != 0 {
		t.Fatalf("number = %v, want %v", block.Number, knownBlockNumber)
	}
	if len(block.Transactions) != 151 {
		t.Fatalf("tx count = %d, want 151", len(block.Transactions))
	}
	if block.Transactions[0] != knownSepoliaTx {
		t.Fatalf("first tx = %s, want %s", block.Transactions[0], knownSepoliaTx)
	}
}

func TestLiteBlockByHash_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	block, err := ec.LiteBlockByHash(ctx, knownBlockHash)
	if err != nil {
		t.Fatal(err)
	}
	if block.Number.Cmp(knownBlockNumber) != 0 {
		t.Fatalf("number = %v, want %v", block.Number, knownBlockNumber)
	}
	if len(block.Transactions) != 151 {
		t.Fatalf("tx count = %d, want 151", len(block.Transactions))
	}
}

func TestLiteBlockByNumber_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, err := ec.LiteBlockByNumber(ctx, big.NewInt(0x7fffffffffffffff))
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

// --- BlockByNumber / BlockByHash ---

func TestBlockByNumber_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	block, err := ec.BlockByNumber(ctx, knownBlockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if block.Hash() != knownBlockHash {
		t.Fatalf("hash = %s, want %s", block.Hash(), knownBlockHash)
	}
	if block.NumberU64() != knownBlockNumber.Uint64() {
		t.Fatalf("number = %d, want %d", block.NumberU64(), knownBlockNumber.Uint64())
	}
	if len(block.Transactions()) != 151 {
		t.Fatalf("tx count = %d, want 151", len(block.Transactions()))
	}
}

func TestBlockByHash_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	block, err := ec.BlockByHash(ctx, knownBlockHash)
	if err != nil {
		t.Fatal(err)
	}
	if block.NumberU64() != knownBlockNumber.Uint64() {
		t.Fatalf("number = %d, want %d", block.NumberU64(), knownBlockNumber.Uint64())
	}
}

func TestBlockByNumber_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, err := ec.BlockByNumber(ctx, big.NewInt(0x7fffffffffffffff))
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

// --- TransactionByHash ---

func TestTransactionByHash_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	tx, isPending, err := ec.TransactionByHash(ctx, knownSepoliaTx)
	if err != nil {
		t.Fatal(err)
	}
	if isPending {
		t.Fatal("expected confirmed tx, got pending")
	}
	if tx.Hash() != knownSepoliaTx {
		t.Fatalf("hash = %s, want %s", tx.Hash(), knownSepoliaTx)
	}
	if *tx.To() != knownTxTo {
		t.Fatalf("to = %s, want %s", tx.To(), knownTxTo)
	}
}

func TestTransactionByHash_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, _, err := ec.TransactionByHash(ctx, fakeHash)
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

// --- TransactionCount / TransactionInBlock ---

func TestTransactionCount_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	count, err := ec.TransactionCount(ctx, knownBlockHash)
	if err != nil {
		t.Fatal(err)
	}
	if count != 151 {
		t.Fatalf("TransactionCount = %d, want 151", count)
	}
}

func TestTransactionInBlock_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	tx, err := ec.TransactionInBlock(ctx, knownBlockHash, 0)
	if err != nil {
		t.Fatal(err)
	}
	if tx.Hash() != knownSepoliaTx {
		t.Fatalf("hash = %s, want %s", tx.Hash(), knownSepoliaTx)
	}
}

func TestTransactionInBlock_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, err := ec.TransactionInBlock(ctx, knownBlockHash, 9999)
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

// --- TransactionReceipt ---

func TestTransactionReceipt_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	receipt, err := ec.TransactionReceipt(ctx, knownSepoliaTx)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.TxHash != knownSepoliaTx {
		t.Fatalf("TxHash = %s, want %s", receipt.TxHash, knownSepoliaTx)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("Status = %d, want 1", receipt.Status)
	}
	if receipt.BlockNumber == nil || receipt.BlockNumber.Cmp(knownBlockNumber) != 0 {
		t.Fatalf("BlockNumber = %v, want %v", receipt.BlockNumber, knownBlockNumber)
	}
	if len(receipt.Logs) == 0 {
		t.Fatal("expected at least one log")
	}
	if receipt.Logs[0].Topics[0] != transferTopic {
		t.Fatalf("log topic = %s, want %s", receipt.Logs[0].Topics[0], transferTopic)
	}
}

func TestTransactionReceipt_NotFound_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	_, err := ec.TransactionReceipt(ctx, fakeHash)
	if err != ethereum.NotFound {
		t.Fatalf("expected ethereum.NotFound, got %v", err)
	}
}

func TestTransactionReceiptAs_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	type LiteReceipt struct {
		TxHash           common.Hash    `json:"transactionHash"`
		Status           hexutil.Uint64 `json:"status"`
		BlockNumber      hexutil.Uint64 `json:"blockNumber"`
		GasUsed          hexutil.Uint64 `json:"gasUsed"`
		TransactionIndex hexutil.Uint64 `json:"transactionIndex"`
	}

	receipt, err := TransactionReceiptAs[*LiteReceipt](ctx, ec, knownSepoliaTx)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.TxHash != knownSepoliaTx {
		t.Fatalf("TxHash = %s, want %s", receipt.TxHash, knownSepoliaTx)
	}
	if receipt.BlockNumber == 0 {
		t.Fatal("expected non-zero BlockNumber")
	}
	if receipt.GasUsed == 0 {
		t.Fatal("expected non-zero GasUsed")
	}
}

// --- BlockReceipts ---

//type TxEntry struct {
//	BlockHash         common.Hash      `json:"blockHash"`
//	BlockNumber       *BigInt          `json:"blockNumber"`
//	ContractAddress   *common.Address  `json:"contractAddress"`
//	CumulativeGasUsed *hexutil.Big     `json:"cumulativeGasUsed"`
//	EffectiveGasPrice *BigInt          `json:"effectiveGasPrice"`
//	From              common.Address   `json:"from"`
//	GasUsed           etherscan.Uint64 `json:"gasUsed"`
//	LogsBloom         types.Bloom      `json:"logsBloom"`
//	Status            etherscan.Uint64 `json:"status"`
//	To                *common.Address  `json:"to"`
//	TransactionHash   common.Hash      `json:"transactionHash"`
//	TransactionIndex  etherscan.Uint64 `json:"transactionIndex"`
//	Type              etherscan.Uint64 `json:"type"`
//
//	Timeboosted   bool              `json:"timeboosted"`   // arb特有，sepolia没有
//	GasUsedForL1  *etherscan.Uint64 `json:"gasUsedForL1"`  // arb特有，sepolia没有
//	L1BlockNumber *BigInt           `json:"l1BlockNumber"` // arb特有，sepolia没有
//}
//
//type BigInt = hexutil.Big
//
//func TestBlockReceipts_Sepolia2(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping RPC integration test")
//	}
//	ec := dialSepolia(t)
//	defer ec.Close()
//	ctx, cancel := sepoliaCtx()
//	defer cancel()
//
//	tag := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(knownBlockNumber.Int64()))
//	receipts, err := BlockReceiptsAs[*TxEntry](ctx, ec, tag)
//	if err != nil {
//		t.Fatal(err)
//	}
//	spew.Dump(receipts)
//	if len(receipts) != 151 {
//		t.Fatalf("receipt count = %d, want 151", len(receipts))
//	}
//	if receipts[0].TransactionHash != knownSepoliaTx {
//		t.Fatalf("first receipt TxHash = %s, want %s", receipts[0].TransactionHash, knownSepoliaTx)
//	}
//}

func TestBlockReceipts_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	tag := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(knownBlockNumber.Int64()))
	receipts, err := ec.BlockReceipts(ctx, tag)
	if err != nil {
		t.Fatal(err)
	}
	if len(receipts) != 151 {
		t.Fatalf("receipt count = %d, want 151", len(receipts))
	}
	if receipts[0].TxHash != knownSepoliaTx {
		t.Fatalf("first receipt TxHash = %s, want %s", receipts[0].TxHash, knownSepoliaTx)
	}
}

// --- BalanceAt ---

func TestBalanceAt_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	// Contract address should have 0 ETH balance.
	bal, err := ec.BalanceAt(ctx, knownTxTo, nil)
	if err != nil {
		t.Fatal(err)
	}
	if bal == nil {
		t.Fatal("expected non-nil balance")
	}
	// We don't assert exact value since it could change, just that call succeeds.
}

// --- NonceAt ---

func TestNonceAt_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	nonce, err := ec.NonceAt(ctx, knownTxSender, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Sender has sent many transactions, nonce should be > 0.
	if nonce == 0 {
		t.Fatal("expected non-zero nonce for active account")
	}
}

// --- CodeAt ---

func TestCodeAt_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	// knownTxTo is an ERC-20 contract, should have code.
	code, err := ec.CodeAt(ctx, knownTxTo, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) == 0 {
		t.Fatal("expected non-empty code for contract address")
	}

	// EOA should have no code.
	code, err = ec.CodeAt(ctx, knownTxSender, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 0 {
		t.Fatalf("expected empty code for EOA, got %d bytes", len(code))
	}
}

// --- StorageAt ---

func TestStorageAt_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	// Slot 0 of the ERC-20 contract (typically totalSupply or name).
	storage, err := ec.StorageAt(ctx, knownTxTo, common.Hash{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(storage) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(storage))
	}
}

// --- FilterLogs ---

func TestFilterLogs_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	logs, err := ec.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: &knownBlockHash,
		Addresses: []common.Address{knownTxTo},
		Topics:    [][]common.Hash{{transferTopic}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Fatalf("log count = %d, want 1", len(logs))
	}
	if logs[0].TxHash != knownSepoliaTx {
		t.Fatalf("log TxHash = %s, want %s", logs[0].TxHash, knownSepoliaTx)
	}
}

// --- SuggestGasPrice ---

func TestSuggestGasPrice_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	price, err := ec.SuggestGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if price == nil || price.Sign() <= 0 {
		t.Fatalf("expected positive gas price, got %v", price)
	}
}

// --- SuggestGasTipCap ---

func TestSuggestGasTipCap_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	tip, err := ec.SuggestGasTipCap(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if tip == nil {
		t.Fatal("expected non-nil tip cap")
	}
}

// --- FeeHistory ---

func TestFeeHistory_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	history, err := ec.FeeHistory(ctx, 1, nil, []float64{25, 75})
	if err != nil {
		t.Fatal(err)
	}
	if history.OldestBlock == nil {
		t.Fatal("expected non-nil OldestBlock")
	}
	if len(history.BaseFee) == 0 {
		t.Fatal("expected non-empty BaseFee")
	}
	if len(history.GasUsedRatio) != 1 {
		t.Fatalf("GasUsedRatio length = %d, want 1", len(history.GasUsedRatio))
	}
	if len(history.Reward) != 1 || len(history.Reward[0]) != 2 {
		t.Fatalf("Reward shape unexpected: %v", history.Reward)
	}
}

// --- EstimateGas ---

func TestEstimateGas_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	gas, err := ec.EstimateGas(ctx, ethereum.CallMsg{
		From:  knownTxSender,
		To:    &to,
		Value: big.NewInt(0),
	})
	if err != nil {
		t.Fatal(err)
	}
	if gas < 21000 {
		t.Fatalf("EstimateGas = %d, expected >= 21000", gas)
	}
}

// --- CallContract ---

func TestCallContract_Sepolia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RPC integration test")
	}
	ec := dialSepolia(t)
	defer ec.Close()
	ctx, cancel := sepoliaCtx()
	defer cancel()

	// Call ERC-20 totalSupply() = 0x18160ddd
	data := common.FromHex("0x18160ddd")
	result, err := ec.CallContract(ctx, ethereum.CallMsg{
		To:   &knownTxTo,
		Data: data,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 32 {
		t.Fatalf("expected 32 bytes return, got %d", len(result))
	}
}
