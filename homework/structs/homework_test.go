package main

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func IsBigEndian() bool {
	return *(*uint16)(unsafe.Pointer(&[2]byte{0xAB, 0xCD})) == 0xCDAB
}

var _isBigEndian bool = IsBigEndian()

func FixEndianess(value uint32) uint32 {
	if _isBigEndian {
		return value>>24 | value<<24 | (value << 8 & 0x00FF0000) | (value >> 8 & 0x0000FF00)
	}
	return value
}

func PackString(s string, dest *byte, destLen int) {
	i := 0
	for runeIdx, char := range s {
		bitIdx := runeIdx * 7
		byteIdx := bitIdx / 8
		shift := bitIdx % 8

		// get 4 bytes from destination for better experience with rune bit operations
		destPtr := (*uint32)(unsafe.Add(unsafe.Pointer(dest), byteIdx))

		// shift 25 bits left to place last 7 bits of rune at the start of the destination
		// then shift `shift` bits right to place start of the 7-bit symbol to correct place
		*destPtr |= FixEndianess(uint32(char) << (25 - shift))

		i++
		if i >= destLen {
			break
		}
	}
}

func UnpackToString(src *byte, srcLen int) string {
	destMaxLen := srcLen * 8 / 7

	destBuilder := strings.Builder{}
	destBuilder.Grow(destMaxLen)

	for charIdx := range destMaxLen {
		bitIdx := charIdx * 7
		byteIdx := bitIdx / 8
		shift := bitIdx % 8

		// get 4 bytes from destination for better experience with rune bit operations
		srcPtr := (*uint32)(unsafe.Add(unsafe.Pointer(src), byteIdx))

		// shift `shift` bits left to place start of the 7-bit symbol to the start of the 32-bits structure
		// then shift 25 bits right to place 7 bits to the end of the destination char
		// then remove part of preceeding symbols
		char := rune((FixEndianess(*srcPtr) >> (25 - shift)) & 127)

		if char == 0 {
			break
		}

		destBuilder.WriteRune(char)
	}

	return destBuilder.String()
}

func WriteInt32ByByteOffset(value int32, buf *byte, offset uint8) {
	destPtr := (*int32)(unsafe.Add(unsafe.Pointer(buf), offset))
	*destPtr = value
}

func ReadInt32ByByteOffset(buf *byte, offset uint8) int32 {
	return *(*int32)(unsafe.Add(unsafe.Pointer(buf), offset))
}

func WriteBitsByOffset(value uint32, bits uint8, buf *byte, bytesOffset uint8, bitsOffset uint8) {
	destPtr := (*uint32)(unsafe.Add(unsafe.Pointer(buf), bytesOffset))

	// move bits we need to the start of 32-bits struct
	*destPtr |= FixEndianess(value << (32 - bits - bitsOffset))
}

func ReadBitsByOffset(bits uint8, buf *byte, bytesOffset uint8, bitsOffset uint8) uint32 {
	data := *(*uint32)(unsafe.Add(unsafe.Pointer(buf), bytesOffset))

	var mask uint32 = 0xFFFFFFFF >> (32 - bits)

	// move bits we need to the end of 32-bits struct
	return (FixEndianess(data) >> (32 - bits - bitsOffset) & mask)
}

func WriteFlagByOffset(value bool, buf *byte, bytesOffset uint8, bitsOffset uint8) {
	var flag uint32
	if value {
		flag = 1
	}

	WriteBitsByOffset(flag, 1, buf, bytesOffset, bitsOffset)
}

func ReadFlagByOffset(buf *byte, bytesOffset uint8, bitsOffset uint8) bool {
	value := ReadBitsByOffset(1, buf, bytesOffset, bitsOffset)
	return value == 1
}

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		PackString(name, (*byte)(unsafe.Pointer(&person.data)), 42)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteInt32ByByteOffset(int32(x), (*byte)(unsafe.Pointer(&person.data)), 37)
		WriteInt32ByByteOffset(int32(y), (*byte)(unsafe.Pointer(&person.data)), 41)
		WriteInt32ByByteOffset(int32(z), (*byte)(unsafe.Pointer(&person.data)), 45)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(gold), 31, (*byte)(unsafe.Pointer(&person.data)), 49, 0)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(mana), 10, (*byte)(unsafe.Pointer(&person.data)), 52, 7)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(health), 10, (*byte)(unsafe.Pointer(&person.data)), 54, 1)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(respect), 4, (*byte)(unsafe.Pointer(&person.data)), 54, 11)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(strength), 4, (*byte)(unsafe.Pointer(&person.data)), 54, 15)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(experience), 4, (*byte)(unsafe.Pointer(&person.data)), 54, 19)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(level), 4, (*byte)(unsafe.Pointer(&person.data)), 54, 23)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		WriteFlagByOffset(true, (*byte)(unsafe.Pointer(&person.data)), 54, 27)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		WriteFlagByOffset(true, (*byte)(unsafe.Pointer(&person.data)), 54, 28)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		WriteFlagByOffset(true, (*byte)(unsafe.Pointer(&person.data)), 54, 29)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		WriteBitsByOffset(uint32(personType), 2, (*byte)(unsafe.Pointer(&person.data)), 54, 30)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

