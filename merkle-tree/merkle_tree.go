package merkletree

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Data  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// if not odd, add dulicate the last item
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// create lead nodes
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	// reaching to the root node
	for i := 0; i < len(data); i++ {
		var newLevel []MerkleNode

		// make to child node for the specific node
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel

	}

	mtree := MerkleTree{&nodes[0]}
	return &mtree
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Right = right
	mNode.Left = left
	return &mNode
}
