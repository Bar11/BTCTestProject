package cmd

import (
	"BTCTestProject/code"
	"fmt"
	"log"
)

func (cli *CLI) CreateBlockChain(address string) {
	if !code.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := code.CreateBlockChain(address)
	defer bc.DB.Close()
	fmt.Println("Blockchain created")
}
