package block

import (
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
	boltDB *bolt.DB // bolt 数据库句柄
	lastBlockHash []byte // 最后一个区块的 hash
}

// 创建区块链函数
func NewBlockChain() *BlockChain {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	bucketName := []byte(BucketName)

	// 创建 Bucket
	db.Update(func (tx *bolt.Tx) error {
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
		
		value := bucket.Get([]byte(LastHashKey))
		// bolt 中没有区块则添加区块，否则取出最后一个区块 hash
		if value == nil {
			// 添加创世块
			block := NewBlock(firstData, []byte{0x0000000000000000})
			err := bucket.Put(block.Hash, block.ToBytes())
			// 存入最后一个区块 hash
			err = bucket.Put([]byte(LastHashKey), block.Hash)

			lastBlockHash = block.Hash
			return err
		} else {
			lastBlockHash = value
		}
		return nil
	})

	blockChain := BlockChain{
		boltDB: db,
		lastBlockHash: lastBlockHash,
	}
	return &blockChain
}

// 添加区块方法
func (blockChain *BlockChain) AddBlock(data string) {
	// 存入数据
	blockChain.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		// 添加区块
		block := NewBlock(data, blockChain.lastBlockHash)
		err := bucket.Put(block.Hash, block.ToBytes())
		// 存入最后一个区块 hash
		err = bucket.Put([]byte(LastHashKey), block.Hash)
		blockChain.lastBlockHash = block.Hash
		return err
	})
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
	boltDB *bolt.DB
	currentHash []byte
}

func NewBlockChainIterator(blockChain *BlockChain) *BlockChainIterator {
	return &BlockChainIterator{
		boltDB: blockChain.boltDB,
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
