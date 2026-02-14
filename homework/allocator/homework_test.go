package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	free := unsafe.Pointer(&memory[0])

	for i := range pointers {
		if pointers[i] == free {
			free = unsafe.Add(free, 1)
			continue
		}

		*(*byte)(free) = *(*byte)(pointers[i])
		*(*byte)(pointers[i]) = 0

		pointers[i] = free
		free = unsafe.Add(free, 1)
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentation2(t *testing.T) {
	var fragmentedMemory = []byte{
		0x00, 0x00, 0x00, 0x00,
		0x11, 0x00, 0x22, 0x00,
		0x00, 0x33, 0x00, 0x00,
		0x00, 0x00, 0x44, 0x55,
		0x00, 0x66, 0x77, 0x00,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[4]),
		unsafe.Pointer(&fragmentedMemory[6]),
		unsafe.Pointer(&fragmentedMemory[9]),
		unsafe.Pointer(&fragmentedMemory[14]),
		unsafe.Pointer(&fragmentedMemory[15]),
		unsafe.Pointer(&fragmentedMemory[17]),
		unsafe.Pointer(&fragmentedMemory[18]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
		unsafe.Pointer(&fragmentedMemory[4]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[6]),
	}

	var defragmentedMemory = []byte{
		0x11, 0x22, 0x33, 0x44,
		0x55, 0x66, 0x77, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
