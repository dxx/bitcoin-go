package block

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const firstData = "Go 区块链"

//// 区块链结构体
//type BlockChain struct {
//	Blocks []*Block
//}
//
//// 创建区块链函数
//func NewBlockChain() *BlockChain {
//	// 添加创世块
//	v1 := NewBlock(firstData, []byte{0x0000000000000000})
//	blockChain := BlockChain{
//		Blocks: []*Block{v1},
//	}
//	return &blockChain
//}
//
//// 添加区块方法
//func (blockChain *BlockChain) AddBlock(data string) {
//	// 获取最后一个区块
//	prevBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
//	v1 := NewBlock(data, prevBlock.Hash)
//	blockChain.Blocks = append(blockChain.Blocks, v1)
//}

const DBPath = "block_bolt.db"
const LastHashKey = "last_block_hash"
const BucketName = "block_bucket"

// 区块链结构体
type BlockChain struct {
	boltDB        *bolt.DB // bolt 数据库句柄
	lastBlockHash []byte   // 最后一个区块的 hash
}

// 创建区块链函数
func NewBlockChain(miner string) *BlockChain {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	bucketName := []byte(BucketName)

	// 创建 Bucket
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		// 判断 Bucket 是否不存在
		if bucket == nil {
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	var lastBlockHash []byte

	// 存入数据
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)

		if bucket.Get([]byte(LastHashKey)) != nil {
			log.Fatal("区块链已存在!")
		}

		// 添加创世块
		// 创世块只有挖矿交易
		coinBase := NewCoinBaseTx(miner, firstData)
		block := NewBlock([]*Transaction{coinBase}, []byte{0x0000000000000000})
		err := bucket.Put(block.Hash, block.ToBytes())
		// 存入最后一个区块 hash
		err = bucket.Put([]byte(LastHashKey), block.Hash)

		lastBlockHash = block.Hash
		return err
	})

	blockChain := BlockChain{
		boltDB:        db,
		lastBlockHash: lastBlockHash,
	}
	return &blockChain
}

// 获取区块链函数
func GetBlockChain() *BlockChain {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var lastBlockHash []byte

	// 获取数据
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		// 判断 Bucket 是否不存在
		if bucket == nil {
			log.Fatal("区块链不存在")
		}
		lastBlockHash = bucket.Get([]byte(LastHashKey))
		return nil
	})

	blockChain := BlockChain{
		boltDB:        db,
		lastBlockHash: lastBlockHash,
	}
	return &blockChain
}

// 添加区块方法
func (blockChain *BlockChain) AddBlock(txs []*Transaction) {
	// 得到交易后第一时间对交易进行校验，过滤调无效交易
	var validTxs []*Transaction
	for _, tx := range txs {
		if blockChain.VerifyTransaction(tx) {
			validTxs = append(validTxs, tx)
		} else {
			fmt.Printf("发现无效的交易: %x\n", tx.TxId)
		}
	}

	// 存入数据
	blockChain.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		// 添加区块
		block := NewBlock(validTxs, blockChain.lastBlockHash)
		err := bucket.Put(block.Hash, block.ToBytes())
		// 存入最后一个区块 hash
		err = bucket.Put([]byte(LastHashKey), block.Hash)
		blockChain.lastBlockHash = block.Hash
		return err
	})
}

// 查找 UTXO
func (blockChain *BlockChain) FindMyUTXOs(publicKeyHash []byte) []UTXOInfo {
	var UTXOInfos []UTXOInfo
	// 0x111 => {0, 1}
	spentUTXOs := make(map[string][]int)
	it := blockChain.Iterator()
	// 遍历区块
	for block := it.Next() ; block != nil ; block = it.Next() {
		// 遍历交易
		for _, tx := range block.Transactions {
			// 挖矿机交易，跳过
			if tx.IsCoinBase() == false {
				// 遍历 input
				for _, input := range tx.TxInputs {
					if bytes.Equal(HashPublicKey(input.PublicKey), publicKeyHash) {
						key := string(input.TxId)
						// 保存输入脚本对应输出脚本的 index，使用 txID 作为 key
						spentUTXOs[key] = append(spentUTXOs[key], input.Index)
					}
				}
			}

			// 遍历 output
		OUTPUT:
			for i, output := range tx.TxOutputs {
				// 过滤 消耗的 output
				key := string(tx.TxId)
				indexes := spentUTXOs[key]
				if len(indexes) > 0 {
					for _, j := range indexes {
						if j == i {
							// 忽略
							continue OUTPUT
						}
					}
				}

				// 查找属于 address 的 output
				if bytes.Equal(output.PublicKeyHash, publicKeyHash) {
					UTXOInfo := UTXOInfo{tx.TxId, i, output}
					UTXOInfos = append(UTXOInfos, UTXOInfo)
				}
			}
		}
	}
	return UTXOInfos
}

