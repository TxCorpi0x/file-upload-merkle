package merkle

import (
	"bytes"

	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
)

// Proof is a proof of a Merkle tree.
type Proof struct {
	Hashes hash.HashList `json:"hashes"`
	Index  uint64        `json:"index"`
}

// initializes a new proof object out of hashes and index.
func newProof(hashes hash.HashList, idx uint64) *Proof { return &Proof{Hashes: hashes, Index: idx} }

// Verify if the root hash bytes is equal to generated hash of proof.
func (p *Proof) Verify(data []byte, rootHash hash.Hash, hasher hash.Hasher) (bool, error) {
	proofHash := p.Hash(data, hasher)

	if bytes.Equal(rootHash, proofHash) {
		return true, nil
	}

	return false, nil
}

// Hash returns the proof hash.
func (p *Proof) Hash(data []byte, hasher hash.Hasher) []byte {
	var proofHash []byte

	proofHash = hasher.Hash(data)
	idx := p.Index + (1 << uint(len(p.Hashes)))

	// calculate hashes for the proof
	for _, hash := range p.Hashes {
		if idx%2 == 0 {
			// hash is on the left hand side of tree
			proofHash = hasher.Hash(proofHash, hash)
		} else {
			// hash is on the right hand side of tree
			proofHash = hasher.Hash(hash, proofHash)
		}
		// shift index right by one
		idx >>= 1
	}

	return proofHash
}
