package code

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

// SerializableWallet Solve gob: type elliptic.p256Curve has no exported fields
type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
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
func (wallets *Wallets) GetWallet(address string) Wallet {
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
	var wallets_ map[string]SerializableWallet
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets_) //解码
	if err != nil {
		log.Panic(err)
	}
	wallets.Wallets = make(map[string]*Wallet)
	for k, v := range wallets_ {
		wallets.Wallets[k] = &Wallet{
			PrivateKey: ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}
	return nil
}

func (wallets *Wallets) SaveToFile() {
	mywalletfile := walletfile
	if err := os.MkdirAll(filepath.Dir(mywalletfile), 0700); err != nil {
		log.Panic(err)
	}

	var content bytes.Buffer

	gob.Register(SerializableWallet{})

	wallets_ := make(map[string]SerializableWallet)
	for k, v := range wallets.Wallets {
		wallets_[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}
	encoder := gob.NewEncoder(&content)
	if err := encoder.Encode(wallets_); err != nil {
		log.Panic(err)
	}

	// 安全优化：文件权限设为仅所有者读写‌:ml-citation{ref="3,6" data="citationList"}
	if err := os.WriteFile(mywalletfile, content.Bytes(), 0600); err != nil {
		log.Panic(err)
	}
}
