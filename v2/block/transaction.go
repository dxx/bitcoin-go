package block

import (
    "bytes"
    "crypto/sha256"
    "encoding/gob"
    "fmt"
    "log"
)

// 交易输入
// 指明交易发起人可支付金额的来源
// 引用 UTXO 的交易 id
// 消费 UTXO 在 output 中的索引
// 解锁脚本（签名，公钥）

// 交易输出
// 包含资金接收方的相关信息
// 交易金额
// 锁定脚本（对方公钥的哈希值）

// 输入交易
type TxInput struct {
    TxId    []byte // 交易 id
    Index   int    // output 的索引
    Address string // 解锁脚本，先使用地址来模拟
}

// 输出交易
type TxOutput struct {
    Value   float64 // 转账金额
    Address string  // 锁定脚本
}

// 交易结构
type Transaction struct {
    TxId      []byte     // 交易 id
    TxInputs  []TxInput  // 所有 inputs
    TxOutputs []TxOutput // 所有 outputs
}

// UTXO 结构
type UTXOInfo struct {
    TxId   []byte   // 交易 id
    Index  int      // output 索引
    Output TxOutput // output
}

// 设置交易 id
// 对 Transaction 进行 hash 运算
func (tx *Transaction) SetTxID() {
    // 使用 gob 对 tx 序列化
    var buffer bytes.Buffer
    encoder := gob.NewEncoder(&buffer)
    err := encoder.Encode(tx)
    if err != nil {
        log.Panic(err)
    }
    hash := sha256.Sum256(buffer.Bytes())
    tx.TxId = hash[:]
}

// 判断是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
    // 挖矿交易特点
    // 1.只有一个 input
    // 2.txId 为 nil
    // 3.Index 等于 -1
    if len(tx.TxInputs) == 1 && tx.TxInputs[0].TxId == nil &&
        tx.TxInputs[0].Index == -1 {
        return true
    }
    return false
}

// 挖矿交易
// 传入挖矿人
func NewCoinBaseTx(miner string, data string) *Transaction {
    // 在之后的程序中需要识别一个交易是否为 CoinBase ，所以初始化一些特殊值
    inputs := []TxInput{{nil, -1, data}}
    outputs := []TxOutput{{12.5, miner}}
    tx := &Transaction{nil, inputs, outputs}
    tx.SetTxID()
    return tx
}

// 创建普通交易
// 1.遍历账本，找到输入付款人合适的金额，即对应的 outputs
// 2.如果金额不足以转账，创建交易失败
// 3.将 outputs 转成 inputs
// 4.创建属于收款人的 output
// 5.如果有找零，创建属于付款人的 output
// 6.设置交易 id
// 7.返回交易结构
func NewTransaction(from, to string, amount float64, blockChain *BlockChain) *Transaction {
    // 能用的 UTXO
    UTXOs := make(map[string][]int)
    // UTXO 存储的金额
    resValue := 0.0
    UTXOs, resValue = blockChain.FindNeedUTXOs(from, amount)

    // 金额不足以转账，创建交易失败
    if resValue < amount {
        fmt.Println("余额不足，交易失败!")
        return nil
    }

    var inputs []TxInput
    // 将 outputs 转成 inputs
    // UTXOs: 0x111 => {0, 1}
    for txId, indexes := range UTXOs {
        for _, index := range indexes {
            inputs = append(inputs, TxInput{[]byte(txId), index, from})
        }
    }
    var outputs []TxOutput
    // 创建属于收款人的 output
    output := TxOutput{amount, to}
    outputs = append(outputs, output)

    if resValue > amount {
        // 如果有找零，创建属于付款人的 output
        outputs = append(outputs, TxOutput{resValue - amount, from})
    }

    tx := &Transaction{nil, inputs, outputs}
    // 设置交易 id
    tx.SetTxID()
    return tx
}
