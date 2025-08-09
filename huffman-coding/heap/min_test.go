package heap_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/velosypedno/information-coding-systems/huffman-coding/heap"
)

func TestInsert(t *testing.T) {
	h := heap.NewMinHeap(func(a, b int) bool {
		return a < b
	})
	h.Insert(6)
	require.Equal(t, 6, h.Peek())
	h.Insert(3)
	h.Insert(4)
	require.Equal(t, 3, h.Peek())
	h.Insert(1)
	require.Equal(t, 1, h.Peek())
}

func TestExtractMin(t *testing.T) {
	h := heap.NewMinHeap(func(a, b int) bool {
		return a < b
	})
	h.Insert(6)
	h.Insert(3)
	h.Insert(4)
	h.Insert(1)
	require.Equal(t, 1, h.ExtractMin())
	require.Equal(t, 3, h.ExtractMin())
	require.Equal(t, 4, h.ExtractMin())
	require.Equal(t, 6, h.ExtractMin())
}

func TestCombineInsertExtract(t *testing.T) {
	h := heap.NewMinHeap(func(a, b int) bool {
		return a < b
	})
	h.Insert(6)
	h.Insert(3)
	h.Insert(4)
	h.Insert(1)
	require.Equal(t, 1, h.ExtractMin())
	require.Equal(t, 3, h.ExtractMin())
	require.Equal(t, 4, h.ExtractMin())
	h.Insert(2)
	h.Insert(3)
	require.Equal(t, 2, h.ExtractMin())
	require.Equal(t, 3, h.ExtractMin())
	require.Equal(t, 6, h.ExtractMin())
}

func TestExtractMinPanicsOnEmptyHeap(t *testing.T) {
	h := heap.NewMinHeap(func(a, b int) bool {
		return a < b
	})
	require.PanicsWithValue(t, "heap is empty", func() {
		h.ExtractMin()
	})
}

func TestPeekPanicsOnEmptyHeap(t *testing.T) {
	h := heap.NewMinHeap(func(a, b int) bool {
		return a < b
	})
	require.PanicsWithValue(t, "heap is empty", func() {
		h.Peek()
	})
}
