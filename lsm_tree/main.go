package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fileName := "storage/log.txt" // You could fix the code to get it from the command line argument.

	storage, err := initializeStorage(fileName)
	if err != nil {
		fmt.Printf("Error while initializing a storage: %v\n", err)
		return
	}

	fmt.Println("+++++ Simple Database +++++")
	fmt.Println("====================================")
	fmt.Println("Enter the following number to execute them.")
	fmt.Printf("1: Append Enter a pair of key and value seprated by \"%s\".\nAnd then enter newline to pass the word.\n", keyValueDelimiter)
	fmt.Println("2: GetEnter a key to find data.")
	fmt.Println("3: Show all the key-value pairs.")
	fmt.Println("0: Quit")

	for {
		fmt.Println("====================================")
		option, err := readStdin("Enter the option number: ")
		if err != nil {
			fmt.Printf("Error while reading an option: %v", err)
			continue
		}

		switch option {
		case "1":
			input, err := readStdin("Enter a pair of key value: ")
			if err != nil {
				fmt.Printf("Error while reading an input: %v\n", err)
				continue
			}

			key, value := splitIntoKeyValue(input)
			if err := storage.Set(key, value); err != nil {
				fmt.Printf("Error while setting a value: %v\n", err)
				continue
			}
			fmt.Printf("key: %s\n", key)

		case "2":
			input, err := readStdin("Enter a key: ")
			if err != nil {
				fmt.Printf("Error while reading a key: %v\n", err)
				continue
			}

			value, err := storage.Get(input)
			if err != nil {
				fmt.Printf("Error while getting a value: %v\n", assertGetError(err))
				continue
			}
			fmt.Printf("Value: %s\n", value)
		case "3":
			fmt.Println("All the key-value pairs are as follows.")
			storage.PrintHashMap()

		case "0":
			fmt.Println("Bye!")
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func initializeStorage(fileName string) (*storage, error) {
	fmt.Println("Loading the data...")
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hashMap := make(map[string]int64)
	reader := bufio.NewReader(file)
	var byteOffset int64 = 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		fmt.Printf("line: %s | byte offset: %d\n", strings.TrimSuffix(line, "\n"), byteOffset)
		key, _ := splitIntoKeyValue(line)
		hashMap[key] = byteOffset
		byteOffset += int64(len(line))
	}

	return New(fileName, hashMap, byteOffset), nil

}

func readStdin(message string) (string, error) {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func splitIntoKeyValue(input string) (string, string) {
	parts := strings.Split(input, keyValueDelimiter)
	return parts[0], parts[1]
}

func assertGetError(err error) error {
	notFoundError := &NotFoundError{}
	if errors.As(err, &notFoundError) {
		return notFoundError
	}
	return err
}
