package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/velosypedno/information-coding-systems/shennon-fano-coding/codes"
	"github.com/velosypedno/information-coding-systems/shennon-fano-coding/tree"
)

type ByFreq []tree.Pair

func (a ByFreq) Len() int {
	return len(a)
}

func (a ByFreq) Less(i, j int) bool {
	return a[i].Freq < a[j].Freq
}

func (a ByFreq) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

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

	pairs := make([]tree.Pair, 0, len(freqMap))
	for char, freq := range freqMap {
		pairs = append(pairs, tree.Pair{Char: char, Freq: freq})
	}
	sort.Sort(sort.Reverse(ByFreq(pairs)))

	tree := tree.NewShannonFanoTree(pairs)
	codesMap := codes.NewShennonFanoCodesMap(tree)

	for r, code := range codesMap {
		var codeBuilder strings.Builder
		for i := code.Length - 1; i >= 0; i-- {
			isZero := (0 == (code.Code&(1<<i))>>i)
			if isZero {
				codeBuilder.WriteRune('0')
			} else {
				codeBuilder.WriteRune('1')
			}
		}
		fmt.Printf("%c, %v - \n", r, codeBuilder.String())
	}

	averageLength := 0.0
	for char, code := range codesMap {
		averageLength += float64(code.Length) * freqMap[char]
	}
	fmt.Printf("Average bit length: %f\n", averageLength)

	encoded, bitLength := codes.EncodeShennonFano(text, codesMap)
	if err := codes.SaveEncodedContent(encodedFile, encoded, bitLength); err != nil {
		fmt.Printf("Error saving encoded content: %v\n", err)
		os.Exit(1)
	}

	if err := codes.SaveCodesMap(codesFile, codesMap); err != nil {
		fmt.Printf("Error saving codes map: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Encoding complete!")
	fmt.Printf("Encoded bits: %d\n", bitLength)
}
