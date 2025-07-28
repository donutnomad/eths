package main

import (
	"fmt"
	"github.com/donutnomad/eths/etherscan"
	"github.com/samber/mo"
	"log"
)

func main() {
	// 创建 Logs 客户端
	client := etherscan.NewEtherscanClient("YourApiKeyToken")

	// 示例1: 根据地址获取事件日志
	fmt.Println("=== 示例1: 根据地址获取事件日志 ===")
	addressOpts := etherscan.GetLogsByAddressOptions{
		Address:   "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
		FromBlock: mo.Some(uint64(12878196)),
		ToBlock:   mo.Some(uint64(12878196)),
		Page:      mo.Some(1),
		Offset:    mo.Some(1000),
	}

	result1, err := client.GetLogsByAddress(addressOpts)
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("状态: %s\n", result1.Status)
		fmt.Printf("消息: %s\n", result1.Message)
		fmt.Printf("日志数量: %d\n", len(result1.Result))
		if len(result1.Result) > 0 {
			fmt.Printf("第一个日志地址: %s\n", result1.Result[0].Address)
			fmt.Printf("交易哈希: %s\n", result1.Result[0].TransactionHash)
		}
	}

	fmt.Println()

	// 示例2: 根据主题获取事件日志
	fmt.Println("=== 示例2: 根据主题获取事件日志 ===")
	topicsOpts := etherscan.GetLogsByTopicsOptions{
		FromBlock:  12878196,
		ToBlock:    12879196,
		Topic0:     mo.Some("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		Topic1:     mo.Some("0x0000000000000000000000000000000000000000000000000000000000000000"),
		Topic01Opr: mo.Some("and"),
		Page:       mo.Some(1),
		Offset:     mo.Some(1000),
	}

	result2, err := client.GetLogsByTopics(topicsOpts)
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("状态: %s\n", result2.Status)
		fmt.Printf("消息: %s\n", result2.Message)
		fmt.Printf("日志数量: %d\n", len(result2.Result))
		if len(result2.Result) > 0 {
			fmt.Printf("第一个日志的主题数量: %d\n", len(result2.Result[0].Topics))
		}
	}

	fmt.Println()

	// 示例3: 根据地址和主题获取事件日志
	fmt.Println("=== 示例3: 根据地址和主题获取事件日志 ===")
	addressTopicsOpts := etherscan.GetLogsByAddressAndTopicsOptions{
		Address:    "0x59728544b08ab483533076417fbbb2fd0b17ce3a",
		FromBlock:  15073139,
		ToBlock:    15074139,
		Topic0:     mo.Some("0x27c4f0403323142b599832f26acd21c74a9e5b809f2215726e244a4ac588cd7d"),
		Topic1:     mo.Some("0x00000000000000000000000023581767a106ae21c074b2276d25e5c3e136a68b"),
		Topic01Opr: mo.Some("and"),
		Page:       mo.Some(1),
		Offset:     mo.Some(1000),
	}

	result3, err := client.GetLogsByAddressAndTopics(addressTopicsOpts)
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("状态: %s\n", result3.Status)
		fmt.Printf("消息: %s\n", result3.Message)
		fmt.Printf("日志数量: %d\n", len(result3.Result))
		if len(result3.Result) > 0 {
			log := result3.Result[0]
			fmt.Printf("地址: %s\n", log.Address)
			fmt.Printf("区块号: %s\n", log.BlockNumber)
			fmt.Printf("时间戳: %s\n", log.TimeStamp)
			fmt.Printf("Gas 价格: %s\n", log.GasPrice)
			fmt.Printf("Gas 使用量: %s\n", log.GasUsed)
		}
	}

	fmt.Println()

	// 示例4: 最简单的用法（仅必需参数）
	fmt.Println("=== 示例4: 最简单的用法 ===")
	simpleOpts := etherscan.GetLogsByAddressOptions{
		Address: "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
	}

	result4, err := client.GetLogsByAddress(simpleOpts)
	if err != nil {
		log.Printf("错误: %v", err)
	} else {
		fmt.Printf("使用默认参数获取到 %d 个日志\n", len(result4.Result))
	}
}
