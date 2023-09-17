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

type Index struct {
	// The file path to store the data.
	filePath string
	// A hash map containing byte offsets of the line in the file for each key.
	hashMap map[string]int64
	// Last byte offset of the file.
	lastByteOffset int64
	// A mutex to lock the file.
	mutex *sync.Mutex
}

var (
	// Variables here can be overwritten in unit tests.
	// https://github.com/GoogleCloudPlatform/esp-v2/blob/master/src/go/gcsrunner/fetch_config.go#L38C1-L38C1
	osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) { return os.OpenFile(name, flag, perm)}

)

var KeyValueDelimiter = ":"
var RecordDelimiter = "\n"
var byteRecordDelimiter = byte('\n')

func NewIndex(filePath string, hashMap map[string]int64, lastByteOffset int64) *Index {
	return &Index{
		filePath:       filePath,
		hashMap:        hashMap,
		lastByteOffset: lastByteOffset,
		mutex:          &sync.Mutex{},
	}
}

func (i *Index) Set(key string, value string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if strings.Contains(value, KeyValueDelimiter) || strings.Contains(value, RecordDelimiter) {
		return fmt.Errorf("value cannot contain (%s) or (%s)", KeyValueDelimiter, RecordDelimiter)
	}

	file, err := osOpenFile(i.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	file, err := os.Open(i.filePath)
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

func (i *Index) Name() string {
	return i.filePath
}

func (i *Index) Size() int64 {
	return i.lastByteOffset
}