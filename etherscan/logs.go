package etherscan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/samber/mo"
)

// LogsClient 提供对 Etherscan Logs API 的访问
type LogsClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// NewLogsClient 创建一个新的 LogsClient 实例
func NewLogsClient(apiKey string) *LogsClient {
	return &LogsClient{
		client:  resty.New(),
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/v2/api",
	}
}

// Hash 表示一个 32 字节的哈希值
type Hash [32]byte

// NewHashFromHex 从十六进制字符串创建 Hash
func NewHashFromHex(hexStr string) (Hash, error) {
	var hash Hash

	// 移除 0x 前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// 检查长度，如果不足 64 个字符，在前面补零
	if len(hexStr) < 64 {
		hexStr = strings.Repeat("0", 64-len(hexStr)) + hexStr
	} else if len(hexStr) > 64 {
		return hash, fmt.Errorf("hex string too long: %d characters", len(hexStr))
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return hash, fmt.Errorf("invalid hex string: %w", err)
	}

	copy(hash[:], bytes)
	return hash, nil
}

// String 返回哈希的十六进制字符串表示（带 0x 前缀）
func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

// Hex 返回哈希的十六进制字符串表示（不带 0x 前缀）
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// Bytes 返回哈希的字节切片
func (h Hash) Bytes() []byte {
	return h[:]
}

// IsZero 检查哈希是否为零值
func (h Hash) IsZero() bool {
	return h == Hash{}
}

// MarshalJSON 实现 JSON 序列化
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 实现 JSON 反序列化
func (h *Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}

	hash, err := NewHashFromHex(hexStr)
	if err != nil {
		return err
	}

	*h = hash
	return nil
}

type LogEntrySlice []LogEntry

func (s LogEntrySlice) FindByTopic0(topic0 [32]byte) LogEntrySlice {
	return lo.Filter(s, func(item LogEntry, index int) bool {
		return item.Topics[0] == topic0
	})
}

// LogEntry 表示一个事件日志条目
type LogEntry struct {
	Address          string `json:"-"`
	Topics           []Hash `json:"-"`
	Data             []byte `json:"-"`
	BlockNumber      uint64 `json:"-"`
	TimeStamp        uint64 `json:"-"`
	GasPrice         uint64 `json:"-"`
	GasUsed          uint64 `json:"-"`
	LogIndex         uint64 `json:"-"`
	TransactionHash  Hash   `json:"-"`
	TransactionIndex uint64 `json:"-"`
}

