package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type storage struct {
	fileName string
	mutex    *sync.Mutex
}

type NotFoundError struct {
	message string
}

func New() *storage {
	return &storage{
		fileName: "storage/log.txt",
		mutex:    &sync.Mutex{},
	}
}

func (s *storage) Set(key string, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	if err != nil {
		return err
	}

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

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := len(lines) - 1; i >= 0; i-- {
		// parse the key from the line. it is split by ":"
		if strings.Split(lines[i], ":")[0] == key {
			return strings.Split(lines[i], ":")[1], nil
		}
	}
	return "", NewNotFoundError(key)
}

func NewNotFoundError(key string) *NotFoundError {
	return &NotFoundError{
		message: fmt.Sprintf("no value corresponding to the key: %s", key),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}
