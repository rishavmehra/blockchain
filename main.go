package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// func main() {
// 	name := "100xdevs"
// 	result := sha256.Sum256([]byte(name))
// 	fmt.Printf("%x", result)
// }

func findHashwithPrefix(prefix string) (string, string) {
	input := 0
	for {
		inputstr := "100xdevs" + strconv.Itoa(input)
		hash := sha256.Sum256([]byte(inputstr))
		hashStr := hex.EncodeToString(hash[:])
		if len(hashStr) >= len(prefix) && hashStr[:len(prefix)] == prefix {
			return inputstr, hashStr
		}
		input++
	}
}
func main() {
	input, hash := findHashwithPrefix("00000")
	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Hash: %s\n", hash)
}
