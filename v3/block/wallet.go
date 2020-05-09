package block

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"os"
)

const WalletFilename = "wallet.dat"

// 定义钱包结构
type Wallets struct {
	WalletMap map[string]*WalletKeyPair
}

// 创建钱包
func NewWallets() *Wallets {
	var wallets *Wallets
	// 从文件中加载数据
	wallets, err := loadFromFile()
	if err != nil {
		panic(err)
	}
	return wallets
}

func (wallets *Wallets) CreateWallet() string {
	walletKeyPair := NewWalletKeyPair()
	address := walletKeyPair.GetAddress()

	wallets.WalletMap[address] = walletKeyPair

	// 保存到文件
	err := saveToFile(wallets)
	if err != nil {
		return ""
	}

	return address
}

func (wallets *Wallets) ListAddress() []string {
	var addresses []string
	for address, _ := range wallets.WalletMap {
		addresses = append(addresses, address)
	}
	return addresses
}

// 校验地址格式是否正确
func IsValidAddress(address string) bool {
	decodeInfo := base58.Decode(address)
	if len(decodeInfo) != 25 {
		fmt.Printf("地址长度不正确!\n")
		return false
	}
	i := len(decodeInfo)-4
	// 21 个字节
	payload := decodeInfo[:i]
	// 进行两次 hash
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])
	// 取前 4 个字节
	checksum1 := hash[:4]
	checksum2 := decodeInfo[i:]
	// 比较
	return bytes.Equal(checksum1, checksum2)
}

func saveToFile(wallets *Wallets) error {
	var buffer bytes.Buffer
	// 注册接口类型
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wallets)
	if err != nil {
		fmt.Printf("Wallet 序列化失败, err: %v\n", err)
		return err
	}
	content := buffer.Bytes()
	err = ioutil.WriteFile(WalletFilename, content, 0600)
	if err != nil {
		fmt.Printf("Wallet 保存失败, err: %v\n", err)
		return err
	}
	return nil
}

func loadFromFile() (*Wallets, error) {
	var wallets Wallets
	_, err := os.Stat(WalletFilename)
	if os.IsNotExist(err) {
		wallets = Wallets{make(map[string]*WalletKeyPair)}
		return &wallets, nil
	}
	content, err := ioutil.ReadFile(WalletFilename)
	if err != nil {
		fmt.Printf("读取文件失败, err: %v\n", err)
		return nil, err
	}
	// 注册接口类型
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)
	if err != nil {
		fmt.Printf("Wallet 反序列化失败, err: %v\n", err)
		return nil, err
	}
	return &wallets, nil
}
