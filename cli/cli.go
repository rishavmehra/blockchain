package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	base58 "github.com/rishavmehra/blockchain/Base58"
	"github.com/rishavmehra/blockchain/blockchain"
	"github.com/rishavmehra/blockchain/wallet"
	"github.com/rishavmehra/blockchain/wallets"
)

type CLI struct {
	BC blockchain.Blockchain
}

func (cli *CLI) createblockchain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("error: address is not valid")
	}
	bc := blockchain.CreateBlockchain(address)
	defer bc.DB.Close()
	UTXOSet := blockchain.UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}

func (clid *CLI) createWallet() {
	wallets, _ := wallets.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("error: address is not valid")
	}

	bc := blockchain.NewBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	defer bc.DB.Close()

	balance := 0
	pubKeyHash := base58.Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s' : %d\n", address, balance)
}

func (cli *CLI) listAddresses() {
	wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("error: sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("error: receiver address is not valid")
	}

	bc := blockchain.NewBlockchain()
	UTXOSet := blockchain.UTXOSet{bc}
	defer bc.DB.Close()

	tx := blockchain.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := blockchain.NewCoinbaseTx(from, "")
	txs := []*blockchain.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)

	fmt.Println("Success!")
}

func (cli *CLI) printChain() {

	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	bci := bc.Iterator()
	for {
		blk := bci.Next()

		fmt.Printf("========= Block %x =========\n", blk.Hash)
		fmt.Printf("Prev. Hash %x\n", blk.PrevBlockHash)
		fmt.Printf("Hash %x\n", blk.Hash)
		pow := blockchain.NewProofOfWork(blk)
		fmt.Printf("PoW %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range blk.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(blk.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("  createblockchain -address ADDRESS - create blockchain and send gensis block reward to Address")
	fmt.Println("  createwallet - genenrate new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses -Lists all the addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	CreateBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	CreateWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
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

	case "createblockchain":
		err := CreateBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := CreateWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressCmd.Parse(os.Args[2:])
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
	if CreateWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressCmd.Parsed() {
		cli.listAddresses()
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
