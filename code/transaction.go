package code

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 挖矿奖励
const subsidy = 10

// 输入
type TXInput struct {
	Txid      []byte //c存储了交易的id
	Vout      int    // 保存该交易中的一个output索引
	ScriptSig string //保存用户的钱包地址
}

// 输出
type TXOutput struct {
	Value        int //保存了币
	ScriptPubKey string
}

// 交易：编号，输入，输出
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// 检查交易失误是否为coinbase,挖矿奖励
func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// 设置交易ID,从二进制数据中
func (tx *Transaction) SetID() {
	var encode bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encode)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encode.Bytes())
	tx.ID = hash[:]

}

// 是否可以解锁输出
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func (input *TXInput) CanUnlockOutPutWith(unlockingData string) bool {
	return input.ScriptSig == unlockingData
}

// 挖矿交易
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{
		nil,
		[]TXInput{txin},
		[]TXOutput{txout},
	}
	return &tx
}

// 转账交易
func NewUTXTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	acc, validOutputs := bc.FindSpendableOutPuts(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//遍历无效输出
	for txid, out := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range out {
			// 输入的交易
			input := TXInput{txID, out, from}
			// 输出的交易
			inputs = append(inputs, input)
		}
	}
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})

	}
	return &Transaction{}
}
