package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rishavmehra/blockchain/block"
	"github.com/rishavmehra/blockchain/blockchain"
)

type CLI struct {
	Bc *blockchain.Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println("addblock - data BLOCK_DATA - add a block to the chain")
	fmt.Println("printchain - print all the blocks of the blockchain")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.Bc.AddBlock(data)
	fmt.Println("Success!!")
}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()

	for {
		blk := bci.Next()

		fmt.Printf("Prev. Hash %x\n", blk.PrevBlockHash)
		fmt.Printf("Data %s\n", blk.Data)
		fmt.Printf("Hash %x\n", blk.Hash)
		pow := block.NewProofOfWork(blk)
		fmt.Printf("PoW %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(blk.PrevBlockHash) == 0 {
			break
		}

	}

}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