// logEntryRaw 用于解析原始 JSON 响应
type logEntryRaw struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TimeStamp        string   `json:"timeStamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

// parseHexToUint64 将十六进制字符串转换为 uint64
func parseHexToUint64(hexStr string) (uint64, error) {
	if hexStr == "" || hexStr == "0x" {
		return 0, nil
	}

	// 移除 0x 前缀
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	if hexStr == "" {
		return 0, nil
	}

	return strconv.ParseUint(hexStr, 16, 64)
}

// parseHexToBytes 将十六进制字符串转换为字节切片
func parseHexToBytes(hexStr string) ([]byte, error) {
	if hexStr == "" || hexStr == "0x" {
		return []byte{}, nil
	}

	// 移除 0x 前缀
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}

	if hexStr == "" {
		return []byte{}, nil
	}

	// 如果长度为奇数，在前面添加一个0
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	return hex.DecodeString(hexStr)
}

// UnmarshalJSON 自定义 JSON 反序列化方法
func (l *LogEntry) UnmarshalJSON(data []byte) error {
	var raw logEntryRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Address = raw.Address

	// 解析 Data 字段
	if dataBytes, err := parseHexToBytes(raw.Data); err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	} else {
		l.Data = dataBytes
	}

	// 解析 TransactionHash
	if raw.TransactionHash != "" {
		transactionHash, err := NewHashFromHex(raw.TransactionHash)
		if err != nil {
			return fmt.Errorf("failed to parse transactionHash: %w", err)
		}
		l.TransactionHash = transactionHash
	}

	// 解析 Topics
	l.Topics = make([]Hash, len(raw.Topics))
	for i, topic := range raw.Topics {
		hash, err := NewHashFromHex(topic)
		if err != nil {
			return fmt.Errorf("failed to parse topic[%d]: %w", i, err)
		}
		l.Topics[i] = hash
	}

	// 解析十六进制数字
	var err error
	if l.BlockNumber, err = parseHexToUint64(raw.BlockNumber); err != nil {
		return fmt.Errorf("failed to parse blockNumber: %w", err)
	}

	if l.TimeStamp, err = parseHexToUint64(raw.TimeStamp); err != nil {
		return fmt.Errorf("failed to parse timeStamp: %w", err)
	}

	if l.GasPrice, err = parseHexToUint64(raw.GasPrice); err != nil {
		return fmt.Errorf("failed to parse gasPrice: %w", err)
	}

	if l.GasUsed, err = parseHexToUint64(raw.GasUsed); err != nil {
		return fmt.Errorf("failed to parse gasUsed: %w", err)
	}

	if l.LogIndex, err = parseHexToUint64(raw.LogIndex); err != nil {
		return fmt.Errorf("failed to parse logIndex: %w", err)
	}

	if l.TransactionIndex, err = parseHexToUint64(raw.TransactionIndex); err != nil {
		return fmt.Errorf("failed to parse transactionIndex: %w", err)
	}

	return nil
}

// LogsResponse 表示 API 响应
type LogsResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Result  []LogEntry `json:"result"`
}

// GetLogsByAddressOptions 用于根据地址获取日志的选项
type GetLogsByAddressOptions struct {
	Address   string         // 必需：要检查日志的地址
	FromBlock mo.Option[int] // 可选：开始搜索的区块号(包含)
	ToBlock   mo.Option[int] // 可选：停止搜索的区块号(包含)
	Page      mo.Option[int] // 可选：页码(最小为1)
	Offset    mo.Option[int] // 可选：每页显示的记录数，最大1000
	ChainID   mo.Option[int] // 可选：链ID，默认为1（以太坊主网）
}

// GetLogsByTopicsOptions 用于根据主题获取日志的选项
type GetLogsByTopicsOptions struct {
	FromBlock  int               // 必需：开始搜索的区块号
	ToBlock    int               // 必需：停止搜索的区块号
	Topic0     mo.Option[string] // 可选：主题0
	Topic1     mo.Option[string] // 可选：主题1
	Topic2     mo.Option[string] // 可选：主题2
	Topic3     mo.Option[string] // 可选：主题3
	Topic01Opr mo.Option[string] // 可选：主题0和1之间的操作符 (and|or)
	Topic12Opr mo.Option[string] // 可选：主题1和2之间的操作符 (and|or)
	Topic23Opr mo.Option[string] // 可选：主题2和3之间的操作符 (and|or)
	Topic02Opr mo.Option[string] // 可选：主题0和2之间的操作符 (and|or)
	Topic03Opr mo.Option[string] // 可选：主题0和3之间的操作符 (and|or)
	Topic13Opr mo.Option[string] // 可选：主题1和3之间的操作符 (and|or)
	Page       mo.Option[int]    // 可选：页码
	Offset     mo.Option[int]    // 可选：每页显示的记录数，最大1000
	ChainID    mo.Option[int]    // 可选：链ID，默认为1（以太坊主网）
}

// GetLogsByAddressAndTopicsOptions 用于根据地址和主题获取日志的选项
type GetLogsByAddressAndTopicsOptions struct {
	Address    string            // 必需：要检查日志的地址
	FromBlock  int               // 必需：开始搜索的区块号
	ToBlock    int               // 必需：停止搜索的区块号
	Topic0     mo.Option[string] // 可选：主题0
	Topic1     mo.Option[string] // 可选：主题1
	Topic2     mo.Option[string] // 可选：主题2
	Topic3     mo.Option[string] // 可选：主题3
	Topic01Opr mo.Option[string] // 可选：主题0和1之间的操作符 (and|or)
	Topic12Opr mo.Option[string] // 可选：主题1和2之间的操作符 (and|or)
	Topic23Opr mo.Option[string] // 可选：主题2和3之间的操作符 (and|or)
	Topic02Opr mo.Option[string] // 可选：主题0和2之间的操作符 (and|or)
	Topic03Opr mo.Option[string] // 可选：主题0和3之间的操作符 (and|or)
	Topic13Opr mo.Option[string] // 可选：主题1和3之间的操作符 (and|or)
	Page       mo.Option[int]    // 可选：页码
	Offset     mo.Option[int]    // 可选：每页显示的记录数，最大1000
	ChainID    mo.Option[int]    // 可选：链ID，默认为1（以太坊主网）
}

// buildQueryParams 构建查询参数的通用方法
func (c *LogsClient) buildQueryParams(params map[string]any) map[string]string {
	queryParams := map[string]string{
		"module": "logs",
		"action": "getLogs",
		"apikey": c.apiKey,
	}

	for key, value := range params {
		switch v := value.(type) {
		case string:
			if v != "" {
				queryParams[key] = v
			}
		case int:
			queryParams[key] = strconv.Itoa(v)
		case mo.Option[string]:
			if v.IsPresent() {
				queryParams[key] = v.MustGet()
			}
		case mo.Option[int]:
			if v.IsPresent() {
				queryParams[key] = strconv.Itoa(v.MustGet())
			}
		}
	}

	return queryParams
}

// GetLogsByAddress 根据地址获取事件日志
func (c *LogsClient) GetLogsByAddress(opts GetLogsByAddressOptions) (*LogsResponse, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("address parameter is required")
	}

	params := map[string]any{
		"address":   opts.Address,
		"fromBlock": opts.FromBlock,
		"toBlock":   opts.ToBlock,
		"page":      opts.Page,
		"offset":    opts.Offset,
		"chainid":   opts.ChainID.OrElse(1), // 默认为以太坊主网
	}

	queryParams := c.buildQueryParams(params)

	var response LogsResponse
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if response.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	return &response, nil
}

// GetLogsByTopics 根据主题获取事件日志
func (c *LogsClient) GetLogsByTopics(opts GetLogsByTopicsOptions) (*LogsResponse, error) {
	params := map[string]any{
		"fromBlock":    opts.FromBlock,
		"toBlock":      opts.ToBlock,
		"topic0":       opts.Topic0,
		"topic1":       opts.Topic1,
		"topic2":       opts.Topic2,
		"topic3":       opts.Topic3,
		"topic0_1_opr": opts.Topic01Opr,
		"topic1_2_opr": opts.Topic12Opr,
		"topic2_3_opr": opts.Topic23Opr,
		"topic0_2_opr": opts.Topic02Opr,
		"topic0_3_opr": opts.Topic03Opr,
		"topic1_3_opr": opts.Topic13Opr,
		"page":         opts.Page,
		"offset":       opts.Offset,
		"chainid":      opts.ChainID.OrElse(1),
	}

	queryParams := c.buildQueryParams(params)

	var response LogsResponse
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if response.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	return &response, nil
}

// GetLogsByAddressAndTopics 根据地址和主题获取事件日志
func (c *LogsClient) GetLogsByAddressAndTopics(opts GetLogsByAddressAndTopicsOptions) (*LogsResponse, error) {
	if opts.Address == "" {
		return nil, fmt.Errorf("address parameter is required")
	}

	params := map[string]any{
		"address":      opts.Address,
		"fromBlock":    opts.FromBlock,
		"toBlock":      opts.ToBlock,
		"topic0":       opts.Topic0,
		"topic1":       opts.Topic1,
		"topic2":       opts.Topic2,
		"topic3":       opts.Topic3,
		"topic0_1_opr": opts.Topic01Opr,
		"topic1_2_opr": opts.Topic12Opr,
		"topic2_3_opr": opts.Topic23Opr,
		"topic0_2_opr": opts.Topic02Opr,
		"topic0_3_opr": opts.Topic03Opr,
		"topic1_3_opr": opts.Topic13Opr,
		"page":         opts.Page,
		"offset":       opts.Offset,
		"chainid":      opts.ChainID.OrElse(1),
	}

	queryParams := c.buildQueryParams(params)

	var response LogsResponse
	resp, err := c.client.R().
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(c.baseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	if response.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	return &response, nil
}
