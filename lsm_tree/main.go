package main

import (
	"errors"
	"fmt"
)

func main() {
	storage := New()

	fmt.Println("Enter words. Enter newline to pass the word. Enter \"QUIT\" to break.")

	for {
		var input string
		fmt.Print("Enter: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf("Error while reading an input: %v", err)
			continue
		}
		if input == "QUIT" {
			break
		}

		key, err := storage.Set(input)
		if err != nil {
			fmt.Printf("Error while setting a value: %v", err)
			continue
		}
		fmt.Printf("The key: %s\n", key)
	}

	fmt.Println("Enter a key to find data. Enter \"QUIT\" to break.")
	for {
		var input string
		fmt.Print("Enter: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf("Error while reading a key: %v", err)
			continue
		}
		if input == "QUIT" {
			break
		}

		value, err := storage.Get(input)
		if err != nil {
			fmt.Printf("Error while getting a value: %v", assertGetError(err))
		}

		fmt.Printf("Value: %s\n", value)
	}
}

func assertGetError(err error) error {
	notFoundError := &NotFoundError{}
	if errors.As(err, &notFoundError) {
		return notFoundError
	}
	return err
}
