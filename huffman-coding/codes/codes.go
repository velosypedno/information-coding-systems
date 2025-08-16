package codes

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/velosypedno/information-coding-systems/huffman-coding/tree"
)

type HuffmanCode struct {
	Code   uint64
	Length int
}

func NewHuffmanCodesMap(n *tree.Node[tree.Pair]) map[rune]HuffmanCode {
	m := make(map[rune]HuffmanCode)

	var dfs func(node *tree.Node[tree.Pair], code uint64, length int)
	dfs = func(node *tree.Node[tree.Pair], code uint64, length int) {
		if node == nil {
			return
		}
		if node.Left == nil && node.Right == nil {
			m[node.Value.Char] = HuffmanCode{Code: code, Length: length}
			return
		}
		dfs(node.Left, code<<1, length+1)
		dfs(node.Right, (code<<1)+1, length+1)
	}
	dfs(n, 0, 0)
	return m
}

type HuffmanEncoded struct {
	Data      []byte
	BitLength int
}

func EncodeHuffman(input string, codes map[rune]HuffmanCode) ([]byte, int) {
	var result []byte
	var currentByte byte
	bitPos := 0
	totalBits := 0

	for _, ch := range input {
		code := codes[ch]

		for i := code.Length - 1; i >= 0; i-- {
			bit := (code.Code >> i) & 1
			currentByte = currentByte<<1 | byte(bit)
			bitPos++
			totalBits++

			if bitPos == 8 {
				result = append(result, currentByte)
				currentByte = 0
				bitPos = 0
			}
		}
	}

	if bitPos > 0 {
		currentByte = currentByte << (8 - bitPos)
		result = append(result, currentByte)
	}

	return result, totalBits
}

func DecodeHuffman(encoded []byte, bitLength int, codes map[rune]HuffmanCode) (string, error) {
	type codeKey struct {
		code   uint64
		length int
	}
	reverse := make(map[codeKey]rune)
	for ch, hc := range codes {
		reverse[codeKey{hc.Code, hc.Length}] = ch
	}

	var result []rune
	var currentCode uint64
	currentLen := 0
	bitsRead := 0

	for _, b := range encoded {
		for i := 7; i >= 0; i-- {
			if bitsRead >= bitLength {
				break
			}
			bit := (b >> i) & 1
			currentCode = (currentCode << 1) | uint64(bit)
			currentLen++
			bitsRead++

			if ch, ok := reverse[codeKey{currentCode, currentLen}]; ok {
				result = append(result, ch)
				currentCode = 0
				currentLen = 0
			}
		}
	}

	if currentLen != 0 {
		return "", fmt.Errorf("incomplete Huffman code at the end")
	}

	return string(result), nil
}

func SaveEncodedContent(fileName string, input []byte, bitLength int) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var header [8]byte
	binary.BigEndian.PutUint64(header[:], uint64(bitLength))
	_, err = file.Write(header[:])
	if err != nil {
		return err
	}
	_, err = file.Write(input)
	if err != nil {
		return err
	}
	return nil
}

func SaveCodesMap(fileName string, codes map[rune]HuffmanCode) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(codes)
	if err != nil {
		return err
	}
	return nil
}

func ReadCodesMap(fileName string) (map[rune]HuffmanCode, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var codes map[rune]HuffmanCode
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&codes)
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func ReadEncodedContent(filename string) ([]byte, int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, 0, err
	}

	if len(data) < 8 {
		return nil, 0, os.ErrInvalid
	}

	bitLength := int(binary.BigEndian.Uint64(data[:8]))
	encoded := data[8:]

	return encoded, bitLength, nil
}
