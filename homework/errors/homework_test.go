package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Unwrap() []error {
	return e.errs
}

func (e *MultiError) Error() string {
	b := &strings.Builder{}
	b.Grow(len(e.errs) + 2)

	fmt.Fprintf(b, "%d errors occured:\n", len(e.errs))

	for _, err := range e.errs {

		fmt.Fprintf(b, "\t* %v", err)
	}

	b.WriteString("\n")

	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	multiErr, ok := err.(*MultiError)
	if !ok {
		return Append(&MultiError{}, errs...)
	}

	multiErr.errs = append(multiErr.errs, errs...)
	return multiErr
}

func TestMultiError(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	var err error
	err = Append(err, err1)
	err = Append(err, err2)

	assert.True(t, errors.Is(err, err1))
	assert.True(t, errors.Is(err, err2))
	assert.False(t, errors.Is(err, err3))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
