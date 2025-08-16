package main

import (
	"fmt"
	"os"

	"github.com/velosypedno/information-coding-systems/huffman-coding/codes"
	"github.com/velosypedno/information-coding-systems/huffman-coding/tree"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: program <input_file> <encoded_file> <codes_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	encodedFile := os.Args[2]
	codesFile := os.Args[3]

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

	nodes := make([]tree.Node[tree.Pair], 0, len(freqMap))
	for char, freq := range freqMap {
		nodes = append(nodes, tree.Node[tree.Pair]{Value: tree.Pair{Char: char, Freq: freq}})
	}

	root := tree.NewHuffmanTree(nodes)
	huffmanCodes := codes.NewHuffmanCodesMap(root)
	encoded, bitLength := codes.EncodeHuffman(text, huffmanCodes)

	if err := codes.SaveEncodedContent(encodedFile, encoded, bitLength); err != nil {
		fmt.Printf("Error saving encoded content: %v\n", err)
		os.Exit(1)
	}

	if err := codes.SaveCodesMap(codesFile, huffmanCodes); err != nil {
		fmt.Printf("Error saving codes map: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Encoding complete!")
	fmt.Printf("Encoded bits: %d\n", bitLength)
}
