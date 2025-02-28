package code

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

const version = byte(0x00)
const walletfile = "wallet.dat"

// 检测地址长度
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// 创建一个钱包
func NewWallet() *Wallet {
	privatekey, publickey := newKeyPair()
	wallet := Wallet{privatekey, publickey}
	return &wallet
}

// 生成公钥私钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	// 生成密钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publickey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publickey
}

// 公钥的校验
func checkSum(data []byte) []byte {
	firstSHA := sha256.Sum256(data) // 加密校验
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

// 公钥的哈希处理
func HashPubKey(pubkey []byte) []byte {
	publicsha256 := sha256.Sum256(pubkey)     // 处理公钥
	R160Hash := crypto.RIPEMD160.New()        // 创建一个哈希算法对象
	_, err := R160Hash.Write(publicsha256[:]) // 写入处理
	if err != nil {
		log.Panic(err)
	}
	publicR160Hash := R160Hash.Sum(nil) // 叠加运算
	return publicR160Hash
}

// 获取钱包地址
func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checkSum(versionedPayload) // 检测版本与公钥
	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

// 校验钱包地址
func ValidateAddress(address string) bool {
	publicHash := Base58Decode([]byte(address))
	actualChecksum := publicHash[len(publicHash)-addressChecksumLen:]
	version := publicHash[0]
	publicHash = publicHash[1 : len(publicHash)-addressChecksumLen]
	targetCheckSum := checkSum(append([]byte{version}, publicHash...))
	return bytes.Compare(actualChecksum, targetCheckSum) == 0
}
