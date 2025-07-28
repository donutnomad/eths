package main

import (
	"fmt"
	"log"
	"time"

	"github.com/donutnomad/eths/etherscan"
	"github.com/samber/mo"
)

func main() {
	// 创建 Logs 客户端
	client := etherscan.NewEtherscanClient("YourApiKeyToken")

	// 示例: 根据地址获取事件日志并展示解析后的数据类型
	fmt.Println("=== 优化后的数据类型示例 ===")
	opts := etherscan.GetLogsByAddressOptions{
		Address:   "0xbd3531da5cf5857e7cfaa92426877b022e612cf8",
		FromBlock: mo.Some(uint64(12878196)),
		ToBlock:   mo.Some(uint64(12878196)),
		Page:      mo.Some(1),
		Offset:    mo.Some(10),
	}

	result, err := client.GetLogsByAddress(opts)
	if err != nil {
		log.Printf("错误: %v", err)
		return
	}

	fmt.Printf("状态: %s\n", result.Status)
	fmt.Printf("消息: %s\n", result.Message)
	fmt.Printf("日志数量: %d\n", len(result.Result))

	if len(result.Result) > 0 {
		logEntry := result.Result[0]

		fmt.Println("\n=== 第一个日志的详细信息 ===")
		fmt.Printf("地址: %s\n", logEntry.Address)
		fmt.Printf("区块号: %d (十进制)\n", logEntry.BlockNumber)
		fmt.Printf("时间戳: %d (Unix时间戳)\n", logEntry.TimeStamp)

		// 将时间戳转换为可读格式
		timestamp := time.Unix(int64(logEntry.TimeStamp), 0)
		fmt.Printf("时间: %s\n", timestamp.Format("2006-01-02 15:04:05 UTC"))

		fmt.Printf("Gas 价格: %d wei\n", logEntry.GasPrice)
		fmt.Printf("Gas 价格: %.2f Gwei\n", float64(logEntry.GasPrice)/1e9)
		fmt.Printf("Gas 使用量: %d\n", logEntry.GasUsed)
		fmt.Printf("日志索引: %d\n", logEntry.LogIndex)
		fmt.Printf("交易索引: %d\n", logEntry.TransactionIndex)

		fmt.Printf("交易哈希: %s\n", logEntry.TransactionHash.String())

		// 展示 Data 字段的使用
		fmt.Printf("数据长度: %d 字节\n", len(logEntry.Data))
		if len(logEntry.Data) > 0 {
			fmt.Printf("数据内容 (hex): %x\n", logEntry.Data)
			fmt.Printf("数据内容 (带0x): 0x%x\n", logEntry.Data)
		} else {
			fmt.Printf("数据内容: 空 (无数据)\n")
		}

		fmt.Printf("主题数量: %d\n", len(logEntry.Topics))
		for i, topic := range logEntry.Topics {
			fmt.Printf("主题[%d]: %s\n", i, topic.String())
			fmt.Printf("主题[%d] (无前缀): %s\n", i, topic.Hex())
			fmt.Printf("主题[%d] 是否为零: %t\n", i, topic.IsZero())
		}

		// 类型安全的数学运算示例
		fmt.Println("\n=== 数学运算示例 ===")
		totalGasCost := logEntry.GasPrice * logEntry.GasUsed
		fmt.Printf("总 Gas 费用: %d wei\n", totalGasCost)
		fmt.Printf("总 Gas 费用: %.6f ETH\n", float64(totalGasCost)/1e18)

		// 区块时间比较
		fmt.Printf("区块是否在2021年之后: %t\n", logEntry.TimeStamp > 1609459200) // 2021-01-01 的时间戳

		// Hash 类型的高级使用
		fmt.Println("\n=== Hash 类型高级特性 ===")
		if len(logEntry.Topics) > 0 {
			topic := logEntry.Topics[0]
			fmt.Printf("第一个主题的字节表示: %x\n", topic.Bytes())
			fmt.Printf("第一个主题的长度: %d 字节\n", len(topic.Bytes()))

			// 比较哈希
			zeroHash := etherscan.Hash{}
			fmt.Printf("是否为零哈希: %t\n", topic == zeroHash)

			// 创建新的哈希进行比较
			if newHash, err := etherscan.NewHashFromHex(topic.String()); err == nil {
				fmt.Printf("哈希比较: %t\n", topic == newHash)
			}
		}
	}

	// 展示类型安全的好处
	fmt.Println("\n=== 类型安全示例 ===")
	if len(result.Result) > 0 {
		logEntry := result.Result[0]

		// 现在可以直接进行数值比较，无需字符串转换
		if logEntry.GasPrice > 50000000000 { // 50 Gwei
			fmt.Println("⚠️  这是一个高 Gas 价格的交易")
		} else {
			fmt.Println("✅ 这是一个正常 Gas 价格的交易")
		}

		// 区块时间范围检查
		oneHour := uint64(3600)
		currentTime := uint64(time.Now().Unix())
		if currentTime-logEntry.TimeStamp < oneHour {
			fmt.Println("🕒 这是一个最近的交易")
		} else {
			fmt.Println("🕰️  这是一个较早的交易")
		}

		// Hash 类型操作示例
		fmt.Println("\n=== Hash 类型操作示例 ===")
		txHash := logEntry.TransactionHash
		fmt.Printf("交易哈希完整形式: %s\n", txHash.String())
		fmt.Printf("交易哈希简洁形式: %s...%s\n",
			txHash.String()[:10],
			txHash.String()[len(txHash.String())-8:])

		// 检查特定的事件类型 (ERC-20 Transfer)
		if len(logEntry.Topics) > 0 {
			transferEventHash, _ := etherscan.NewHashFromHex("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
			if logEntry.Topics[0] == transferEventHash {
				fmt.Println("🎯 这是一个 ERC-20 Transfer 事件")

				// 解析 Transfer 事件的数据
				if len(logEntry.Data) >= 32 {
					// ERC-20 Transfer 的数据通常是 amount (32 bytes)
					fmt.Printf("转账金额 (raw bytes): %x\n", logEntry.Data[:32])

					// 将字节转换为大整数（这里简化展示前8字节作为uint64）
					if len(logEntry.Data) >= 8 {
						var amount uint64
						for i := 0; i < 8 && i < len(logEntry.Data); i++ {
							amount = (amount << 8) | uint64(logEntry.Data[len(logEntry.Data)-8+i])
						}
						fmt.Printf("转账金额 (部分解析): %d\n", amount)
					}
				}
			}
		}

		// Data 字段的实用工具函数示例
		fmt.Println("\n=== Data 字段工具函数示例 ===")
		if len(logEntry.Data) > 0 {
			// 检查数据是否全为零
			allZero := true
			for _, b := range logEntry.Data {
				if b != 0 {
					allZero = false
					break
				}
			}
			fmt.Printf("数据是否全为零: %t\n", allZero)

			// 获取数据的校验和（简单示例）
			checksum := byte(0)
			for _, b := range logEntry.Data {
				checksum ^= b
			}
			fmt.Printf("数据校验和 (XOR): 0x%02x\n", checksum)
		}
	}
}
