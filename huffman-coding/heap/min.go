package heap

import (
	"fmt"
)

type Less[T any] func(a, b T) bool

type MinHeap[T any] struct {
	data []T
	less Less[T]
}

func NewMinHeap[T comparable](comp Less[T]) *MinHeap[T] {
	return &MinHeap[T]{
		data: make([]T, 0),
		less: comp,
	}
}

func (h *MinHeap[T]) parent(i int) int {
	return (i - 1) / 2
}

func (h *MinHeap[T]) left(i int) int {
	return 2*i + 1
}

func (h *MinHeap[T]) right(i int) int {
	return 2*i + 2
}

func (h *MinHeap[T]) heapifyUp(i int) {
	for i > 0 {
		parent := h.parent(i)
		if h.less(h.data[i], h.data[parent]) {
			h.data[i], h.data[parent] = h.data[parent], h.data[i]
			i = parent
		} else {
			break
		}
	}
}

func (h *MinHeap[T]) heapifyDown(i int) {
	for {
		left := h.left(i)
		right := h.right(i)
		smallest := i

		if left < len(h.data) && h.less(h.data[left], h.data[smallest]) {
			smallest = left
		}
		if right < len(h.data) && h.less(h.data[right], h.data[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.data[i], h.data[smallest] = h.data[smallest], h.data[i]
		i = smallest
	}
}

func (h *MinHeap[T]) Len() int {
	return len(h.data)
}

func (h *MinHeap[T]) Insert(v T) {
	h.data = append(h.data, v)
	h.heapifyUp(len(h.data) - 1)
}

func (h *MinHeap[T]) ExtractMin() T {
	min := h.data[0]
	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	h.heapifyDown(0)
	return min
}

func (h *MinHeap[T]) Peek() T {
	return h.data[0]
}

func (h *MinHeap[T]) String() string {
	return fmt.Sprint(h.data)
}
