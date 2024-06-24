package merkle

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"

	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
)

type Input [][]byte

// Tree is the type for merkle tree.
type Tree struct {
	// is the implemented Hasher interface for the desired hashing algorithm (e.g. Sha256).
	hasher hash.Hasher
	// Input is the source data to create the merkle tree.
	Input Input `json:"input"`
	// Nodes carries leaves and branches.
	Nodes hash.HashList `json:"nodes"`
}

// NewTree creates a new merkle tree using the provided information.
func NewTree(data Input, hasher hash.Hasher) (*Tree, error) {
	branchesLen := int(math.Exp2(math.Ceil(math.Log2(float64(len(data))))))

	// We pad our data length up to the power of 2.
	nodes := make(hash.HashList, branchesLen*2)

	tree := &Tree{
		hasher: hasher,
		Input:  data,
	}

	// We put the leaves after the branches in the slice of nodes.
	tree.fillHashes(nodes[branchesLen : branchesLen+len(tree.Input)])

	// Pad the space left after the leaves.
	for i := len(tree.Input) + branchesLen; i < len(nodes); i++ {
		nodes[i] = make([]byte, hasher.Len())
	}

	// Branches.
	tree.fillBranchHashes(nodes, branchesLen)

	tree.Nodes = nodes

	return tree, nil
}

// Proof generates proof for the node with the input content.
func (t *Tree) Proof(data []byte) (*Proof, error) {
	// Find the idx of the data
	idx, err := t.indexOf(data)
	if err != nil {
		return nil, err
	}

	return t.proofByIndex(idx)
}

// returns proof of node by input index.
func (t *Tree) proofByIndex(idx uint64) (*Proof, error) {
	if uint64(len(t.Input)) <= idx {
		return nil, errors.New("index out of range")
	}

	// calculate the minimum needed proof hashes to be able to check if a
	// hash belongs to a certain merkle tree: log(n)
	proofLen := int(math.Ceil(math.Log2(float64(len(t.Input)))))
	hashes := make(hash.HashList, proofLen)

	cur := 0
	minI := uint64(math.Pow(2, float64(1))) - 1
	for i := idx + uint64(len(t.Nodes)/2); minI < i; i /= 2 {
		hashes[cur] = t.Nodes[i^1]
		cur++
	}

	return newProof(hashes, idx), nil
}

// finds the index of the data to be proven in the merkle tree.
func (t *Tree) indexOf(input []byte) (uint64, error) {
	for i, data := range t.Input {
		if bytes.Equal(data, input) {
			return uint64(i), nil
		}
	}

	return 0, errors.New("input content was not found in the merkle")
}

// fills the hashes list of the tree pointer.
func (t *Tree) fillHashes(hashes hash.HashList) {
	for i := range t.Input {
		hashes[i] = t.hasher.Hash(t.Input[i])
	}
}

// fills branches with the corresponding hashes.
func (t *Tree) fillBranchHashes(nodes hash.HashList, leafOffset int) {
	for leafIdx := leafOffset - 1; leafIdx > 0; leafIdx-- {
		left := nodes[leafIdx*2]
		right := nodes[leafIdx*2+1]

		nodes[leafIdx] = t.hasher.Hash(left, right)
	}
}

// Root returns merkle root hash.
func (t *Tree) Root() []byte {
	return t.Nodes[1]
}

// RootHex returns calculated hexadecimal value of the binary root.
func (t *Tree) RootHex() string {
	return hex.EncodeToString(t.Root())
}
