package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	var res []uintptr

	checked := make(map[uintptr]struct{})

	for _, stack := range stacks {
		for _, ptr := range stack {
			var tracePtr func(ptr uintptr)
			tracePtr = func(ptr uintptr) {
				if ptr == 0x00 {
					return
				}

				if _, ok := checked[ptr]; ok {
					return
				}
				checked[ptr] = struct{}{}

				res = append(res, ptr)

				var ptrptr = (*uintptr)(unsafe.Pointer(ptr))
				if ptrptr != nil {
					tracePtr(*ptrptr)
				}
			}
			tracePtr(ptr)
		}
	}

	return res
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3
	var heapPointer5 *int = &heapObjects[3]
	var heapPointer6 **int = &heapPointer5
	var heapPointer7 *int = &heapObjects[4]
	var heapPointer8 **int = &heapPointer7
	var heapPointer9 ***int = &heapPointer8

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapPointer3)),
			0x00, 0x00, uintptr(unsafe.Pointer(&heapPointer6)), 0x00,
			0x00, uintptr(unsafe.Pointer(&heapPointer9)), 0x00, 0x00,
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
		uintptr(unsafe.Pointer(&heapPointer6)),
		uintptr(unsafe.Pointer(&heapPointer5)),
		uintptr(unsafe.Pointer(&heapPointer9)),
		uintptr(unsafe.Pointer(&heapPointer8)),
		uintptr(unsafe.Pointer(&heapPointer7)),
		uintptr(unsafe.Pointer(&heapObjects[4])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
