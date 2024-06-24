package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
	"github.com/TxCorpi0x/file-upload-merkle/types"
)

type HttpDownloader struct {
	client   *http.Client
	baseURL  string
	rootHash hash.Hash
}

func NewHttpDownloader(httpClient *http.Client, baseURL string, rootHash hash.Hash) *HttpDownloader {
	return &HttpDownloader{
		client:   httpClient,
		baseURL:  baseURL,
		rootHash: rootHash,
	}
}

func (h *HttpDownloader) DownloadFileAt(index int, destination *os.File) (err error) {
	downloadResponse, err := http.Get(fmt.Sprintf("%s/download/%d", h.baseURL, index))
	if err != nil {
		err = fmt.Errorf("%w: error sending GET /download request: %s", errFailedDownload, err)

		return
	}
	defer func() { _ = downloadResponse.Body.Close() }()

	if downloadResponse.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("%w: file not found at index %d", errFailedDownload, index)

		return
	}

	proofResponse, err := http.Get(fmt.Sprintf("%s/proof/%d", h.baseURL, index))
	if err != nil {
		err = fmt.Errorf("%w: error sending GET /proof request: %s", errFailedDownload, err)

		return
	}
	defer func() { _ = proofResponse.Body.Close() }()

	var merkleProof types.MerkleProofResponse
	if err = json.NewDecoder(proofResponse.Body).Decode(&merkleProof); err != nil {
		err = fmt.Errorf("%w: error decoding merkle proof response body: %s", errFailedDownload, err)

		return
	}

	fileContent, err := io.ReadAll(downloadResponse.Body)
	if err != nil {
		err = fmt.Errorf("%w: error reading download response body %s", errFailedDownload, err)

		return
	}

	verified, err := merkleProof.MerkleProof.Verify(fileContent, h.rootHash, hash.NewSha256())
	if err != nil {
		err = fmt.Errorf("%w: merkle root does not match: %s", errFailedProveHash, err)
	}
	if !verified {
		err = fmt.Errorf("%w: merkle root does not match: %s", errFailedDownload, h.rootHash)

		return
	}

	reader := io.NopCloser(bytes.NewReader(fileContent))
	if _, err = io.Copy(destination, reader); err != nil {
		err = fmt.Errorf("%w: error reading downloaded file: %s", errFailedDownload, err)
	}

	return
}
