package wallets

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/rishavmehra/blockchain/wallet"
)

type Wallets struct {
	Wallets map[string]*wallet.Wallet
}

type SerializableWallet struct {
	D         *big.Int
	X, Y      *big.Int
	PublicKey []byte
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*wallet.Wallet)
	err := wallets.LoadFromFile()

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := wallet.NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet
	return address

}

func (ws *Wallets) GetWallet(address string) wallet.Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses

}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(wallet.WriteFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(wallet.WriteFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets map[string]SerializableWallet
	gob.Register(SerializableWallet{})
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = make(map[string]*wallet.Wallet)
	for k, v := range wallets {
		ws.Wallets[k] = &wallet.Wallet{
			PrivateKey: ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.X,
					Y:     v.Y,
				},
				D: v.D,
			},
			PublicKey: v.PublicKey,
		}
	}

	return nil

}

func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(SerializableWallet{})
	wallets := make(map[string]SerializableWallet)
	for k, v := range ws.Wallets {
		wallets[k] = SerializableWallet{
			D:         v.PrivateKey.D,
			X:         v.PrivateKey.PublicKey.X,
			Y:         v.PrivateKey.PublicKey.Y,
			PublicKey: v.PublicKey,
		}
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(wallet.WriteFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}