// 获取余额
func (blockChain *BlockChain) GetBalance(address string) {
	publicKeyHash := Lock(address)
	UTXOInfos := blockChain.FindMyUTXOs(publicKeyHash)
	total := 0.0
	for _, UTXOInfo := range UTXOInfos {
		total += UTXOInfo.Output.Value
	}
	fmt.Printf("%s的余额为%f", address, total)
}

// 遍历账本，找到属于付款人的合适金额
func (blockChain *BlockChain) FindNeedUTXOs(publicKeyHash []byte, amount float64) (map[string][]int, float64) {
	UTXOs := make(map[string][]int)
	resValue := 0.0

	UTXOInfos := blockChain.FindMyUTXOs(publicKeyHash)
	for _, UTXOInfo := range UTXOInfos {
		key := string(UTXOInfo.TxId)
		UTXOs[key] = append(UTXOs[key], UTXOInfo.Index)
		resValue += UTXOInfo.Output.Value
		// 判断金额是否能够进行交易
		if resValue >= amount {
			// 可以进行交易，跳出循环
			break
		}
	}

	// 返回 UTXO 和 金额
	return UTXOs, resValue
}

// 交易签名
func (blockChain *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	fmt.Printf("对交易进行签名...\n")

	// 如果是挖矿交易，直接返回
	if tx.IsCoinBase() {
		return
	}

	prevTxs := blockChain.FindTransaction(tx)

	tx.Sign(privateKey, prevTxs)
}

// 校验签名
func (blockChain *BlockChain) VerifyTransaction(tx *Transaction) bool {
	fmt.Printf("对交易进行校验...\n")

	// 如果是挖矿交易，直接返回 true
	if tx.IsCoinBase() {
		return true
	}

	prevTxs := blockChain.FindTransaction(tx)

	return tx.Verify(prevTxs)
}

// 查找 input 引用的交易信息
func (blockChain *BlockChain) FindTransaction(tx *Transaction) map[string]*Transaction {
	prevTxs := make(map[string]*Transaction)
	it := blockChain.Iterator()
	// 遍历区块
	for block := it.Next() ; block != nil ; block = it.Next() {
		// 遍历交易
		for _, transaction := range block.Transactions {
			// 遍历输入
			for _, input := range tx.TxInputs {
				if bytes.Equal(input.TxId, transaction.TxId) {
					prevTxs[string(input.TxId)] = tx
					fmt.Printf("找到交易 %x\n", input.TxId)
				} else {
					fmt.Printf("未找到交易 %x\n", input.TxId)
				}
			}
		}
	}
	return prevTxs
}

// 获取迭代器
func (blockChain *BlockChain) Iterator() *BlockChainIterator {
	return NewBlockChainIterator(blockChain)
}

func (blockChain *BlockChain) Release() {
	_ = blockChain.boltDB.Close()
}

func (blockChain *BlockChain) Clear() {
	blockChain.boltDB.Close()
	err := os.Remove(DBPath)
	if err != nil {
		log.Fatal(err)
	}
}

// 迭代器
type BlockChainIterator struct {
	boltDB      *bolt.DB
	currentHash []byte
}

func NewBlockChainIterator(blockChain *BlockChain) *BlockChainIterator {
	return &BlockChainIterator{
		boltDB:      blockChain.boltDB,
		currentHash: blockChain.lastBlockHash,
	}
}

func (it *BlockChainIterator) Next() *Block {
	var block *Block
	it.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		value := bucket.Get(it.currentHash)
		if value != nil {
			block = &Block{}
			block.ToBlock(value)
			// 更新当前 hash
			it.currentHash = block.PrevHash
		}
		return nil
	})
	return block
}

