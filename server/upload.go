package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TxCorpi0x/file-upload-merkle/merkle"
	"github.com/TxCorpi0x/file-upload-merkle/merkle/hash"
	"github.com/TxCorpi0x/file-upload-merkle/storage"
	"github.com/TxCorpi0x/file-upload-merkle/types"
)

func NewUploadHandler(repository storage.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		// limit maxMultipartMemory
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			httpError(w, http.StatusBadRequest, fmt.Errorf("unable to parse multipart form: %s", err))

			return
		}

		if err := repository.DeleteAllFiles(r.Context()); err != nil {
			httpError(w, http.StatusInternalServerError, fmt.Errorf("error while resetting storage: %s", err))

			return
		}

		var uploadedFiles []types.UploadedFile
		var blocks [][]byte

		files := r.MultipartForm.File["files"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				httpError(w, http.StatusBadRequest, fmt.Errorf("unable to open file: %s", err))

				return
			}
			_ = file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				httpError(w, http.StatusBadRequest, fmt.Errorf("unable to read file: %s", err))

				return
			}

			i, err := repository.StoreFile(r.Context(), storage.StoredFile{
				Name:    fileHeader.Filename,
				Content: data,
			})
			if err != nil {
				httpError(w, http.StatusInternalServerError, err)

				return
			}

			uploadedFiles = append(uploadedFiles, types.UploadedFile{
				Name:  fileHeader.Filename,
				Index: i,
			})

			blocks = append(blocks, data)
		}

		merkleTree, err := merkle.NewTree(blocks, hash.NewSha256())
		if err != nil {
			httpError(w, http.StatusInternalServerError, err)

			return
		}

		if err = repository.StoreTree(r.Context(), merkleTree); err != nil {
			httpError(w, http.StatusInternalServerError, fmt.Errorf("unable to store the merkle tree: %s", err))

			return
		}

		if err := httpOkJson(w, types.UploadedFilesResponse{UploadedFiles: uploadedFiles}); err != nil {
			httpError(w, http.StatusInternalServerError, err)

			return
		}
	}
}
