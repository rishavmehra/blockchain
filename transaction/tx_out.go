package transaction

import (
	"bytes"

	base58 "github.com/rishavmehra/blockchain/Base58"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
	// ScriptPubKey string // public key(locking script) - this locking utxo
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
