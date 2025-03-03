package code

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// 新建一个钱包，或者获取已存在的钱包
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}

func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	wallets.Wallets[address] = wallet
	return address
}

// 获取所有钱包地址
func (wallets *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// 获取一个钱包
func (wallets *Wallets) GetWallte(address string) Wallet {
	return *wallets.Wallets[address]
}

// 从文件中读取钱包
func (wallets *Wallets) LoadFromFile() error {

	mywalletfile := walletfile
	if _, err := os.Stat(mywalletfile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(mywalletfile)
	if err != nil {
		log.Panic(err)
	}
	// 读取二进制文件并解析
	var wallets_ Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets) //解码
	if err != nil {
		log.Panic(err)
	}
	wallets.Wallets = wallets_.Wallets
	return nil
}

// 钱包保存到文件
func (wallets *Wallets) SaveToFile() {
	var content bytes.Buffer
	mywalletfile := walletfile
	gob.Register(elliptic.P256()) // 注册加密算法
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(mywalletfile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
