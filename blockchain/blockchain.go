package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/rishavmehra/blockchain/block"
	"github.com/rishavmehra/blockchain/transaction"
	"github.com/rishavmehra/blockchain/wallet"
	"github.com/rishavmehra/blockchain/wallets"
)

const dbfile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := block.NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialization())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) FindUnSpendTransactions(pubKeyHash []byte) []transaction.Transaction {
	var unspendTXs []transaction.Transaction
	spendTXOs := make(map[string][]int)
	//interate over all the blocks in the blockchain
	bci := bc.Iterator()

	for {
		// get the next block in the blockchain
		block := bci.Next()
		// iterate over all the transactions in the block
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs: // iterate over all the outputs in the transaction
			for outIdx, out := range tx.Vout {
				// check if the output is already spent
				if spendTXOs[txID] != nil {
					for _, spendOut := range spendTXOs[txID] {
						if spendOut == outIdx {
							continue Outputs
						}
					}
				}
				// check if the output can be unlocked with the address
				if out.IsLockedWithKey(pubKeyHash) {
					unspendTXs = append(unspendTXs, *tx)
				}
			}
			// if the transaction is not a coinbase transaction
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spendTXOs[inTxID] = append(spendTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		// if the block is the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspendTXs
}

// / method returns the unspent transaction outputs for a given address
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transaction.TxOutput {
	var UTXOs []transaction.TxOutput
	unspendTransactions := bc.FindUnSpendTransactions(pubKeyHash)

	// iterate over all the transactions in the block
	for _, tx := range unspendTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *transaction.Transaction {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	wallets, err := wallets.NewWallets()
	if err != nil {
		log.Panic(err)
	}
	Wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(Wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Bhai tare pass etna bitcoin he nahi h")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := transaction.TxInput{txID, out, nil, Wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *transaction.NewTxOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *transaction.NewTxOutput(acc-amount, from))
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, Wallet.PrivateKey)

	return &tx
}

func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspendOutputs := make(map[string][]int)
	unspendTXs := bc.FindUnSpendTransactions(pubKeyHash)
	allTxTOgether := 0

Work:
	for _, tx := range unspendTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIDx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && allTxTOgether < amount {
				allTxTOgether = allTxTOgether + out.Value
				unspendOutputs[txID] = append(unspendOutputs[txID], outIDx)

				if allTxTOgether >= amount {
					break Work
				}
			}
		}
	}

	return allTxTOgether, unspendOutputs

}

func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No exiting blockchain found, create one first. ")
		os.Exit(1)
	}

	var Tip []byte
	db, err := bolt.Open(dbfile, 0666, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		Tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{Tip, db}

	return &bc
}

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exits")
		os.Exit(1)
	}

	var Tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := transaction.NewCoinbaseTx(address, genesisCoinbaseData)
		genesis := block.NewGenesisBlock(cbtx)

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
		Tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{Tip, db}
	return &bc

}

func (bc *Blockchain) FindTransaction(ID []byte) (transaction.Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return transaction.Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	tx.Sign(privKey, prevTxs)
}

func (bc *Blockchain) VerifyTransaction(tx *transaction.Transaction) bool {
	prevTxs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	return tx.Verify(prevTxs)
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.DB}

	return bci
}

func (i *BlockchainIterator) Next() *block.Block {
	var blk *block.Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		blk = block.DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = blk.PrevBlockHash
	return blk
}