/*
0       1       2        3      4       5       6       7
1000000010000000100000001000000010000000100000001000000010000000
[name___________________________________________________________
8       9       10      11      12      13      14      15
1000000010000000100000001000000010000000100000001000000010000000
16      17      18      19      20      21      22      23
1000000010000000100000001000000010000000100000001000000010000000
24      25      26      27      28      29      30      31
1000000010000000100000001000000010000000100000001000000010000000
32      33      34      35      36      37      37      39
1000000010000000100000001000000010000000100000001000000010000000
_________________________________name]**[x______________________
40      41      42      43      44      45      46      47
1000000010000000100000001000000010000000100000001000000010000000
_______][y____________________________y][z______________________
48      49      50      51      52      53      54      55
1000000010000000100000001000000010000000100000001000000010000000
_______][gold_____________________gold][mana____][health__][  ][  ->  [respect][strength
56      57
1000000010000000
  ][  ][  ][ ][]                                                  ->  strength][experience][level][flags: hasGun, hasFamilty, hasHouse][personType]
*/

type GamePerson struct {
	data [58]byte
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}

	for _, option := range options {
		option(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	return UnpackToString((*byte)(unsafe.Pointer(&p.data)), 42)
}

func (p *GamePerson) X() int {
	return int(ReadInt32ByByteOffset((*byte)(unsafe.Pointer(&p.data)), 37))
}

func (p *GamePerson) Y() int {
	return int(ReadInt32ByByteOffset((*byte)(unsafe.Pointer(&p.data)), 41))
}

func (p *GamePerson) Z() int {
	return int(ReadInt32ByByteOffset((*byte)(unsafe.Pointer(&p.data)), 45))
}

func (p *GamePerson) Gold() int {
	return int(ReadBitsByOffset(31, (*byte)(unsafe.Pointer(&p.data)), 49, 0))
}

func (p *GamePerson) Mana() int {
	return int(ReadBitsByOffset(10, (*byte)(unsafe.Pointer(&p.data)), 52, 7))
}

func (p *GamePerson) Health() int {
	return int(ReadBitsByOffset(10, (*byte)(unsafe.Pointer(&p.data)), 54, 1))
}

func (p *GamePerson) Respect() int {
	return int(ReadBitsByOffset(4, (*byte)(unsafe.Pointer(&p.data)), 54, 11))
}

func (p *GamePerson) Strength() int {
	return int(ReadBitsByOffset(4, (*byte)(unsafe.Pointer(&p.data)), 54, 15))
}

func (p *GamePerson) Experience() int {
	return int(ReadBitsByOffset(4, (*byte)(unsafe.Pointer(&p.data)), 54, 19))
}

func (p *GamePerson) Level() int {
	return int(ReadBitsByOffset(4, (*byte)(unsafe.Pointer(&p.data)), 54, 23))
}

func (p *GamePerson) HasGun() bool {
	return ReadFlagByOffset((*byte)(unsafe.Pointer(&p.data)), 54, 27)
}

func (p *GamePerson) HasFamilty() bool {
	return ReadFlagByOffset((*byte)(unsafe.Pointer(&p.data)), 54, 28)
}

func (p *GamePerson) HasHouse() bool {
	return ReadFlagByOffset((*byte)(unsafe.Pointer(&p.data)), 54, 29)
}

func (p *GamePerson) Type() int {
	return int(ReadBitsByOffset(2, (*byte)(unsafe.Pointer(&p.data)), 54, 30))
}

func TestGamePerson(t *testing.T) {
	fmt.Printf("Size of struct: %d\n", unsafe.Sizeof(GamePerson{}))

	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

func TestPackString(t *testing.T) {
	tests := map[string]struct {
		s      string
		maxLen int
	}{
		"test case #1": {
			s:      "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc",
			maxLen: 42,
		},
		"test case #2": {
			s:      "0123456789",
			maxLen: 42,
		},
		"test case #3": {
			s:      "a1b2c3",
			maxLen: 8,
		},
		"test case #4": {
			s:      "ab",
			maxLen: 2,
		},
		"test case #5": {
			s:      "abc",
			maxLen: 7,
		},
		"test case #6": {
			s:      "abcd",
			maxLen: 4,
		},
		"test case #7": {
			s:      "abcdefgh",
			maxLen: 4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			packed := make([]byte, test.maxLen)

			packedPtr := unsafe.SliceData(packed)

			PackString(test.s, packedPtr, test.maxLen)
			result := UnpackToString(packedPtr, test.maxLen)

			if len(test.s) > test.maxLen {
				test.s = test.s[:test.maxLen]
			}

			assert.Equal(t, test.s, result)
		})
	}
}

func TestFixEndianess(t *testing.T) {
	tests := map[string]struct {
		src      uint32
		expected uint32
	}{
		"test case #1": {
			src:      0x00000000,
			expected: 0x00000000,
		},
		"test case #2": {
			src:      0x1A2B3C4D,
			expected: 0x4D3C2B1A,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FixEndianess(test.src)
			assert.Equal(t, test.expected, actual)
		})
	}
}
