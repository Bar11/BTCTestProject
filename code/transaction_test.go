package code

import (
	"fmt"
	"testing"
)

func TestNewCoinBaseTX(t *testing.T) {
	address := "17ZeVq9YmtLMgfrdcF7PyfXETaMBxA6JnE"
	cbtx := NewCoinBaseTX(address, genesisCoinbaseData)
	fmt.Println(cbtx.ID)
	fmt.Println(cbtx.Vin)
	fmt.Println(cbtx.Vout)
	fmt.Println("Txid:", cbtx.Vin[0].Txid)
	fmt.Println("PubKey:", cbtx.Vin[0].PubKey)
	fmt.Println("Vout:", cbtx.Vin[0].Vout)
	fmt.Println("Signature:", cbtx.Vin[0].Signature)
	fmt.Println("Value:", cbtx.Vout[0].Value)
	fmt.Println("PubKeyHash:", cbtx.Vout[0].PubKeyHash)
}
