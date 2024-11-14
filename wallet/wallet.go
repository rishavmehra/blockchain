package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	base58 "github.com/rishavmehra/blockchain/Base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	Version            = byte(0x00)
	WriteFile          = "wallet.dat"
	AddressCheckSumLen = 4
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := NewKeypair()
	wallet := Wallet{private, public}

	return &wallet
}

func NewKeypair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubkey
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checkSum := CheckSum(versionedPayload)

	fullPayload := append(versionedPayload, checkSum...)
	address := base58.Encode(fullPayload)
	return address
}

func ValidateAddress(address string) bool {
	pubKeyHash := base58.Decode([]byte(address))
	actualCheckSum := pubKeyHash[len(pubKeyHash)-AddressCheckSumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressCheckSumLen]
	targetChecksum := CheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualCheckSum, targetChecksum) == 0
}

func HashPubKey(pubkey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubkey)

	// need to improve the ripemd
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func CheckSum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	SecondSHA := sha256.Sum256(firstSHA[:])

	return SecondSHA[:AddressCheckSumLen]
}
