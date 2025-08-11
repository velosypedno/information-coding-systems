package queue_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/velosypedno/information-coding-systems/huffman-coding/queue"
)

func TestMinPriorityQueue_BasicOperations(t *testing.T) {
	q := queue.NewMinPriorityQueue[string]()

	require.True(t, q.IsEmpty())
	require.Equal(t, 0, q.Size())

	q.Enqueue("low", 5)
	q.Enqueue("medium", 3)
	q.Enqueue("high", 1)

	require.False(t, q.IsEmpty())
	require.Equal(t, 3, q.Size())

	item := q.Peek()
	require.Equal(t, "high", item.Value)
	require.Equal(t, 1.0, item.Priority)
	require.Equal(t, 3, q.Size())

	item = q.Dequeue()
	require.Equal(t, "high", item.Value)
	require.Equal(t, 1.0, item.Priority)

	item = q.Dequeue()
	require.Equal(t, "medium", item.Value)
	require.Equal(t, 3.0, item.Priority)

	item = q.Dequeue()
	require.Equal(t, "low", item.Value)
	require.Equal(t, 5.0, item.Priority)

	require.True(t, q.IsEmpty())
	require.Equal(t, 0, q.Size())
}

func TestMinPriorityQueue_EmptyQueuePanics(t *testing.T) {
	q := queue.NewMinPriorityQueue[int]()

	require.Panics(t, func() {
		q.Dequeue()
	})

	require.Panics(t, func() {
		q.Peek()
	})
}
