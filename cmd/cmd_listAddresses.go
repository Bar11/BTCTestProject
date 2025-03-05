package cmd

import (
	"BTCTestProject/code"
	"fmt"
)

func (cli *CLI) ListAddresses() {
	wallets, err := code.NewWallets()
	if err != nil {
		panic(err)
	}
	addresses := wallets.GetAddresses()
	for _, addr := range addresses {
		fmt.Println(addr)
	}
}
