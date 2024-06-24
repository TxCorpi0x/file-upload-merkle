package hash

// Hash is a type alias for bytes.
type Hash []byte

// HashList is a type alias fo array of bytes.
type HashList [][]byte

// Hasher is the interface which needs to be implemented by the desired hashing algorithm.
type Hasher interface {
	// Hash hashes the input bytes and returns hashed bytes.
	Hash(data ...[]byte) Hash
	// Len returns constant length of hashing algorithm.
	Len() int
}
