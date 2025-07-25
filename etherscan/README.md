# Etherscan Logs API Client

这是一个用于调用 Etherscan Logs API 的 Go 客户端库，使用 resty 发起 HTTP 请求。

## 功能特性

- ✅ 支持根据地址获取事件日志
- ✅ 支持根据主题获取事件日志  
- ✅ 支持根据地址和主题组合获取事件日志
- ✅ 智能参数处理（使用 `mo.Option` 处理可选参数）
- ✅ 完整的错误处理
- ✅ 类型安全的 API 响应
- ✅ 完备的单元测试
- ✅ 并发安全

## 安装

```bash
go get github.com/go-resty/resty/v2
go get github.com/samber/mo
go get github.com/stretchr/testify
```

## 快速开始

### 1. 创建客户端

```go
import "github.com/ubuntu/cuti/eths/etherscan"

client := etherscan.NewLogsClient("your-api-key")
```

### 2. 根据地址获取日志

```go
import "github.com/samber/mo"

opts := etherscan.GetLogsByAddressOptions{
    Address:   "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
    FromBlock: mo.Some(12878196),  // 可选参数
    ToBlock:   mo.Some(12878196),  // 可选参数
    Page:      mo.Some(1),         // 可选参数
    Offset:    mo.Some(1000),      // 可选参数
}

result, err := client.GetLogsByAddress(opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("获取到 %d 个日志\n", len(result.Result))
```

### 3. 根据主题获取日志

```go
opts := etherscan.GetLogsByTopicsOptions{
    FromBlock:  12878196,  // 必需参数
    ToBlock:    12879196,  // 必需参数
    Topic0:     mo.Some("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
    Topic1:     mo.Some("0x0000000000000000000000000000000000000000000000000000000000000000"),
    Topic01Opr: mo.Some("and"),  // 主题0和1之间的操作符
}

result, err := client.GetLogsByTopics(opts)
```

### 4. 根据地址和主题获取日志

```go
opts := etherscan.GetLogsByAddressAndTopicsOptions{
    Address:    "0x59728544b08ab483533076417fbbb2fd0b17ce3a",  // 必需
    FromBlock:  15073139,  // 必需
    ToBlock:    15074139,  // 必需
    Topic0:     mo.Some("0x27c4f0403323142b599832f26acd21c74a9e5b809f2215726e244a4ac588cd7d"),
    Topic1:     mo.Some("0x00000000000000000000000023581767a106ae21c074b2276d25e5c3e136a68b"),
    Topic01Opr: mo.Some("and"),
}

result, err := client.GetLogsByAddressAndTopics(opts)
```

## API 参数说明

### 必需参数 vs 可选参数

使用 `mo.Option[T]` 类型明确区分必需参数和可选参数：

- **必需参数**: 直接使用基础类型（如 `string`, `int`）
- **可选参数**: 使用 `mo.Option[T]` 类型

```go
// 必需参数
Address: "0x123..."     

// 可选参数
Page: mo.Some(1)        // 有值
Offset: mo.None[int]()  // 无值（将被忽略）
```

### 支持的主题操作符

在使用多个主题时，可以指定它们之间的逻辑关系：

- `topic0_1_opr`: 主题0和1之间的操作符
- `topic1_2_opr`: 主题1和2之间的操作符  
- `topic2_3_opr`: 主题2和3之间的操作符
- `topic0_2_opr`: 主题0和2之间的操作符
- `topic0_3_opr`: 主题0和3之间的操作符
- `topic1_3_opr`: 主题1和3之间的操作符

操作符值: `"and"` 或 `"or"`

## 数据结构

### LogEntry

```go
type LogEntry struct {
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
```

### LogsResponse

```go
type LogsResponse struct {
    Status  string     `json:"status"`
    Message string     `json:"message"`
    Result  []LogEntry `json:"result"`
}
```

## 错误处理

客户端提供详细的错误信息：

```go
result, err := client.GetLogsByAddress(opts)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "address parameter is required"):
        // 参数验证错误
    case strings.Contains(err.Error(), "API returned error"):
        // API 返回的业务错误
    case strings.Contains(err.Error(), "failed to make request"):
        // 网络请求错误
    default:
        // 其他错误
    }
}
```

## 运行测试

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=.
```

## 完整示例

查看 `example/main.go` 文件获取完整的使用示例。

## 限制说明

- 每次查询最多返回 1000 条记录
- 需要有效的 Etherscan API 密钥
- 支持以太坊主网和测试网（通过 ChainID 参数）
- 默认使用以太坊主网（ChainID = 1）

## 依赖项

- `github.com/go-resty/resty/v2`: HTTP 客户端
- `github.com/samber/mo`: 可选值处理
- `github.com/stretchr/testify`: 测试工具（仅测试时）