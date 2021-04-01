package block

import (
    "flag"
    "fmt"
    "strconv"
)

type CLI struct {
    //blockChain *BlockChain
}

func (cli *CLI) Run() {
    var data string
    var addr string
    var send bool
    var list bool
    var clear bool
    flag.StringVar(&data, "create-block-chain", "", "创建区块链")
    flag.StringVar(&addr, "get-balance", "", "获取余额")
    flag.BoolVar(&send, "send", false, "转账（付款人 收款人 转账金额 miner）")
    flag.BoolVar(&list, "list", false, "显示所有区块")
    flag.BoolVar(&clear, "clear", false, "删除所有区块")
    // 解析命令行参数写入注册的flag里
    flag.Parse()

    var blockChain *BlockChain

    switch {
        case data != "":
            // 创建区块链
            NewBlockChain(data)
            fmt.Println("创建区块链成功!!!")
        case addr != "":
            // 获取余额
            blockChain = GetBlockChain()
            blockChain.GetBalance(addr)
        case send:
            // 转账
            blockChain = GetBlockChain()
            args := flag.Args()
            if len(args) == 4 {
                sender := args[0]
                receiver := args[1]
                amount, _ := strconv.ParseFloat(args[2], 64)
                miner := args[3]

                // 创建挖矿交易
                coinBase := NewCoinBaseTx(miner, firstData)
                txs := []*Transaction{coinBase}
                // 普通交易
                tx := NewTransaction(sender, receiver, amount, blockChain)
                if tx != nil {
                    txs = append(txs, tx)
                } else {
                    fmt.Println("无效交易!")
                }

                // 添加区块
                blockChain.AddBlock(txs)
            } else {
                fmt.Println("参数错误!")
            }
        case list:
            // 打印区块
            blockChain = GetBlockChain()
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
                // fmt.Printf("Data: %s\n", blockData.Data) TODO
                fmt.Printf("Hash: %x\n", blockData.Hash)

                pow := NewProofOfWork(blockData)
                fmt.Printf("IsValid: %v\n", pow.IsValid())
            }
        case clear:
            // 删除区块
            blockChain = GetBlockChain()
            blockChain.Clear()
            fmt.Println("删除成功!!!")
        default:
            flag.PrintDefaults()
    }
    if blockChain != nil {
        blockChain.Release()
    }
}
