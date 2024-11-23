package decentralized

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/rishavmehra/blockchain/blockchain"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string          // address to receive mining reward
var knowNodes = []string{":3000"} // this central node, every node must to know where to connect to initially

type Version struct {
	Ver        int
	BestHeight int    // store the length of the node's blockchain
	AddrFrom   string // address of the sender
}

func CommandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]

}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := blockchain.NewBlockchain(nodeID)
	if nodeAddress != knowNodes[0] {

	}
}

func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Println("%s is not available \n ", addr)
		var updatedNodes []string

	}
}

func SendVersion(addr string, bc *blockchain.Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddress})

	request := append(CommandToBytes("version"), payload...)

	send
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()

}
