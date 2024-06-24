package storage

import (
	"context"
	"github.com/TxCorpi0x/file-upload-merkle/merkle"
	"sync"
)

var _ Repository = (*InMemoryStorage)(nil)

type InMemoryStorage struct {
	mu    sync.RWMutex
	seq   int
	files map[int]StoredFile
	tree  *merkle.Tree
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		files: make(map[int]StoredFile),
	}
}

func (s *InMemoryStorage) StoreFile(_ context.Context, file StoredFile) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	file.Index = s.seq
	s.files[s.seq] = file

	return s.seq, nil
}

func (s *InMemoryStorage) RetrieveFileByIndex(_ context.Context, i int) (storedFile StoredFile, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storedFile, found := s.files[i]
	if !found {
		err = ErrStoredFileNotFound
	}

	return
}

func (s *InMemoryStorage) DeleteAllFiles(_ context.Context) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq = 0
	s.files = make(map[int]StoredFile)

	return nil
}

func (s *InMemoryStorage) StoreTree(_ context.Context, tree *merkle.Tree) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tree = tree
	return nil
}

func (s *InMemoryStorage) RetrieveTree(_ context.Context) (*merkle.Tree, error) {
	return s.tree, nil
}
