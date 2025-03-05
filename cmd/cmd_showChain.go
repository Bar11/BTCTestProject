package cmd

import (
	"BTCTestProject/code"
	"fmt"
	"strconv"
)

func (cli *CLI) ShowBlockChain() {
	bc := code.NewBlockchain()
	defer bc.DB.Close()

	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := code.NewProofOfWork(block) //工作量证明
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println("\n\n")
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
