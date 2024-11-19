package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	base58 "github.com/rishavmehra/blockchain/Base58"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
	// ScriptPubKey string // public key(locking script) - this locking utxo
}

type TxOutputs struct {
	Outputs []TxOutput
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

func (outs TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	fmt.Println("here is error")
	if err != nil {
		log.Panic(err)
	}
	return outputs
}
