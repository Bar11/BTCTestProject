package cmd

import (
	"BTCTestProject/code"
	"fmt"
)

func (cli *CLI) CreateWallet() {
	wallets, _ := code.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Println("you wallet address: " + address)

}
