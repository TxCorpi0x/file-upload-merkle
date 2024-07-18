package merkle

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"

	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
)

// type alias for input data, slice of bytes.
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
	tree := &Tree{hasher: hasher, Input: data}

	// calculate branches length of tree according to the input data
	branchesLen := tree.BranchesLen()

	// if we have x branches this means that we have double nodes.
	nodes := make(hash.HashList, branchesLen*2)

	// fill hashes extract bottom line leaves
	addedNodes := nodes[branchesLen : branchesLen+len(tree.Input)]
	tree.fillHashes(addedNodes)

	// allocate nodes hashes for leaves.
	for i := len(tree.Input) + branchesLen; i < len(nodes); i++ {
		nodes[i] = make([]byte, hasher.Len())
	}

	// fill branch hashes according to the left and right nodes
	tree.fillBranchHashes(nodes, branchesLen)

	tree.Nodes = nodes

	return tree, nil
}

// LevelsLen calculates the levels length of the tree according to the data length.
// number of levels of a merkle tree follow Log2(n) since the number of nodes doubles every level
// e.g 1M leaves Log2(1M) = 20
// e.g 2M leaves Log2(2M) = 30
func (t *Tree) LevelsLen() float64 {
	return math.Ceil(math.Log2(float64(len(t.Input))))
}

// BranchesLen calculates the total number of branches in the tree.
// calculate the branch length according to the levels
// this ensures equal left and right branches in a tree.
// e.g if there are 8 nodes, we have 3 levels and 8 branches.
// e.g if there are 9 nodes, we have 4 levels and 16 branches.
func (t *Tree) BranchesLen() int {
	return int(math.Exp2(t.LevelsLen()))
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
	// hash belongs to a certain merkle tree
	// we need at least one node from each level of the tree in order to build the root.
	proofLen := int(t.LevelsLen())
	hashes := make(hash.HashList, proofLen)

	cur := 0
	minI := uint64(1) // the root
	startNodeIndex := idx + uint64(len(t.Nodes)/2)
	// loop through each level and pick one node.
	for i := startNodeIndex; minI < i; i /= 2 {
		// get the even member of each branch and put into the proof hash,
		// this leads us to the root node eventually.
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
