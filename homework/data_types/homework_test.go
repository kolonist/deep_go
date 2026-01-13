package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian(number uint32) uint32 {
	//           octets:   3 2 1 0
	const mask uint32 = 0x000000ff

	return (number&mask)<<24 | // octet 0 to 3
		(number>>8&mask)<<16 | // octet 1 to 2
		(number>>16&mask)<<8 | // octet 2 to 1
		(number >> 24 & mask) // octet 3 to 0
}

func ToLittleEndianGeneric[T uint16 | uint32 | uint64](number T) T {
	var mask T = 0xff
	sizeBytes := uint(unsafe.Sizeof(number))

	var result T

	// the same as ToLittleEndian but do it byte by byte
	// as we don't know how many bytes are in number
	for byteIdx := range sizeBytes {
		rShift := byteIdx * 8
		lShift := (sizeBytes - byteIdx - 1) * 8

		result |= (number >> rShift & mask) << lShift
	}

	return result
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #6": {
			number: 0x12345678,
			result: 0x78563412,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestGenericСonversion(t *testing.T) {
	tests16 := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case #4": {
			number: 0x0102,
			result: 0x0201,
		},
		"test case #5": {
			number: 0x1234,
			result: 0x3412,
		},
	}

	tests32 := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #6": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #7": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #8": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #9": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #10": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #11": {
			number: 0x12345678,
			result: 0x78563412,
		},
	}

	tests64 := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #12": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #13": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #14": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case #15": {
			number: 0x00000000FFFFFFFF,
			result: 0xFFFFFFFF00000000,
		},
		"test case #16": {
			number: 0x0000FFFF0000FFFF,
			result: 0xFFFF0000FFFF0000,
		},
		"test case #17": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
		"test case #18": {
			number: 0x0123456789ABCDEF,
			result: 0xEFCDAB8967452301,
		},
	}

	for name, test := range tests16 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	for name, test := range tests32 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	for name, test := range tests64 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
