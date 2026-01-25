package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type BTree[K cmp.Ordered, V any] struct {
	Key   K
	Value V
	Empty bool
	Left  *BTree[K, V]
	Right *BTree[K, V]
}

func NewBTree[K cmp.Ordered, V any]() *BTree[K, V] {
	return &BTree[K, V]{
		Empty: true,
	}
}

func (t *BTree[K, V]) AddOrUpdate(key K, value V) bool {
	if t.Empty || t.Key == key {
		t.Key = key
		t.Value = value

		wasEmpty := t.Empty
		t.Empty = false

		return wasEmpty
	}

	if key < t.Key {
		if t.Left == nil {
			t.Left = NewBTree[K, V]()
		}
		return t.Left.AddOrUpdate(key, value)
	}

	if t.Right == nil {
		t.Right = NewBTree[K, V]()
	}
	return t.Right.AddOrUpdate(key, value)
}

func (t *BTree[K, V]) Delete(key K) bool {
	if t.Key == key {
		t.Empty = true
		return true
	}

	var node *BTree[K, V]
	if key < t.Key {
		node = t.Left
	} else {
		node = t.Right
	}

	if node == nil {
		return false
	}

	return node.Delete(key)
}

func (t *BTree[K, V]) IsSet(key K) bool {
	if t.Key == key {
		return !t.Empty
	}

	var node *BTree[K, V]
	if key < t.Key {
		node = t.Left
	} else {
		node = t.Right
	}

	if node == nil {
		return false
	}

	return node.IsSet(key)
}

func (t *BTree[K, V]) ForEach(action func(K, V)) {
	if t.Left != nil {
		t.Left.ForEach(action)
	}

	if !t.Empty {
		action(t.Key, t.Value)
	}

	if t.Right != nil {
		t.Right.ForEach(action)
	}
}

type OrderedMap[K cmp.Ordered, V any] struct {
	tree *BTree[K, V]
	len  int
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		tree: NewBTree[K, V](),
	}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	if m.tree.AddOrUpdate(key, value) {
		m.len++
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	if m.tree.Delete(key) {
		m.len--
	}
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.tree.IsSet(key)
}

func (m *OrderedMap[K, V]) Size() int {
	return m.len
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.tree.ForEach(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
