# File Upload Merkle

## Usage Example

Build the binary

```bash
make build-client
```

Run a containerized server on port `8080`.

```bash
make start-server # or `docker compose up` to attach to process. or ./fxmerkle server
```

Sample uploads: The following command, uploads the files under `.runtime/files` to the server and puts the corresponding hash to the `.runtime/merkleroot` file for further processing and file verification.

```bash
make test-upload # ./fxmerkle client upload .runtime/files
```

Download the file at index and verify the proof received from server to the file content.

```bash
make test-download # ./fxmerkle client download 1
```

Stop containerized server

```bash
make stop-server # or `docker compose down` to attach to process. 
```

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
- Usage of [DAG](https://github.com/heimdalr/dag) instead of merkle tree(the same as IPFS).

### Implementation

- The implemented Merkle tree does not support multi-proof, it can be implemented to support multiple proofs verification.
- Support multi-chunk file upload to support large files.
- Support insertion and deletion using [bm](https://github.com/sorpaas/bm) in-place tree modification.
- Use persistent data store instead of memory storage in the server.

## Resources

Modified version of [merkle-file-uploader](https://github.com/fabiobozzo/merkle-file-uploader).
