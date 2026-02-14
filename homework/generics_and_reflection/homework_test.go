package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

// Никакой рекурсии или разыменования указателей не поддерживается, равно как не поддерживается
// и сериализация примитивных типов
func Serialize[T any](obj T) string {
	const tagName string = "properties"
	const omitEmptyTag string = "omitempty"

	res := &strings.Builder{}

	objValue := reflect.ValueOf(obj)

	for field, val := range objValue.Fields() {
		tag, ok := field.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")

		if len(parts) > 1 && parts[1] == omitEmptyTag && val.IsZero() {
			continue
		}

		strVal := fmt.Sprint(val)
		strLen := len(parts[0]) + len(strVal) + 2

		res.Grow(strLen)

		res.WriteString(parts[0])
		res.WriteString("=")
		res.WriteString(strVal)
		res.WriteString("\n")
	}

	return res.String()[:res.Len()-1]
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestGenericSerialization(t *testing.T) {
	type Struct1 struct {
		F1 string   `properties:"f1"`
		F2 int      `properties:"f2,omitempty"`
		F3 []uint64 `properties:"f3"`
		F4 *bool    `properties:"f4,omitempty"`
	}

	struct1 := Struct1{
		F2: 0,
		F3: []uint64{1, 2, 3},
	}
	expected1 := "f1=\nf3=[1 2 3]"

	actual1 := Serialize(struct1)
	assert.Equal(t, actual1, expected1)

	type Struct2 struct {
		F1 string  `properties:"f1,omitempty"`
		F2 Struct1 `properties:"f2,omitempty"`
	}

	struct2 := Struct2{
		F1: "kek",
	}
	expected2 := "f1=kek"

	actual2 := Serialize(struct2)
	assert.Equal(t, actual2, expected2)

	struct2 = Struct2{
		F1: "lol",
		F2: Struct1{},
	}
	expected2 = "f1=lol"

	actual2 = Serialize(struct2)
	assert.Equal(t, actual2, expected2)

	struct2 = Struct2{
		F2: Struct1{
			F2: 42,
		},
	}
	expected2 = "f2={ 42 [] <nil>}"

	actual2 = Serialize(struct2)
	fmt.Println(actual2)
	assert.Equal(t, actual2, expected2)
}
