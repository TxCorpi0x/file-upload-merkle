package hash

import (
	"golang.org/x/crypto/sha3"
)

const sha256Len = 32

// Sha256 is the 256-bit SHA3 hashing method.
type Sha256 struct{}

// New256 creates a new 256-bit SHA3 hashing method.
func NewSha256() *Sha256 {
	return &Sha256{}
}

// Hash generates a SHA3 hash from input byte arrays.
func (algo *Sha256) Hash(data ...[]byte) Hash {
	var hash [sha256Len]byte
	if len(data) == 1 {
		hash = sha3.Sum256(data[0])
	} else {
		concatDataLen := 0
		for _, d := range data {
			concatDataLen += len(d)
		}
		concatData := make([]byte, concatDataLen)
		curOffset := 0
		for _, d := range data {
			copy(concatData[curOffset:], d)
			curOffset += len(d)
		}
		hash = sha3.Sum256(concatData)
	}

	return hash[:]
}

// Len returns constant length of the hashing algorithm.
func (*Sha256) Len() int {
	return sha256Len
}
