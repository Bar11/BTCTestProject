package code

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// 定义区块

type Block struct {
	Timestamp     int64          //时间线
	Transactions  []*Transaction //交易的集合
	PrevBlockHash []byte         //上一块数据的哈希
	Hash          []byte         //当前数据块的哈希
	Nonce         int            //工作量证明
}

// 对于交易实现哈希计算
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// 设定结构体对象的哈希
func (block *Block) SetHash() {
	// 处理当前时间，转化为十进制字符串在转化为字节集合
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	// 叠加要哈希的数据
	headers := bytes.Join([][]byte{block.PrevBlockHash, timestamp}, []byte{})
	// 计算出哈希地址
	hash := sha256.Sum256(headers)
	// 设置哈希
	block.Hash = hash[:]
}

// 创建一个区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	// block是一个指针，取得一个对象初始化之后的地址
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	block.SetHash()
	return block

}

// 创建一个创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 对象转化为二进制字节集，可以写入文件
func (block *Block) Serialize() []byte {
	var result bytes.Buffer            //开辟内存，存放字节集合
	encoder := gob.NewEncoder(&result) //编码对象创建
	err := encoder.Encode(block)       //编码操作
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 读取文件，读到二进制字节集，二进制字节集转化为对象
func Deserialize(data []byte) (*Block, error) {
	var block Block                                  // 对象存储用于字节转化的对象
	decoder := gob.NewDecoder(bytes.NewReader(data)) //解码
	err := decoder.Decode(&block)                    //尝试解码
	if err != nil {
		log.Panic(err)
	}
	return &block, nil
}

//func (block *Block) HashTransactions() []byte {
//
//}
