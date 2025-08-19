package tree

import (
	"fmt"
)

type Node[T any] struct {
	Left  *Node[T]
	Right *Node[T]
	Value T
}

func (n *Node[T]) String() string {
	if n == nil {
		return "()"
	}
	if n.Left == nil && n.Right == nil {
		return fmt.Sprintf("(%v)", n.Value)
	}
	return fmt.Sprintf("(%v %v %v)", n.Value, n.Left, n.Right)
}

type Pair struct {
	Char rune
	Freq float64
}

func NewShannonFanoTree(pairs []Pair) *Node[rune] {
	if len(pairs) == 1 {
		return &Node[rune]{Value: pairs[0].Char}
	}
	var totalSum float64
	for _, p := range pairs {
		totalSum += p.Freq
	}

	var splitSum float64
	splitIndex := 0
	for i, p := range pairs {
		splitSum += p.Freq
		if splitSum >= totalSum/2 {
			splitIndex = i + 1
			break
		}
	}

	if splitIndex == 0 || splitIndex == len(pairs) {
		if len(pairs) > 1 {
			splitIndex = 1
		} else {
			return &Node[rune]{Value: pairs[0].Char}
		}
	}

	left := pairs[:splitIndex]
	right := pairs[splitIndex:]

	leftNode := NewShannonFanoTree(left)
	rightNode := NewShannonFanoTree(right)

	return &Node[rune]{Left: leftNode, Right: rightNode, Value: 0}
}
