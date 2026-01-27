package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Реализуйте функцию `func Map(data []int, action func(int) int) []int`,
// которая принимает функцию `action` и срез данных `data`.
// Функция `Map` должна применить функцию `action` к каждому элементу среза `data`
// и вернуть новый срез с результатами.
func Map[T, R any](data []T, action func(T) R) []R {
	if data == nil {
		return nil
	}

	result := make([]R, len(data))

	for i, v := range data {
		result[i] = action(v)
	}

	return result
}

// Реализуйте функцию `func Filter(data []int, action func(int) bool) []int`,
// которая принимает функцию `action` и срез данных `data`. Функция `Filter` должна вернуть новый срез,
// содержащий только те элементы `data`, для которых функция `action` возвращает `true`.
func Filter[T any](data []T, action func(T) bool) []T {
	if data == nil {
		return nil
	}

	result := make([]T, 0, len(data))

	for _, v := range data {
		if action(v) {
			result = append(result, v)
		}
	}

	return result
}

// Реализуйте функцию `func Reduce(data []int, initial int, action func(int, int) int) int`,
// которая принимает функцию `action` (функцию двух аргументов), срез данных `data`
// и начальное значение `initial`. Функция `Reduce` должна применить функцию `action` к каждому элементу `data`
// и начальному значению `initial`, накапливая результат.
func Reduce[T, R any](data []T, initial R, action func(T, R) R) R {
	if data == nil {
		return initial
	}

	result := initial

	for _, v := range data {
		result = action(v, result)
	}

	return result
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) int
		result []int
	}{
		"nil numbers": {
			action: func(number int) int {
				return -number
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) int {
				return -number
			},
			result: []int{},
		},
		"inc numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number + 1
			},
			result: []int{2, 3, 4, 5, 6},
		},
		"double numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) int {
				return number * number
			},
			result: []int{1, 4, 9, 16, 25},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		data   []int
		action func(int) bool
		result []int
	}{
		"nil numbers": {
			action: func(number int) bool {
				return number == 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(number int) bool {
				return number == 1
			},
			result: []int{},
		},
		"even numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(number int) bool {
				return number%2 == 0
			},
			result: []int{2, 4},
		},
		"positive numbers": {
			data: []int{-1, -2, 1, 2},
			action: func(number int) bool {
				return number > 0
			},
			result: []int{1, 2},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(test.data, test.action)
			assert.True(t, reflect.DeepEqual(test.result, result))
		})
	}
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		initial int
		data    []int
		action  func(int, int) int
		result  int
	}{
		"nil numbers": {
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"empty numbers": {
			data: []int{},
			action: func(lhs, rhs int) int {
				return 0
			},
		},
		"sum of numbers": {
			data: []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 15,
		},
		"sum of numbers with initial value": {
			initial: 10,
			data:    []int{1, 2, 3, 4, 5},
			action: func(lhs, rhs int) int {
				return lhs + rhs
			},
			result: 25,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(test.data, test.initial, test.action)
			assert.Equal(t, test.result, result)
		})
	}
}
