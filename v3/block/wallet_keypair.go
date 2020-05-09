package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 创建密钥对，保存私钥和公钥
type WalletKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey []byte // 使用 gob 编码序列化的公钥
}

func NewWalletKeyPair() *WalletKeyPair {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	publicKey := privateKey.PublicKey

	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(publicKey)
	if err != nil {
		log.Panic(err)
	}

	return &WalletKeyPair{privateKey, buffer.Bytes()}
}

// 获取钱包地址
func (walletKeyPair *WalletKeyPair) GetAddress() string {
	var address string
	// 20 个字节
	publicKeyHash := HashPublicKey(walletKeyPair.PublicKey)

	// 1 个字节
	version := []byte{0x00}
	// 21 个字节
	payload := append(version, publicKeyHash...)

	// 进行两次 sha256
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	// 取前 4 个字节
	checksum := secondHash[:4]
	// 25 个字节
	payload = append(payload, checksum...)
	address = base58.Encode(payload)
	return address
}

func HashPublicKey(publicKey []byte) []byte {
	sha256Hash := sha256.Sum256(publicKey)

	ripe := ripemd160.New()
	ripe.Write(sha256Hash[:])
	// 20 个字节
	publicKeyHash := ripe.Sum(nil)
	return publicKeyHash
}

// 给定一个地址，反向推出公钥哈希
func Lock(address string) []byte {
	// 25 个字节
	b := base58.Decode(address)
	// 去掉第一个字节，和后四个字节
	publicKeyHash := b[1:len(b)-4]
	return publicKeyHash
}
