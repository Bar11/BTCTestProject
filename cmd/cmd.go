package cmd

import (
	"BTCTestProject/code"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	Blockchain *code.Blockchain
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t addblock 向区块链增加快")
	fmt.Println("\t showchain 显示区块链")
	fmt.Println("\t addblock 向区块链增加快")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(transactions []*code.Transaction) {
	cli.Blockchain.AddBlock(transactions)
	fmt.Println("Added block successfully!")
}

func (cli *CLI) showBlockChain() {
	bci := cli.Blockchain.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("timestamp:%d\n", block.Timestamp)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("data:%s\n", block.Transactions[0].ID)
		fmt.Printf("Hash:%x\n", block.Hash)
		pow := code.NewProofOfWork(block)
		fmt.Printf("pow %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			//遇到创世区块
			break
		}
	}

}

func (cli *CLI) Run() {
	cli.validateArgs() // 校验
	addblockcmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	showChaincmd := flag.NewFlagSet("showchain", flag.ExitOnError)

	addBlockData := addblockcmd.String("data", "", "Block data")
	switch os.Args[1] {
	case "addblock":
		err := addblockcmd.Parse(os.Args[2:]) //解析参数
		if err != nil {
			log.Panic(err)
		}
	case "showchain":
		err := showChaincmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)

	}
	if addblockcmd.Parsed() {
		if *addBlockData == "" {

			addblockcmd.Usage()
			os.Exit(1)
		} else {
			cli.addBlock(*addBlockData) //增加区块
		}
	}
	if showChaincmd.Parsed() {
		cli.showBlockChain() // 显示区块链
	}
}
