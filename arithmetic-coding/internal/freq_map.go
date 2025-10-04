package internal

import (
	"fmt"
	"strings"
)

type FreqMap map[rune]int

func NewFreqMapFromString(data string) FreqMap {
	minSize := 256
	runeToFreq := make(map[rune]int, minSize)
	for _, r := range data {
		runeToFreq[r]++
	}
	return runeToFreq
}

func (m FreqMap) String() string {
	var sb strings.Builder
	sb.WriteString("Frequency map of symbols in text:\n")
	sb.WriteString("-------------------------------\n")

	maxSymbolsInline := 5
	counter := 0

	for symbol, freq := range m {
		var sym string
		switch symbol {
		case '\n':
			sym = `\n`
		case '\t':
			sym = `\t`
		case ' ':
			sym = `' '`
		default:
			sym = string(symbol)
		}

		sb.WriteString(fmt.Sprintf("(%-3s: %4d)  ", sym, freq))
		counter++
		if counter == maxSymbolsInline {
			sb.WriteRune('\n')
			counter = 0
		}
	}

	if counter != 0 {
		sb.WriteRune('\n')
	}

	return sb.String()
}
