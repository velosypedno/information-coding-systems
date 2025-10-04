package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/velosypedno/information-coding-systems/arithmetic-coding/internal"
)

func checkOsArgs() error {
	if len(os.Args) < 4 {
		return errors.New("not enough arguments")
	}
	return nil
}

func parseArgs() (string, string, string, error) {
	err := checkOsArgs()
	return os.Args[1], os.Args[2], os.Args[3], err
}

func main() {
	// step 1: parse args
	inputFileName, outputFileName, freqTableFileName, err := parseArgs()
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Usage: program <input_file> <output_file>")
		os.Exit(1)
	}

	// step 2: read input file
	inputBytes, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	inputString := string(inputBytes)

	// step 3: build model
	freqMap := internal.NewFreqMapFromString(inputString)
	freqTableFile, err := os.Create(freqTableFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer freqTableFile.Close()
	freqTableEncoder := gob.NewEncoder(freqTableFile)
	err = freqTableEncoder.Encode(freqMap)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	model := internal.NewModel(freqMap)
	fmt.Println(model)

	// step 4: encode
	var lowerBound uint32 = 0
	var upperBound uint32 = ^lowerBound
	underflowCounter := 0

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println(outputFile)
	defer outputFile.Close()
	bitStream := internal.NewOutputStream(outputFile)

	symbols := []rune(inputString)
	// fmt.Println("Start encoding")
	// fmt.Printf("(Lower: %b, Upper: %b)\n", lowerBound, upperBound)

	for i := 0; i < len(symbols)+1; i++ {
		var symbol rune
		if i == len(symbols) {
			//fmt.Print("EOF\n")
			symbol = internal.EOFSymbol
		} else {
			symbol = symbols[i]
			//fmt.Printf("%c\n", symbol)
		}
		var currentRange uint64 = uint64(upperBound) + 1 - uint64(lowerBound)
		symbolLowBound, symbolHighBound := model.GetBounds(symbol)
		upperBound = lowerBound + uint32((currentRange*symbolHighBound)/model.TotalCumulative) - 1
		lowerBound = lowerBound + uint32((currentRange*symbolLowBound)/model.TotalCumulative)

		// fmt.Printf("Compute bounds: (%b, %b)\n", lowerBound, upperBound)

		for {
			if upperBound>>31 == lowerBound>>31 {
				b := upperBound >> 31
				// fmt.Printf("Most significant bits matches: %d\n", b)
				bitStream.WriteBit(b)
				for i := 0; i < underflowCounter; i++ {
					if b == 0 {
						bitStream.WriteBit(1)
					} else {
						bitStream.WriteBit(0)
					}
				}

				underflowCounter = 0
				upperBound <<= 1
				upperBound |= 1
				lowerBound <<= 1
			} else if (lowerBound>>30&1) == 1 && (upperBound>>30&1) == 0 {
				underflowCounter++

				upperBound = upperBound << 1
				upperBound = upperBound | (1 << 31)
				upperBound = upperBound | 1

				lowerBound = lowerBound << 1
				lowerBound = lowerBound & ((1 << 31) - 1)
			} else {
				break
			}
		}
		if symbol == internal.EOFSymbol {
			break
		}
	}
	bitStream.WriteBit(0)
	bitStream.WriteBit(1)
	bitStream.FlushToByte(1)
}
