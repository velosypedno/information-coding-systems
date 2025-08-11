package queue

import "github.com/velosypedno/information-coding-systems/huffman-coding/heap"

type Item[T any] struct {
	Value    any
	Priority float64
}

type MinPriorityQueue[T any] struct {
	heap *heap.MinHeap[Item[T]]
}

func NewMinPriorityQueue[T any]() *MinPriorityQueue[T] {
	return &MinPriorityQueue[T]{
		heap: heap.NewMinHeap(func(a, b Item[T]) bool { return a.Priority < b.Priority }),
	}
}

func (q *MinPriorityQueue[T]) Enqueue(value T, priority float64) {
	q.heap.Insert(Item[T]{
		Value:    value,
		Priority: priority,
	})
}

func (q *MinPriorityQueue[T]) Dequeue() Item[T] {
	if q.heap.Len() == 0 {
		panic("queue is empty")
	}
	return q.heap.ExtractMin()
}

func (q *MinPriorityQueue[T]) Peek() Item[T] {
	if q.heap.Len() == 0 {
		panic("queue is empty")
	}
	return q.heap.Peek()
}

func (q *MinPriorityQueue[T]) IsEmpty() bool {
	return q.heap.Len() == 0
}

func (q *MinPriorityQueue[T]) Size() int {
	return q.heap.Len()
}
