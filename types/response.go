package types

import "github.com/TxCorpi0x/file-upload-merkle/merkle"

// UploadedFile struct used for file names and index.
type UploadedFile struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

// UploadedFilesResponse is the http response for file uploader server endpoint.
type UploadedFilesResponse struct {
	UploadedFiles []UploadedFile `json:"uploadedFiles"`
}

// MerkleProofResponse is the http response of downloader server endpoint to get the proof of downloaded file.
type MerkleProofResponse struct {
	MerkleProof merkle.Proof `json:"merkleProof"`
}
