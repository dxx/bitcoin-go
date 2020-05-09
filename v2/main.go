package main

import (
	"bitcoin-go/v2/block"
)

func main() {
	cli := block.CLI{}
	cli.Run()

	//blockChain.AddBlock("新的区块")
	//
	//iterator := blockChain.Iterator()
	//
	//blockData := iterator.Next()
	//for ; blockData != nil ; blockData = iterator.Next() {
	//	fmt.Printf("====================================\n")
	//	fmt.Printf("Version: %d\n", blockData.Version)
	//	fmt.Printf("PrevHash: %x\n", blockData.PrevHash)
	//	fmt.Printf("MerKleRoot: %x\n", blockData.MerKleRoot)
	//	fmt.Printf("Timestamp: %d\n", blockData.Timestamp)
	//	fmt.Printf("Difficulty: %d\n", blockData.Difficulty)
	//	fmt.Printf("Nonce: %d\n", blockData.Nonce)
	//	fmt.Printf("Data: %s\n", blockData.Data)
	//	fmt.Printf("Hash: %x\n", blockData.Hash)
	//
	//	pow := v1.NewProofOfWork(blockData)
	//	fmt.Printf("IsValid: %v\n", pow.IsValid())
	//}
	//
}
