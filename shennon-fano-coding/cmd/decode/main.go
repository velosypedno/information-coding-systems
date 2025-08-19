package main

import (
	"fmt"
	"os"

	"github.com/velosypedno/information-coding-systems/shennon-fano-coding/codes"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: program <encoded_file> <codes_file> <output_file>")
		os.Exit(1)
	}

	encodedFile := os.Args[1]
	codesFile := os.Args[2]
	outputFile := os.Args[3]

	encodedData, bitLength, err := codes.ReadEncodedContent(encodedFile)
	if err != nil {
		fmt.Printf("Error reading encoded file: %v\n", err)
		os.Exit(1)
	}

	huffmanCodes, err := codes.ReadCodesMap(codesFile)
	if err != nil {
		fmt.Printf("Error reading codes map: %v\n", err)
		os.Exit(1)
	}

	decodedText, err := codes.DecodeHuffman(encodedData, bitLength, huffmanCodes)
	if err != nil {
		fmt.Printf("Error decoding Huffman data: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, []byte(decodedText), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Decoding complete! Output written to", outputFile)
}
