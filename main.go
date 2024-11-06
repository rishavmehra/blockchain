package main

import (
	"github.com/rishavmehra/blockchain/blockchain"
	"github.com/rishavmehra/blockchain/cli"
)

func main() {
	// bc := blockchain.NewBlockchain()
	// bc.AddBlock("Send 1 BTC to Ivan")
	// bc.AddBlock("Send 2 more BTC to Ivan")

	// for _, block := range bc.Blocks {
	// 	fmt.Printf("Previous Block Hash: %x\n", block.PrevBlockHash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Nonce: %x\n", block.Nonce)
	// 	fmt.Printf("Timestamp: %x\n", block.Timestamp)
	// 	fmt.Printf("Block Hash: %x\n", block.Hash)
	// 	pow := blk.NewProofOfWork(block)
	// 	fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
	// 	fmt.Println()
	// }

	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := &cli.CLI{bc}
	cli.Run()

}
