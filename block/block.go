package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/rishavmehra/blockchain/transaction"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Serialization() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil
	}
	return &block
}

func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
