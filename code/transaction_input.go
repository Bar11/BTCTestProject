package code

import "bytes"

// 输入
type TXInput struct {
	Txid      []byte //c存储了交易的id
	Vout      int    // 保存该交易中的一个output索引
	Signature []byte //签名
	PubKey    []byte // 公钥
}

// key 监测地址与交易
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockinghash := HashPubKey(in.PubKey)
	return bytes.Compare(lockinghash, pubKeyHash) == 0
}
