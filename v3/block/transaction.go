package block

import (
    "bytes"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "encoding/gob"
    "fmt"
    "log"
    "math/big"
    "strings"
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
    TxId  []byte // 交易 id
    Index int    // output 的索引
    // Address string // 解锁脚本，先使用地址来模拟
    Signature []byte // 签名
    PublicKey []byte // 公钥
}

// 输出交易
type TxOutput struct {
    Value float64 // 转账金额
    // Address string  // 锁定脚本
    PublicKeyHash []byte // 公钥哈希
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

const reward = 12.5

// 挖矿交易
// 传入挖矿人
func NewCoinBaseTx(miner string, data string) *Transaction {
    // 在之后的程序中需要识别一个交易是否为 CoinBase ，所以初始化一些特殊值
    //inputs := []TxInput{{nil, -1, data}}
    //outputs := []TxOutput{{12.5, miner}}
    inputs := []TxInput{{nil, -1, nil, []byte(data)}}
    outputs := []TxOutput{{reward, Lock(miner)}}

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

    // 获取钱包，找出公钥和私钥
    wallets := NewWallets()
    walletMap := wallets.WalletMap
    walletKeyPair := walletMap[from]
    if walletKeyPair == nil {
        fmt.Println("付款人地址错误，交易失败!")
        return nil
    }
    privateKey := walletKeyPair.PrivateKey
    publicKey := walletKeyPair.PublicKey

    publicKeyHash := Lock(from)

    UTXOs, resValue = blockChain.FindNeedUTXOs(publicKeyHash, amount)

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
            inputs = append(inputs, TxInput{[]byte(txId), index, nil, publicKey})
        }
    }
    var outputs []TxOutput
    // 创建属于收款人的 output
    output := TxOutput{amount, Lock(to)}
    outputs = append(outputs, output)

    if resValue > amount {
        // 如果有找零，创建属于付款人的 output
        outputs = append(outputs, TxOutput{resValue - amount, Lock(from)})
    }

    tx := &Transaction{nil, inputs, outputs}
    // 设置交易 id
    tx.SetTxID()

    blockChain.SignTransaction(tx, privateKey)

    return tx
}

// 签名
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, txs map[string]*Transaction) {
    fmt.Printf("签名...\n")
    // 复制交易
    copyTx := tx.Copy()
    // 遍历复制出来的 inputs 找到所引用的交易
    for i := 0; i < len(copyTx.TxInputs); i++ {
        input := copyTx.TxInputs[i]
        prevTx := txs[string(input.TxId)]
        // 将 PublicKeyHash 赋值给 PublicKey
        copyTx.TxInputs[i].PublicKey = prevTx.TxOutputs[input.Index].PublicKeyHash
        // 对交易进行 hash 运算，求出交易的 hash
        copyTx.SetTxID()
        signData := copyTx.TxId

        fmt.Printf("对数据 %x 进行签名\n", signData)

        // 将当前签完名的 PublicKey 设置为 nil
        input.PublicKey = nil

        // 对交易 hash 进行签名
        r, s, err := ecdsa.Sign(rand.Reader, privateKey, signData)
        if err != nil {
            fmt.Printf("对交易进行签名失败, err: %v\n", err)
        }
        // 拼接 r s
        signature := append(r.Bytes(), s.Bytes()...)
        // 赋值给原始交易的 Signature 字段
        tx.TxInputs[i].Signature = signature
    }
}

// 校验签名
func (tx *Transaction) Verify(txs map[string]*Transaction) bool {
    fmt.Printf("校验...\n")
    // 复制交易
    copyTx := tx.Copy()
    // 遍历 inputs 找到所引用的交易
    for i := 0; i < len(tx.TxInputs); i++ {
        input := tx.TxInputs[i]
        prevTx := txs[string(input.TxId)]
        // 将 PublicKeyHash 赋值给 PublicKey
        copyTx.TxInputs[i].PublicKey = prevTx.TxOutputs[input.Index].PublicKeyHash
        // 对交易进行 hash 运算，求出交易的 hash
        copyTx.SetTxID()
        // 这里的 verifyData 就是需要校验的数据
        verifyData := copyTx.TxId

        fmt.Printf("对数据 %x 进行校验\n", verifyData)

        // 将当前签完名的 PublicKey 设置为 nil
        copyTx.TxInputs[i].PublicKey = nil

        // 获取 signature
        signature := input.Signature
        // 裁切成签名后的 r 和 s
        var r big.Int
        var s big.Int
        r.SetBytes(signature[:len(signature) / 2])
        s.SetBytes(signature[len(signature) / 2:])

        // 获取使用了 gob 编码后的 publicKey
        pubKey := input.PublicKey
        // 反序列化成 PublicKey
        var publicKey *ecdsa.PublicKey
        gob.Register(elliptic.P256())
        decoder := gob.NewDecoder(bytes.NewReader(pubKey))
        err := decoder.Decode(&publicKey)
        if err != nil {
            fmt.Printf("反序列化 PublickKey 错误: %v\n", err)
        }
        // 校验
        if !ecdsa.Verify(publicKey, verifyData, &r, &s) {
            // 只要有一输入交易校验未通过就返回
            return false
        }

    }
    return true
}

// 复制交易
// 将每个输入交易 Signature 和 PublicKey 设置为 nil
// 输出交易不变
func (tx *Transaction) Copy() Transaction {
    var inputs []TxInput
    var outputs []TxOutput

    for _, input := range tx.TxInputs {
        txInput := TxInput{input.TxId, input.Index, nil, nil}
        inputs = append(inputs, txInput)
    }

    outputs = tx.TxOutputs

    return Transaction{tx.TxId, inputs, outputs}
}

// 定义 String 方法
func (tx *Transaction) String() string {
    var lines []string
    lines = append(lines, fmt.Sprintf("\n  Transaction %x:", tx.TxId))

    for i, txInput := range tx.TxInputs {
        lines = append(lines, fmt.Sprintf("    Input %d:", i))
        lines = append(lines, fmt.Sprintf("      TxId: %x", txInput.TxId))
        lines = append(lines, fmt.Sprintf("      OutIndex: %d", txInput.Index))
        lines = append(lines, fmt.Sprintf("      Signature: %x", txInput.Signature))
        lines = append(lines, fmt.Sprintf("      PublicKey: %x", txInput.PublicKey))
    }

    for i, txOutput := range tx.TxOutputs {
        lines = append(lines, fmt.Sprintf("    Output %d:", i))
        lines = append(lines, fmt.Sprintf("      Value: %f", txOutput.Value))
        lines = append(lines, fmt.Sprintf("      PublicKeyHash: %x", txOutput.PublicKeyHash))
    }

    return strings.Join(lines, "\n")
}
