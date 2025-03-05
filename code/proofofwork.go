package code

import (
	"BTCTestProject/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 24 // 控制计算难度

type ProofOfWork struct {
	Block  *Block   // 区块
	target *big.Int // 存储计算哈希对比的特定整数
}

// 创建一个工作量证明的挖矿对象
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  //初始化目标整数
	target.Lsh(target, uint(256-targetBits)) //数据转换
	pow := &ProofOfWork{b, target}
	return pow
}

// 准备数据进行挖矿运算
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash, //上一块哈希
			//pow.Block.Transactions,              //当前数据
			utils.IntToHex(pow.Block.Timestamp), //时间十六进制
			utils.IntToHex(int64(targetBits)),   //位数十六进制
			utils.IntToHex(int64(nonce)),        //保存工作量的证明
		}, []byte{},
	)
	return data
}

// 挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	//fmt.Printf("当前挖矿计算的区块数据：%s", pow.Block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //获取要对比的数据
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	//fmt.Println("\n\n")
	return nonce, hash[:]

}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid

}
