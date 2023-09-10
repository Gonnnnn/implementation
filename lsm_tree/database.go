package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
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

func (s *storage) Set(value string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	file, err := os.OpenFile(s.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	key := generateRandomID()
	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	if err != nil {
		return "", err
	}

	return key, nil
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

func generateRandomID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(100_000_000_000))
}
