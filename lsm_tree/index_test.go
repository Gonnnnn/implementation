package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("Test Case 1", func(t *testing.T) {
		// TODO: Make Set testable and Implement a test.
	})
}

func TestGet(t *testing.T) {
    t.Run("Test Case 1", func(t *testing.T) {
		// TODO: Make Get testable and Implement a test.
	})
}

func TestPrintHashMap(t *testing.T) {
	t.Run("Test Case 1", func(t *testing.T) {
		// TODO: Make PrintHashMap testable and Implement a test.
	})
}

func TestName(t *testing.T) {
	t.Run("Returns the name of the index", func(t *testing.T) {
		index := NewIndex("test-file-path", make(map[string]int64), 0)
		name := index.Name()
		assert.Equal(t, name, "test-file-path")
	})
}

func TestSize(t *testing.T) {
	t.Run("Returns the size of the index", func(t *testing.T) {
		index := NewIndex("test-file-path", make(map[string]int64), 0)
		size := index.Size()
		assert.Equal(t, size, int64(0))
	})
}