package internal

import (
	"fmt"
	"slices"
	"strings"
)

const EOFSymbol rune = 0xffffff

type Model struct {
	symbolToIndex   map[rune]int
	indexToSymbol   []rune
	cumulativeLow   []uint32
	TotalCumulative uint64
	Size            int
}

func NewModel(freqMap map[rune]int) *Model {
	freqMap[EOFSymbol] = 1
	indexToSymbol := make([]rune, len(freqMap))
	i := 0
	for symbol := range freqMap {
		indexToSymbol[i] = symbol
		i++
	}
	slices.Sort(indexToSymbol)
	symbolToIndex := make(map[rune]int, len(freqMap))
	for i, symbol := range indexToSymbol {
		symbolToIndex[symbol] = i
	}

	cumulativeLow := make([]uint32, len(freqMap)+1)
	for i := 1; i < len(cumulativeLow); i++ {
		cumulativeLow[i] = cumulativeLow[i-1] + uint32(freqMap[indexToSymbol[i-1]])
	}

	return &Model{
		symbolToIndex:   symbolToIndex,
		indexToSymbol:   indexToSymbol,
		cumulativeLow:   cumulativeLow,
		TotalCumulative: uint64(cumulativeLow[len(freqMap)]),
	}

}

func (m *Model) GetBounds(r rune) (uint64, uint64) {
	index := m.symbolToIndex[r]
	return uint64(m.cumulativeLow[index]), uint64(m.cumulativeLow[index+1])
}

func (m *Model) GetSymbol(point uint64) rune {
	for i, lowBound := range m.cumulativeLow {
		if uint64(lowBound) <= point {
			continue
		} else {
			return m.indexToSymbol[i-1]
		}
	}
	return EOFSymbol
}

func (m *Model) String() string {
	var sb strings.Builder
	sb.WriteString("Model: frequency & cumulative ranges\n")
	sb.WriteString("=================================\n")
	sb.WriteString(fmt.Sprintf("%-8s %-8s %-8s %-8s\n", "Symbol", "Index", "Low", "High"))
	sb.WriteString("---------------------------------\n")

	symbols := make([]rune, len(m.symbolToIndex))
	for r, idx := range m.symbolToIndex {
		symbols[idx] = r
	}

	for idx, r := range symbols {
		low := m.cumulativeLow[idx]
		high := m.cumulativeLow[idx+1]

		var sym string
		switch r {
		case '\n':
			sym = `\n`
		case '\t':
			sym = `\t`
		case ' ':
			sym = `' '`
		case EOFSymbol:
			sym = "EOF"
		default:
			sym = string(r)
		}

		sb.WriteString(fmt.Sprintf("%-8s %-8d %-8d %-8d\n", sym, idx, low, high))
	}

	return sb.String()
}
