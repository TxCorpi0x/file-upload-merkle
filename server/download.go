package server

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/TxCorpi0x/file-upload-merkle/storage"
	"github.com/TxCorpi0x/file-upload-merkle/types"
)

func NewDownloadHandler(repository storage.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		index, err := indexFromRequest(r)
		if err != nil {
			httpError(w, http.StatusBadRequest, err)

			return
		}

		fileContent, err := repository.RetrieveFileByIndex(r.Context(), index)
		if err == storage.ErrStoredFileNotFound {
			httpError(w, http.StatusNotFound, fmt.Errorf("{index} not found: %d", index))

			return
		}
		if err != nil {
			httpError(w, http.StatusInternalServerError, err)

			return
		}

		_, err = w.Write(fileContent.Content)

		return
	}
}

func NewProofHandler(repository storage.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		index, err := indexFromRequest(r)
		if err != nil {
			httpError(w, http.StatusBadRequest, err)

			return
		}

		merkleTree, err := repository.RetrieveTree(r.Context())
		if err != nil {
			httpError(w, http.StatusInternalServerError, err)

			return
		}

		fileByIndex, err := repository.RetrieveFileByIndex(r.Context(), index)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, storage.ErrStoredFileNotFound) {
				statusCode = http.StatusNotFound
			}

			httpError(w, statusCode, err)

			return
		}

		merkleProof, err := merkleTree.Proof(fileByIndex.Content)
		if err != nil {
			httpError(w, http.StatusInternalServerError, err)

			return
		}

		if err = httpOkJson(w, types.MerkleProofResponse{MerkleProof: *merkleProof}); err != nil {
			httpError(w, http.StatusInternalServerError, err)
		}

		return
	}
}

func indexFromRequest(r *http.Request) (index int, err error) {
	vars := mux.Vars(r)
	indexParam, isIndexSet := vars["index"]
	if !isIndexSet {
		err = errors.New("{index} path param is not passed in")

		return
	}

	index, err = strconv.Atoi(indexParam)
	if err != nil {
		err = fmt.Errorf("{index} path param must be numeric: %s", err)
	}

	return
}
