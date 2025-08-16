package tree

import (
	"fmt"

	"github.com/velosypedno/information-coding-systems/huffman-coding/queue"
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
	Freq int
}

func NewHuffmanTree(nodes []Node[Pair]) *Node[Pair] {
	q := queue.NewMinPriorityQueue[Node[Pair]]()
	for _, node := range nodes {
		q.Enqueue(node, float64(node.Value.Freq))
	}
	var root *Node[Pair] = nil
	for q.Size() > 1 {
		n1 := q.Dequeue().Value
		n2 := q.Dequeue().Value
		root = &Node[Pair]{
			Left:  &n1,
			Right: &n2,
			Value: Pair{Char: 0, Freq: n1.Value.Freq + n2.Value.Freq},
		}
		q.Enqueue(*root, float64(root.Value.Freq))
	}
	return root
}
