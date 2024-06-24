package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"os"

	"github.com/TxCorpi0x/file-upload-merkle/merkle"
	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
	"github.com/TxCorpi0x/file-upload-merkle/types"
)

type HttpUploader struct {
	client  *http.Client
	baseURL string
}

func NewHttpUploader(httpClient *http.Client, baseURL string) *HttpUploader {
	return &HttpUploader{
		client:  httpClient,
		baseURL: baseURL,
	}
}

func (h *HttpUploader) UploadFilesFrom(filePaths []string) (
	uploadedFiles []types.UploadedFile,
	merkleRoot string,
	err error,
) {
	requestBody, formDataContentType, err := multipartFormFromFiles(filePaths)
	if err != nil {
		err = fmt.Errorf("%w: error preparing POST request body: %s", errFailedUpload, err)

		return
	}

	response, err := http.Post(fmt.Sprintf("%s/upload", h.baseURL), formDataContentType, &requestBody)
	if err != nil {
		err = fmt.Errorf("%w: error sending POST request: %s", errFailedUpload, err)

		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%w: unexpected http status: %s", errFailedUpload, response.Status)

		return
	}

	var decodedResponse types.UploadedFilesResponse
	if err = json.NewDecoder(response.Body).Decode(&decodedResponse); err != nil {
		err = fmt.Errorf("%w: error decoding json response: %s", errFailedUpload, err)

		return
	}

	defer func() { _ = response.Body.Close() }()

	merkleRoot, err = h.computeMerkleRoot(filePaths)
	if err != nil {
		err = fmt.Errorf("%w: error computing merkle root: %s", errFailedUpload, err)

		return
	}

	return decodedResponse.UploadedFiles, merkleRoot, nil
}

func (h *HttpUploader) computeMerkleRoot(filePaths []string) (merkleRoot string, err error) {
	var blocks merkle.Input
	for _, f := range filePaths {
		var fileContent []byte
		fileContent, err = os.ReadFile(f)
		if err != nil {
			err = fmt.Errorf("%w: error reading file for hashing: %s", errFailedUpload, err)

			return
		}

		blocks = append(blocks, fileContent)
	}

	merkleTree, err := merkle.NewTree(blocks, hash.NewSha256())
	if err != nil {
		return
	}

	return merkleTree.RootHex(), nil
}

func multipartFormFromFiles(filePaths []string) (multipartForm bytes.Buffer, formDataContentType string, err error) {
	multipartWriter := multipart.NewWriter(&multipartForm)

	for _, fp := range filePaths {
		var file *os.File
		file, err = os.Open(fp)
		if err != nil {
			return
		}

		var filePart io.Writer
		filePart, err = multipartWriter.CreateFormFile("files", filepath.Base(fp))
		if err != nil {
			return
		}

		// copy the file content to the form file part
		if _, err = io.Copy(filePart, file); err != nil {
			return
		}

		if err = file.Close(); err != nil {
			return
		}
	}

	// Close the multipart writer to finish building the request body
	if err = multipartWriter.Close(); err != nil {
		return
	}

	formDataContentType = multipartWriter.FormDataContentType()

	return
}
