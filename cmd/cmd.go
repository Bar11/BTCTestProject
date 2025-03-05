package cmd

import (
	"BTCTestProject/code"
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	Blockchain *code.Blockchain
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t createwallet 'create a new wallet' ")
	fmt.Println("\t listaddresses 'list all address' ")
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

func (cli *CLI) Run() {
	cli.validateArgs() // 校验
	getbalance_cmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	listaddresses_cmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	createwallet_cmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
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
	case "listaddresses":
		err := listaddresses_cmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createwallet_cmd.Parse(os.Args[2:])
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
		cli.GetBalance(*getbc_param_address)
	}
	if listaddresses_cmd.Parsed() {
		cli.ListAddresses()
	}
	if createwallet_cmd.Parsed() {
		cli.CreateWallet()
	}
	if showChaincmd.Parsed() {
		cli.ShowBlockChain()
	}
	if createbc_cmd.Parsed() {
		if *createbc_param_address == "" {
			createbc_cmd.PrintDefaults()
		}
		fmt.Println(*createbc_param_address)
		cli.CreateBlockChain(*createbc_param_address)
	}
	if sendcmd.Parsed() {
		if *send_param_amount == 0 || *send_param_from == "" || *send_param_to == "" {
			sendcmd.PrintDefaults()
			os.Exit(1)
		}
		cli.Send(*send_param_from, *send_param_to, *send_param_amount)
	}

}
