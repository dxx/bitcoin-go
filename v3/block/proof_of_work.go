package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
)

// 用来推导难度值
const Bits = 16

// 工作量证明
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 创建工作量证明函数
func NewProofOfWork(block *Block) *ProofOfWork {
	// 难度值，应该是推导出来的，这里先定死
	// targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	// var targetBigInt big.Int
	// targetBigInt.SetString(targetStr, 16)

	// 推导难度值
	// 初始化十六进制数 1
	//  0000000000000000000000000000000000000000000000000000000000000001
	// 将对应的二进制位左移 256 位，（一个十六进制位对应 4 个二进制位）
	//1 0000000000000000000000000000000000000000000000000000000000000000
	// 再将对应的二进制位右移 16 位
	//0 0001000000000000000000000000000000000000000000000000000000000000

	targetBigInt := big.NewInt(1)
	// targetBigInt.Lsh(targetBigInt, 256)
	// targetBigInt.Rsh(targetBigInt, 16)
	targetBigInt.Lsh(targetBigInt, 256-Bits)
	pow := ProofOfWork{
		block:  block,
		target: targetBigInt,
	}
	return &pow
}

// 不断计算 hash
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var hash []byte
	var nonce uint64
	for {
		fmt.Printf("%x\r", hash)
		var bigIntTemp big.Int
		// v1 + nonce 计算 hash
		sha256Hash := sha256.Sum256(toBytes(pow.block, nonce))
		hash = sha256Hash[:]
		bigIntTemp.SetBytes(hash)
		if bigIntTemp.Cmp(pow.target) == -1 { // 判断当前 hash 是否比 target hash 小
			fmt.Printf("挖矿成功! nonce: %d, hash: %x\n", nonce, hash)
			break
		} else {
			// 否则 nonce++
			nonce++
		}
	}
	return hash, nonce
}

// 校验挖矿是否有效
func (pow *ProofOfWork) IsValid() bool {
	var bigInt big.Int
	bigInt.SetBytes(pow.block.Hash)
	return bigInt.Cmp(pow.target) == -1
}

func toBytes(block *Block, nonce uint64) []byte {
	dataBytes := [][]byte{
		uint64ToBytes(block.Version),
		block.PrevHash,
		block.MerKleRoot,
		uint64ToBytes(block.Timestamp),
		uint64ToBytes(block.Difficulty),
		uint64ToBytes(nonce),
	}
	data := bytes.Join(dataBytes, []byte{})
	sha256Bytes := sha256.Sum256(data)
	return sha256Bytes[:]
}

func uint64ToBytes(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
