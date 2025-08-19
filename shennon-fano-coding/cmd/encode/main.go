package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: program <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	text := string(data)
	fmt.Println("File successfully read!")

	freqMap := make(map[rune]int)
	for _, char := range text {
		freqMap[char]++
	}

}
