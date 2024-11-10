package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rishavmehra/blockchain/block"
	"github.com/rishavmehra/blockchain/blockchain"
	"github.com/rishavmehra/blockchain/transaction"
)

type CLI struct {
	BC blockchain.Blockchain
}

func (cli *CLI) createblockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.DB.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain(address)
	defer bc.DB.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s' : %d\n", address, balance)
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println("getbalance -adddress ADDRESS - Get balance of ADDRESS")
	fmt.Println("createblockchain -adddress ADDRESS - create blockchain and send gensis block reward to Address")
	fmt.Println("printchain - Print all the blocks of the blockchain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {

	bc := blockchain.NewBlockchain("")
	defer bc.DB.Close()

	bci := bc.Iterator()
	for {
		blk := bci.Next()

		fmt.Printf("Prev. Hash %x\n", blk.PrevBlockHash)
		fmt.Printf("Hash %x\n", blk.Hash)
		pow := block.NewProofOfWork(blk)
		fmt.Printf("PoW %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(blk.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain(from)
	defer bc.DB.Close()

	tx := blockchain.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Println("success")
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	CreateBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Specify the wallet address to check its balance.")
	CreateBlockchainAddress := CreateBlockchainCmd.String("address", "", "Specify the wallet address to receive the genesis block reward.")
	sendFrom := sendCmd.String("from", "", "Enter the sender's wallet address.")
	sendTo := sendCmd.String("to", "", "Enter the recipient's wallet address.")
	sendAmount := sendCmd.Int("amount", 0, "Specify the amount to transfer.")

	switch os.Args[1] {

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "Createblockchain":
		err := CreateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if CreateBlockchainCmd.Parsed() {
		if *CreateBlockchainAddress == "" {
			CreateBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createblockchain(*CreateBlockchainAddress)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}
