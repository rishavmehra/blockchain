package blockchain

import (
	"github.com/rishavmehra/blockchain/block"
)

type Blockchain struct {
	Blocks []*block.Block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := block.NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGensisBlock() *block.Block {
	return block.NewBlock("Gensis Block", []byte{0})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{NewGensisBlock()}}
}
