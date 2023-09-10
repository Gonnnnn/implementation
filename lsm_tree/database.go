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

var keyValueDelimiter = ":"
var recordDelimiter = "\n"
var byteRecordDelimiter = byte('\n')

func New(filename string, hashMap map[string]int64, lastByteOffset int64) *storage {
	return &storage{
		fileName:       filename,
		hashMap:        hashMap,
		lastByteOffset: lastByteOffset,
		mutex:          &sync.Mutex{},
	}
}

func (s *storage) Set(key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if strings.Contains(value, keyValueDelimiter) || strings.Contains(value, recordDelimiter) {
		return fmt.Errorf("value cannot contain (%s) or (%s)", keyValueDelimiter, recordDelimiter)
	}

	file, err := os.OpenFile(s.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := file.WriteString(fmt.Sprintf("%s%s%s%s", key, keyValueDelimiter, value, recordDelimiter))
	if err != nil {
		return err
	}
	s.hashMap[key] = s.lastByteOffset
	s.lastByteOffset += int64(bytes)

	return nil
}

func (s *storage) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.Open(s.fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	offset, ok := s.hashMap[key]
	if !ok {
		return "", NewNotFoundError(key)
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
			return strings.TrimSuffix(strings.Split(line, ":")[1], recordDelimiter), nil
		}

		byteOffset += int64(len(line))
	}

	return "", errors.New("it should be in the file but it's not")
}

func (s *storage) PrintHashMap() {
	for key, value := range s.hashMap {
		fmt.Printf("key: %s, byte offset: %d\n", key, value)
	}
}

func NewNotFoundError(key string) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf("no value corresponding to the key: %s", key),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}
