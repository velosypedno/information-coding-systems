package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
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
	if err != nil {
		return "", "", "", err
	}
	return os.Args[1], os.Args[2], os.Args[3], nil
}

func main() {
	// step 1: parse args
	encodedFileName, freqTableFileName, outputFileName, err := parseArgs()
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Usage: program <encoded_file> <freq_table_file> <output_file>")
		os.Exit(1)
	}

	// step 2: read encoded file
	encodedFile, err := os.Open(encodedFileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer encodedFile.Close()
	inputStream := internal.NewInputStream(encodedFile)

	// step 3: build model
	freqTableFile, err := os.Open(freqTableFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer freqTableFile.Close()

	freqTableEncoder := gob.NewDecoder(freqTableFile)
	var freqMap internal.FreqMap
	err = freqTableEncoder.Decode(&freqMap)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	model := internal.NewModel(freqMap)
	fmt.Println(model)

	// step 4: encode
	var lowerBound uint32 = 0
	var upperBound uint32 = ^lowerBound
	var encodedBits uint32 = 0
	for i := 0; i < 32; i++ {
		bit, err := inputStream.ReadBit()
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		encodedBits = encodedBits<<1 | (bit & 1)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer outputFile.Close()
	bufWriter := bufio.NewWriter(outputFile)

	fmt.Println("Start decoding")
	fmt.Printf("(Lower: %b, Upper: %b)\n", lowerBound, upperBound)

	i := 0
	for {

		var currentRange uint64 = uint64(upperBound) - uint64(lowerBound) + 1

		var scaledSymbol uint64 = ((uint64(encodedBits)-uint64(lowerBound)+1)*model.TotalCumulative - 1) / currentRange

		symbol := model.GetSymbol(scaledSymbol)
		if i < 10 {
			fmt.Printf("Encode bits: %x\n", encodedBits)
			fmt.Printf("Lower bound: %x\n", lowerBound)
			fmt.Printf("Symbol: %c\n", symbol)
			fmt.Printf("Scaled: %d\n", scaledSymbol)
			i++
		}
		if symbol == internal.EOFSymbol {
			break
		}
		if _, err := bufWriter.WriteRune(symbol); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		symbolLow, symbolHigh := model.GetBounds(symbol)
		ubOffset := (currentRange * symbolHigh) / model.TotalCumulative
		if ubOffset > 0 {
			ubOffset = ubOffset - 1
		}
		upperBound = lowerBound + uint32(ubOffset)
		lowerBound = lowerBound + uint32((currentRange*symbolLow)/model.TotalCumulative)

		// fmt.Printf("Compute bounds: (%b, %b)\n", lowerBound, upperBound)

		for {
			if upperBound>>31 == lowerBound>>31 {
				upperBound <<= 1
				upperBound |= 1
				lowerBound <<= 1

				encodedBits <<= 1
				bit, err := inputStream.ReadBit()
				if err != nil && err != io.EOF {
					fmt.Println("Error", err)
					os.Exit(1)
				}
				encodedBits |= (bit & 1)
			} else if (lowerBound>>30&1) == 1 && (upperBound>>30&1) == 0 {
				upperBound = upperBound << 1
				upperBound = upperBound | (1 << 31)
				upperBound = upperBound | 1

				lowerBound = lowerBound << 1
				lowerBound = lowerBound & ((1 << 31) - 1)

				var mostSignificantBit uint32 = encodedBits >> 31
				var rest uint32 = encodedBits & 0x3fffffff
				bit, err := inputStream.ReadBit()
				if err != nil && err != io.EOF {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				encodedBits = (mostSignificantBit << 31) | (rest << 1) | (bit & 1)
			} else {
				break
			}
		}
		if symbol == internal.EOFSymbol {
			break
		}
	}
	err = bufWriter.Flush()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}
