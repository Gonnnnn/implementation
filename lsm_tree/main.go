package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	storage := New()

	fmt.Println("+++++ Simple Database +++++")
	fmt.Println("====================================")
	fmt.Println("Enter a pair of key and value seprated by \":\".\nAnd then enter newline to pass the word.\nEnter \"QUIT\" to break.")
	fmt.Println("====================================")

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
		key, value := splitIntoKeyValue(input)

		if err := storage.Set(key, value); err != nil {
			fmt.Printf("Error while setting a value: %v", err)
			continue
		}
		fmt.Printf("The key: %s\n", key)
	}

	fmt.Println("====================================")
	fmt.Println("Enter a key to find data. Enter \"QUIT\" to break.")
	fmt.Println("====================================")
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

func splitIntoKeyValue(input string) (string, string) {
	parts := strings.Split(input, ":")
	return parts[0], parts[1]
}

func assertGetError(err error) error {
	notFoundError := &NotFoundError{}
	if errors.As(err, &notFoundError) {
		return notFoundError
	}
	return err
}
