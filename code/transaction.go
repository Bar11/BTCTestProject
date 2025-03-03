package code

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// 挖矿奖励
const subsidy = 100

// 交易：编号，输入，输出
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// 序列化
func (tx *Transaction) Serialize() []byte {
	var encoder bytes.Buffer
	enc := gob.NewEncoder(&encoder) // 编码器
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoder.Bytes()
}

// 反序列化
func (tx *Transaction) DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data)) // 解码器
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}
	return transaction
}

// 对于交易事务进行哈希
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize()) //取得二进制进行哈希计算
	return hash[:]
}

// 签名
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return //挖矿的无需签名
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("error pre transaction")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil

	}

}

// 用于签名的交易事务，裁剪的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, vin.Signature, vin.PubKey})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
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

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Transaction %x\n", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("input_id %d", i))
		lines = append(lines, fmt.Sprintf("TXID %x", input.Txid))
		lines = append(lines, fmt.Sprintf("output %d", input.Vout))
		lines = append(lines, fmt.Sprintf("Signature %d", input.Signature))
		lines = append(lines, fmt.Sprintf("PubKey %x", input.PubKey))
	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("output_id %d", i))
		lines = append(lines, fmt.Sprintf("Value %d", output.Value))
		lines = append(lines, fmt.Sprintf("pubKeyHash %x", output.PubKeyHash))

	}
	return strings.Join(lines, "\n")
}

// 签名认证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("error pre transaction")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256() //加密
	for inID, vin := range tx.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash //设置公钥

		r := big.Int{}
		s := big.Int{}
		signLen := len(vin.Signature) //统计签名长度
		r.SetBytes(vin.Signature[:signLen/2])
		s.SetBytes(vin.Signature[signLen/2:])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:keyLen/2])
		y.SetBytes(vin.PubKey[keyLen/2:])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}
	return true
}

//// 是否可以解锁输出
//func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
//	return out.ScriptPubKey == unlockingData
//}
//
//func (input *TXInput) CanUnlockOutputWith(unlockingData string) bool {
//	return input.ScriptSig == unlockingData
//}

// 挖矿交易
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("forward to %s", to)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{
		nil,
		[]TXInput{txin},
		[]TXOutput{*txout},
	}
	return &tx
}

// 转账交易
func NewUTXTransaction(from, to string, amount int, bc *Blockchain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallte(from)
	pubkeyHash := HashPubKey(wallet.PublicKey)

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
			input := TXInput{txID, out, nil, pubkeyHash}
			// 输出的交易
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, to))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}
