package cmd

import (
	"BTCTestProject/code"
	"fmt"
	"log"
)

func (cli *CLI) GetBalance(address string) {
	if !code.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := code.NewBlockchain()
	defer bc.DB.Close()
	balance := 0
	pubkeyhash_ := code.Base58Decode([]byte(address)) //提取公钥
	pubkeyhash := pubkeyhash_[1 : len(pubkeyhash_)-1]
	UTXO := bc.FindUTXO(pubkeyhash) //查找交易金额
	fmt.Println("UTXO:", UTXO)
	for _, out := range UTXO {
		fmt.Println(out.Value, "&&&&&&")
		balance += out.Value //去除金额
	}
	fmt.Printf("Balance:%d\n", balance)

}
