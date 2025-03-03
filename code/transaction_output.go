package code

import "bytes"

// 输出
type TXOutput struct {
	Value      int //保存了币
	PubKeyHash []byte
}

// 输出锁住的标志
func (out *TXOutput) Lock(address []byte) {
	pubkeyhash := Base58Decode(address)            //编码
	pubkeyhash = pubkeyhash[1 : len(pubkeyhash)-4] //截取有效哈希
	out.PubKeyHash = pubkeyhash
}

// 监测是否被Key锁住
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// 新建一个输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address)) // 加锁
	return txo
}
