package block

import (
	"flag"
	"fmt"
)

type CLI struct {
	blockChain *BlockChain
}

func NewCLI(blockChain *BlockChain) *CLI {
	return &CLI{blockChain: blockChain}
}

func (cli *CLI) Run() {
	var data string
	var list bool
	var clear bool
	// 将命令行参数 -add 的值绑定到 data 变量中
	flag.StringVar(&data, "add", "", "添加区块")
	flag.BoolVar(&list, "list", false, "显示所有区块")
	flag.BoolVar(&clear, "clear", false, "删除所有区块")
	// 解析命令行参数写入注册的flag里
	flag.Parse()

	blockChain := cli.blockChain

	switch {
		case data != "":
			// 添加区块
			blockChain.AddBlock(data)
			fmt.Println("添加区块成功!!!")
		case list:
			// 打印区块
			iterator := blockChain.Iterator()

			blockData := iterator.Next()
			var i int
			for ; blockData != nil ; blockData = iterator.Next() {
				i++
				fmt.Printf("=================%d===================\n", i)
				fmt.Printf("Version: %d\n", blockData.Version)
				fmt.Printf("PrevHash: %x\n", blockData.PrevHash)
				fmt.Printf("MerKleRoot: %x\n", blockData.MerKleRoot)
				fmt.Printf("Timestamp: %d\n", blockData.Timestamp)
				fmt.Printf("Difficulty: %d\n", blockData.Difficulty)
				fmt.Printf("Nonce: %d\n", blockData.Nonce)
				fmt.Printf("Data: %s\n", blockData.Data)
				fmt.Printf("Hash: %x\n", blockData.Hash)

				pow := NewProofOfWork(blockData)
				fmt.Printf("IsValid: %v\n", pow.IsValid())
			}
		case clear:
			// 删除区块
			blockChain.Clear()
			fmt.Println("删除成功!!!")
		default:
			flag.PrintDefaults()
	}
}
