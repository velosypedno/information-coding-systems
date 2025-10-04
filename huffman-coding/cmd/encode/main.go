package main

import (
	"fmt"
	"math"
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

	totalLength := 0
	counterMap := make(map[rune]int)
	for _, char := range text {
		counterMap[char]++
		totalLength++
	}
	freqMap := make(map[rune]float64)
	for char, count := range counterMap {
		freqMap[char] = float64(count) / float64(totalLength)
	}
	entropy := 0.0
	for _, freq := range freqMap {
		entropy += -1 * freq * math.Log2(freq)
	}
	fmt.Printf("Entropy: %f\n", entropy)

	nodes := make([]tree.Node[tree.Pair], 0, len(counterMap))
	for char, freq := range counterMap {
		nodes = append(nodes, tree.Node[tree.Pair]{Value: tree.Pair{Char: char, Freq: freq}})
	}

	root := tree.NewHuffmanTree(nodes)
	huffmanCodes := codes.NewHuffmanCodesMap(root)
	encoded, bitLength := codes.EncodeHuffman(text, huffmanCodes)
	fmt.Println(huffmanCodes)
	for r, code := range huffmanCodes {
		fmt.Printf("%c - %b\n", r, code.Code)
	}

	if err := codes.SaveEncodedContent(encodedFile, encoded, bitLength); err != nil {
		fmt.Printf("Error saving encoded content: %v\n", err)
		os.Exit(1)
	}

	if err := codes.SaveCodesMap(codesFile, huffmanCodes); err != nil {
		fmt.Printf("Error saving codes map: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Encoding complete!")

	averageLength := 0.0
	for char, code := range huffmanCodes {
		averageLength += float64(code.Length) * freqMap[char]
	}
	fmt.Printf("Average bit length: %f\n", averageLength)

	fmt.Printf("Encoded bits: %d\n", bitLength)
}
