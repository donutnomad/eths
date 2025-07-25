package main

import (
	"fmt"
	"log"

	"github.com/donutnomad/eths/etherscan"
	"github.com/samber/mo"
)

func main() {
	// 创建交易客户端
	client := etherscan.NewTransactionClient("YourApiKeyToken")

	fmt.Println("=== Etherscan 交易状态 API 示例 ===")

	// 示例交易哈希
	successTxHash := "0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76"
	failedTxHash := "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a"

	// 示例1: 检查合约执行状态
	fmt.Println("\n=== 1. 检查合约执行状态 ===")

	// 检查成功的交易
	fmt.Printf("检查交易: %s\n", successTxHash)
	contractResult, err := client.GetContractExecutionStatus(etherscan.GetContractExecutionStatusOptions{
		TxHash: successTxHash,
	})
	if err != nil {
		log.Printf("获取合约执行状态失败: %v", err)
	} else {
		fmt.Printf("API 状态: %s\n", contractResult.Status)
		fmt.Printf("API 消息: %s\n", contractResult.Message)
		fmt.Printf("合约执行是否成功: %t\n", contractResult.Result.IsSuccess())
		fmt.Printf("合约执行是否失败: %t\n", contractResult.Result.HasError())
		if contractResult.Result.HasError() {
			fmt.Printf("错误描述: %s\n", contractResult.Result.ErrDescription)
		}
	}

	// 检查失败的交易
	fmt.Printf("\n检查交易: %s\n", failedTxHash)
	contractResult2, err := client.GetContractExecutionStatus(etherscan.GetContractExecutionStatusOptions{
		TxHash:  failedTxHash,
		ChainID: mo.Some(1), // 明确指定以太坊主网
	})
	if err != nil {
		log.Printf("获取合约执行状态失败: %v", err)
	} else {
		fmt.Printf("API 状态: %s\n", contractResult2.Status)
		fmt.Printf("API 消息: %s\n", contractResult2.Message)
		fmt.Printf("合约执行是否成功: %t\n", contractResult2.Result.IsSuccess())
		fmt.Printf("合约执行是否失败: %t\n", contractResult2.Result.HasError())
		if contractResult2.Result.HasError() {
			fmt.Printf("错误描述: %s\n", contractResult2.Result.ErrDescription)
		}
	}

	// 示例2: 检查交易收据状态（拜占庭分叉后）
	fmt.Println("\n=== 2. 检查交易收据状态 ===")
	fmt.Println("注意: 此 API 仅适用于拜占庭分叉后的交易")

	// 检查成功的交易收据
	fmt.Printf("检查交易收据: %s\n", successTxHash)
	receiptResult, err := client.GetTransactionReceiptStatus(etherscan.GetTransactionReceiptStatusOptions{
		TxHash: successTxHash,
	})
	if err != nil {
		log.Printf("获取交易收据状态失败: %v", err)
	} else {
		fmt.Printf("API 状态: %s\n", receiptResult.Status)
		fmt.Printf("API 消息: %s\n", receiptResult.Message)
		fmt.Printf("交易是否成功: %t\n", receiptResult.Result.IsSuccess())
		fmt.Printf("交易是否失败: %t\n", receiptResult.Result.HasError())

		// 根据结果提供用户友好的消息
		if receiptResult.Result.IsSuccess() {
			fmt.Println("✅ 交易执行成功")
		} else {
			fmt.Println("❌ 交易执行失败")
		}
	}

	// 示例3: 批量检查多个交易状态
	fmt.Println("\n=== 3. 批量检查交易状态 ===")

	txHashes := []string{
		"0x513c1ba0bebf66436b5fed86ab668452b7805593c05073eb2d51d3a52f480a76",
		"0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a",
		// 可以添加更多交易哈希
	}

	for i, txHash := range txHashes {
		fmt.Printf("\n检查交易 %d: %s\n", i+1, txHash)

		// 同时检查合约执行状态和交易收据状态
		contractStatus, contractErr := client.GetContractExecutionStatus(etherscan.GetContractExecutionStatusOptions{
			TxHash: txHash,
		})

		receiptStatus, receiptErr := client.GetTransactionReceiptStatus(etherscan.GetTransactionReceiptStatusOptions{
			TxHash: txHash,
		})

		fmt.Printf("交易哈希: %s...%s\n", txHash[:10], txHash[len(txHash)-6:])

		if contractErr == nil && contractStatus != nil {
			if contractStatus.Result.IsSuccess() {
				fmt.Printf("  合约执行: ✅ 成功\n")
			} else {
				fmt.Printf("  合约执行: ❌ 失败 - %s\n", contractStatus.Result.ErrDescription)
			}
		} else {
			fmt.Printf("  合约执行: ⚠️  检查失败 - %v\n", contractErr)
		}

		if receiptErr == nil && receiptStatus != nil {
			if receiptStatus.Result.IsSuccess() {
				fmt.Printf("  交易收据: ✅ 成功\n")
			} else {
				fmt.Printf("  交易收据: ❌ 失败\n")
			}
		} else {
			fmt.Printf("  交易收据: ⚠️  检查失败 - %v\n", receiptErr)
		}
	}

	// 示例4: 实用的交易状态检查函数
	fmt.Println("\n=== 4. 实用工具函数示例 ===")

	checkTransactionStatus := func(txHash string) {
		fmt.Printf("\n🔍 完整检查交易: %s\n", txHash)

		// 尝试获取合约执行状态
		contractResult, contractErr := client.GetContractExecutionStatus(etherscan.GetContractExecutionStatusOptions{
			TxHash: txHash,
		})

		// 尝试获取交易收据状态
		receiptResult, receiptErr := client.GetTransactionReceiptStatus(etherscan.GetTransactionReceiptStatusOptions{
			TxHash: txHash,
		})

		// 综合分析结果
		if contractErr == nil && receiptErr == nil {
			contractSuccess := contractResult.Result.IsSuccess()
			receiptSuccess := receiptResult.Result.IsSuccess()

			if contractSuccess && receiptSuccess {
				fmt.Println("🎉 交易完全成功!")
			} else if !contractSuccess && !receiptSuccess {
				fmt.Printf("💥 交易完全失败: %s\n", contractResult.Result.ErrDescription)
			} else {
				fmt.Println("⚠️  交易状态不一致，请进一步检查")
				fmt.Printf("   合约执行: %t, 交易收据: %t\n", contractSuccess, receiptSuccess)
			}
		} else {
			fmt.Println("⚠️  无法获取完整的交易状态信息")
			if contractErr != nil {
				fmt.Printf("   合约状态错误: %v\n", contractErr)
			}
			if receiptErr != nil {
				fmt.Printf("   收据状态错误: %v\n", receiptErr)
			}
		}
	}

	// 检查示例交易
	checkTransactionStatus(successTxHash)
	checkTransactionStatus(failedTxHash)

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 合约执行状态 API: 检查智能合约交易的执行状态")
	fmt.Println("   - IsError='0': 执行成功")
	fmt.Println("   - IsError='1': 执行失败，并提供错误描述")
	fmt.Println()
	fmt.Println("2. 交易收据状态 API: 检查交易的最终状态（拜占庭分叉后）")
	fmt.Println("   - Status='1': 交易成功")
	fmt.Println("   - Status='0': 交易失败")
	fmt.Println()
	fmt.Println("3. 建议同时使用两个 API 以获得最完整的交易状态信息")
}
