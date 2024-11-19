package blockchain

import (
	"bytes"

	"github.com/rishavmehra/blockchain/wallet"
)

type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
	// ScriptSig string //valid digital signature made using the private key // this provide the data to unlock a utxo for spending(unlocking script)
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
