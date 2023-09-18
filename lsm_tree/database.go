package main

import (
	"errors"
	"fmt"
	"time"
)

type storage struct {
	// The group of indices sorted in an ascending order.
	indices []*Index

	// The path of the directory that files are saved.
	dirPath string

	// The maximum byte size of each index. If the index size is bigger than this, the storage makes a new index.
	maxIndexByteSize int64
}

type NotFoundError struct {
	message string
}

func NewStorage(indices []*Index, dirPath string) *storage {
	return &storage{
		indices: indices,
		dirPath: dirPath,
		maxIndexByteSize: 100,
	}
}

func (s *storage) Set(key string, value string) error {
	index := s.currentIndex()
	return index.Set(key, value)
}

func (s *storage) currentIndex() *Index {
	if len(s.indices) == 0 {
		newIndex := NewIndex(s.newFilePath(), make(map[string]int64), 0)
		s.indices = append(s.indices, newIndex)
		return newIndex
	}

	lastIndex := s.indices[len(s.indices)-1]
	if lastIndex.lastByteOffset > s.maxIndexByteSize{
		newIndex := NewIndex(s.newFilePath(), make(map[string]int64), 0)
		s.indices = append(s.indices, newIndex)
		return newIndex
	}

	return lastIndex
}

func (s *storage) newFilePath() string {
	return fmt.Sprintf("%s/%s", s.dirPath, time.Now().String())
}

func (s *storage) Get(key string) (string, error) {
	for i := 0; i < len(s.indices); i++{
		index := s.indices[len(s.indices) - 1 - i]
		value, err := index.Get(key)
		notFoundError := &NotFoundError{}
		if errors.As(err, &notFoundError){
			continue
		}
		if err != nil {
			return "", err
		}

		return value, nil
	}
	return "", newNotFoundError(key)
}

func (s *storage) PrintHashMap() {
	for _, index := range s.indices {
		fmt.Printf("INDEX NAME: %s | BYTE SIZE: %d\n", index.Name(), index.Size())
		index.PrintHashMap()
	}
}


func newNotFoundError(key string) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf("no value corresponding to the key: %s", key),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}
