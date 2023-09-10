package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type storage struct {
	// TODO: Handle multiple files and indexes.
	// The index to store the byte offsets of the lines in the file.
	index *Index
}

type Index struct {
	// The file name to store the data.
	fileName string
	// A hash map containing byte offsets of the line in the file for each key.
	hashMap map[string]int64
	// Last byte offset of the file.
	lastByteOffset int64
	// A mutex to lock the file.
	mutex *sync.Mutex
}

type NotFoundError struct {
	message string
}

var KeyValueDelimiter = ":"
var RecordDelimiter = "\n"
var byteRecordDelimiter = byte('\n')

func NewStorage(index *Index) *storage {
	return &storage{
		index: index,
	}
}

func NewIndex(fileName string, hashMap map[string]int64, lastByteOffset int64) *Index {
	return &Index{
		fileName:       fileName,
		hashMap:        hashMap,
		lastByteOffset: lastByteOffset,
		mutex:          &sync.Mutex{},
	}
}

func (s *storage) Set(key string, value string) error {
	return s.index.Set(key, value)
}

func (s *storage) Get(key string) (string, error) {
	return s.index.Get(key)
}

func (s *storage) PrintHashMap() {
	s.index.PrintHashMap()
}

func (i *Index) Set(key string, value string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if strings.Contains(value, KeyValueDelimiter) || strings.Contains(value, RecordDelimiter) {
		return fmt.Errorf("value cannot contain (%s) or (%s)", KeyValueDelimiter, RecordDelimiter)
	}

	file, err := os.OpenFile(i.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := file.WriteString(fmt.Sprintf("%s%s%s%s", key, KeyValueDelimiter, value, RecordDelimiter))
	if err != nil {
		return err
	}
	i.hashMap[key] = i.lastByteOffset
	i.lastByteOffset += int64(bytes)

	return nil
}

func (i *Index) Get(key string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	file, err := os.Open(i.fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	offset, ok := i.hashMap[key]
	if !ok {
		return "", newNotFoundError(key)
	}

	reader := bufio.NewReader(file)
	var byteOffset int64 = 0
	for {
		line, err := reader.ReadString(byteRecordDelimiter)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		if byteOffset == int64(offset) {
			return strings.TrimSuffix(strings.Split(line, ":")[1], RecordDelimiter), nil
		}

		byteOffset += int64(len(line))
	}

	return "", errors.New("it should be in the file but it's not")
}

func (i *Index) PrintHashMap() {
	for key, value := range i.hashMap {
		fmt.Printf("key: %s, byte offset: %d\n", key, value)
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
