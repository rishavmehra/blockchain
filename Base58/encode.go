package base58

import (
	"math/big"

	"github.com/rishavmehra/blockchain/utils"
)

const base58Alphabets = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(base58Alphabets)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabets[mod.Int64()])
	}

	utils.ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{base58Alphabets[0]}, result...)
		} else {
			break
		}
	}
	return result
}
