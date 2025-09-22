package etherscan

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

func (e *EtherscanClient) GetLogs(
	ctx context.Context,
	chainID uint64,
	address *common.Address, // 可选(address和opts二选一)
	fromBlock *uint64, // nil则从最开始的区块开始
	toBlock *uint64, // nil则为最新的区块
	page int, // min 1
	offset int, // max 1000
	opts GetLogsOptions,
) (*Response[[]LogEntry], error) {
	var m = map[string]string{
		"page":   strconv.Itoa(max(page, 1)),
		"offset": strconv.Itoa(min(offset, 1000)),
	}
	if fromBlock != nil {
		m["fromBlock"] = strconv.FormatUint(*fromBlock, 10)
	}
	if toBlock != nil {
		m["toBlock"] = strconv.FormatUint(*toBlock, 10)
	}
	if address != nil {
		m["address"] = address.Hex()
	}
	for k, v := range opts.toMap() {
		m[k] = v
	}
	return getResult[[]LogEntry](ctx, e, m, "logs", "getLogs", chainID)
}

type GetLogsOptions struct {
	Topic0     mo.Option[common.Hash] // 可选：主题0
	Topic1     mo.Option[common.Hash] // 可选：主题1
	Topic2     mo.Option[common.Hash] // 可选：主题2
	Topic3     mo.Option[common.Hash] // 可选：主题3
	Topic01Opr mo.Option[bool]        // 可选：主题0和1之间的操作符 (and|or), true: and, false: or
	Topic12Opr mo.Option[bool]        // 可选：主题1和2之间的操作符 (and|or), true: and, false: or
	Topic23Opr mo.Option[bool]        // 可选：主题2和3之间的操作符 (and|or), true: and, false: or
	Topic02Opr mo.Option[bool]        // 可选：主题0和2之间的操作符 (and|or), true: and, false: or
	Topic03Opr mo.Option[bool]        // 可选：主题0和3之间的操作符 (and|or), true: and, false: or
	Topic13Opr mo.Option[bool]        // 可选：主题1和3之间的操作符 (and|or), true: and, false: or
}

func (opts *GetLogsOptions) toMap() map[string]string {
	var params = make(map[string]string)
	var setTopic = func(key string, topic mo.Option[common.Hash]) {
		if topic.IsPresent() {
			params[key] = topic.MustGet().Hex()
		}
	}
	var setOpr = func(key string, b mo.Option[bool]) {
		if b.IsPresent() {
			if b.MustGet() {
				params[key] = "and"
			} else {
				params[key] = "or"
			}
		}
	}
	setTopic("topic0", opts.Topic0)
	setTopic("topic1", opts.Topic1)
	setTopic("topic2", opts.Topic2)
	setTopic("topic3", opts.Topic3)
	setOpr("topic0_1_opr", opts.Topic01Opr)
	setOpr("topic1_2_opr", opts.Topic12Opr)
	setOpr("topic2_3_opr", opts.Topic23Opr)
	setOpr("topic0_2_opr", opts.Topic02Opr)
	setOpr("topic0_3_opr", opts.Topic03Opr)
	setOpr("topic1_3_opr", opts.Topic13Opr)
	return params
}

type LogEntrySlice []LogEntry

func (s LogEntrySlice) FindByTopic0(topic0 [32]byte) LogEntrySlice {
	return lo.Filter(s, func(item LogEntry, index int) bool {
		return item.Topics[0] == topic0
	})
}

type LogEntry struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	Data             HexBs          `json:"data"`
	BlockNumber      Uint64         `json:"blockNumber"`
	BlockHash        common.Hash    `json:"blockHash"`
	Timestamp        Uint64         `json:"timeStamp"`
	GasPrice         Uint64         `json:"gasPrice"`
	GasUsed          Uint64         `json:"gasUsed"`
	LogIndex         Uint64         `json:"logIndex"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex Uint64         `json:"transactionIndex"`
}

func (e *LogEntry) ToLog() *ethTypes.Log {
	return &ethTypes.Log{
		Address:     e.Address,
		Topics:      e.Topics,
		Data:        e.Data,
		BlockNumber: uint64(e.BlockNumber),
		TxHash:      e.TransactionHash,
		TxIndex:     uint(e.TransactionIndex),
		BlockHash:   e.BlockHash,
		Index:       uint(e.LogIndex),
		Removed:     false,
	}
}
