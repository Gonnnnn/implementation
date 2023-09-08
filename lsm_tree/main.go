package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("Enter words. Enter newline to pass the word. Enter an empty string to break.")
	
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
		fmt.Printf("The key: %d\n", set(input))
	}

	fmt.Println("Enter a key to find data. Enter 0 to break.")
	for {
		var input string
		fmt.Print("Enter: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf("Error while reading a key: %v", err)
			continue
		}
		if input == "0" {
			break
		}
		key, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error while converting the key into int type:", err)
			continue
		}

		value, err := get(key)
		if err != nil {
			fmt.Println("No value corresponding to the key.")
			continue
		}

		fmt.Printf("Value: %s\n", value)
	}
}