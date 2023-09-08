package main

import "sync"

type logSnippet struct {
	key int
	value string
}

var increment = 0
var storage = []logSnippet{}

var mutex = &sync.Mutex{}

func set(value string) int {
	mutex.Lock()
	defer mutex.Unlock()

	increment += 1
	storage = append(storage, logSnippet{key:increment, value:value})
	return increment
}

func get(key int) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	for index := len(storage) - 1; index >= 0; index-- {
		snippet := storage[index]
		if snippet.key == key {
			return snippet.value, nil
		}
	}
	return "", nil
}