package main

import (
	"fmt"

	"github.com/rishavmehra/blockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	bc.AddBlock("Send 100 BC to Rishav")
	bc.AddBlock("Send 10Bc to World")

	for _, block := range bc.Blocks {
		fmt.Printf("Previous Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		fmt.Println()
	}
}

// func main() {
// 	name := "100xdevs"
// 	result := sha256.Sum256([]byte(name))
// 	fmt.Printf("%x", result)
// }

//	func findHashwithPrefix(prefix string) (string, string) {
//		input := 0
//		for {
//			inputstr := "100xdevs" + strconv.Itoa(input)
//			hash := sha256.Sum256([]byte(inputstr))
//			hashStr := hex.EncodeToString(hash[:])
//			if len(hashStr) >= len(prefix) && hashStr[:len(prefix)] == prefix {
//				return inputstr, hashStr
//			}
//			input++
//		}
//	}
