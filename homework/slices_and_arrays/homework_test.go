package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type CircularQueue[T Integer] struct {
	values []T
	len    int
	cap    int
	start  int
	end    int
}

func NewCircularQueue[T Integer](size int) CircularQueue[T] {
	if size < 0 {
		panic("CircularQueue should have positive size")
	}

	return CircularQueue[T]{
		values: make([]T, size),
		cap:    size,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	if !q.Empty() {
		q.end = (q.end + 1) % len(q.values)
	}

	q.values[q.end] = value
	q.len++

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.len--

	if !q.Empty() {
		q.start = (q.start + 1) % len(q.values)
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.start]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.end]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.len == q.cap
}

func TestCircularQueue(t *testing.T) {
	var queueSize int
	var queue CircularQueue[int]

	/* queueSize = -1 */
	queueSize = -1
	assert.Panics(t, func() {
		queue = NewCircularQueue[int](queueSize)
	})

	/* queueSize = 0 */
	queueSize = 0
	queue = NewCircularQueue[int](queueSize)

	// очередь с нулевой длиной и полная и пустая одновременно
	assert.True(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.False(t, queue.Push(1))

	/* queueSize = 1 */
	queueSize = 1
	queue = NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.False(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{1}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 1, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(2))

	assert.True(t, reflect.DeepEqual([]int{2}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 2, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	/* queueSize > 1 */
	queueSize = 3
	queue = NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
