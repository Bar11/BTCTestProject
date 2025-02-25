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

func (cli *CLI) createBlockChain(address string) {
	bc := code.CreateBlockChain(address)
	defer bc.DB.Close()
	fmt.Println("Blockchain created")
}

func (cli *CLI) getBalance(address string) {
	bc := code.NewBlockchain(address)
	defer bc.DB.Close()
	balance := 0
	UTXO := bc.FindUTXO(address) //查找交易金额
	for _, out := range UTXO {
		balance += out.Value //去除金额
	}
	fmt.Printf("Address:%s;Balance:%d\n", address, balance)

}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t getbalance -address 'get balance by address' ")
	fmt.Println("\t createblockchain -address 'create blockchain by address'")
	fmt.Println("\t send -from From -to To -amount Amount ")
	fmt.Println("\t showchain 'show chain' ")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := code.NewBlockchain(from)
	defer bc.DB.Close()
	tx := code.NewUTXTransaction(from, to, amount, bc)
	bc.MineBlock([]*code.Transaction{tx}) //挖矿记账
	fmt.Println("New UTX Transaction created")
}

func (cli *CLI) showBlockChain() {
	bc := code.NewBlockchain("")
	defer bc.DB.Close()

	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := code.NewProofOfWork(block) //工作量证明
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs() // 校验
	getbalance_cmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createbc_cmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendcmd := flag.NewFlagSet("send", flag.ExitOnError)
	showChaincmd := flag.NewFlagSet("showchain", flag.ExitOnError)

	getbc_param_address := getbalance_cmd.String("address", "", "get balance by address")
	createbc_param_address := createbc_cmd.String("address", "", "new bc address")
	send_param_amount := sendcmd.Int("amount", 0, "amount to send")
	send_param_from := sendcmd.String("from", "", "from address")
	send_param_to := sendcmd.String("to", "", "to address")

	switch os.Args[1] {
	case "getbalance":
		err := getbalance_cmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createbc_cmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendcmd.Parse(os.Args[2:])
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
	if getbalance_cmd.Parsed() {
		if *getbc_param_address == "" {
			getbalance_cmd.PrintDefaults()
			os.Exit(1)
		}
		cli.getBalance(*getbc_param_address)
	}
	if createbc_cmd.Parsed() {
		if *createbc_param_address == "" {
			createbc_cmd.PrintDefaults()
		}
		cli.createBlockChain(*createbc_param_address)
	}
	if sendcmd.Parsed() {
		if *send_param_amount == 0 || *send_param_from == "" || *send_param_to == "" {
			sendcmd.PrintDefaults()
			os.Exit(1)
		}
		cli.send(*send_param_from, *send_param_to, *send_param_amount)
	}
	if showChaincmd.Parsed() {
		cli.showBlockChain()
	}
}
