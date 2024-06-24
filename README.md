# github.com/TxCorpi0x/file-upload-merkle

## Merkle tree Implementation

`merkle` package contains a simple merkle tree implementation for single proof verification.

## Drawbacks

Addition to the [Limitations](https://github.com/fabiobozzo/merkle-file-uploader?tab=readme-ov-file#limitations-and-future-improvements), the following items can be considered.

### Design

- In Merkle tree we need to rehash the entire merkle tree to calculate root hash, alternative [zero-merkle-tree](https://github.com/cf/zero-merkle-tree-tutorial) can be implemented in Go.
- Decrease the storage allocation by using (Verkle-Tree)[https://github.com/ethereum/go-verkle].
- Usage of [MMR](https://github.com/ComposableFi/go-merkle-trees/tree/main/mmr) instead of merkle tree to minimize the storage.
- Use [libp2p](https://github.com/libp2p/go-libp2p) to store the files in a decentralized manner.
- Integrate IPFS to ensure decentralized file system and proof verification.

### Implementation

- The implemented Merkle tree does not support multi-proof, it can be implemented to support multiple proofs verification.

## Resources

Modified version of [merkle-file-uploader](https://github.com/fabiobozzo/merkle-file-uploader).
