package blockchain

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/rishavmehra/blockchain/block"
)

const dbfile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	// Blocks []*block.Block
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := block.NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialization())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func NewBlockchain() *Blockchain {
	// return &Blockchain{[]*block.Block{block.NewGenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := block.NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialization())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}

	return &bc
}
