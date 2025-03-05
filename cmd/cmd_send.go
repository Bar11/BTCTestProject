package cmd

import (
	"BTCTestProject/code"
	"fmt"
	"log"
)

func (cli *CLI) Send(from, to string, amount int) {
	if !code.ValidateAddress(from) {
		log.Panic("Invalid from address")
	}
	if !code.ValidateAddress(to) {
		log.Panic("Invalid to address")
	}
	bc := code.NewBlockchain()
	defer bc.DB.Close()
	tx := code.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*code.Transaction{tx}) //挖矿记账
	fmt.Println("New UTX Transaction created")
}
