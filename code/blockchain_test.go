package code

import (
	"fmt"
	"testing"
)

func TestCreateBlockChain(t *testing.T) {
	address := "17ZeVq9YmtLMgfrdcF7PyfXETaMBxA6JnE"
	bc := CreateBlockChain(address)
	fmt.Println(bc)

}

func TestBlockchain_FindUTXO(t *testing.T) {
	address := "17ZeVq9YmtLMgfrdcF7PyfXETaMBxA6JnE"
	bc := NewBlockchain()
	defer bc.DB.Close()
	balance := 0
	pubkeyhash_ := Base58Decode([]byte(address)) //提取公钥
	pubkeyhash := pubkeyhash_[1 : len(pubkeyhash_)-1]
	UTXO := bc.FindUTXO(pubkeyhash) //查找交易金额
	fmt.Println(bc.Tip)
	for _, out := range UTXO {
		fmt.Println(out.Value, "&&&&&&")
		balance += out.Value //去除金额
	}
	fmt.Printf("Balance:%d\n", balance)
}
