package storage

import (
	"context"
	"errors"

	"github.com/TxCorpi0x/file-upload-merkle/merkle"
)

var (
	ErrStoredFileNotFound = errors.New("the file is not found in the storage")
)

type StoredFile struct {
	Index   int
	Name    string
	Content []byte
}

type Repository interface {
	StoreFile(context.Context, StoredFile) (int, error)
	RetrieveFileByIndex(context.Context, int) (StoredFile, error)
	DeleteAllFiles(context.Context) error
	StoreTree(context.Context, *merkle.Tree) error
	RetrieveTree(context.Context) (*merkle.Tree, error)
}
