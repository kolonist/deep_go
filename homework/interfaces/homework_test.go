package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Constructor = func() any

type TypeDescriptor struct {
	constructor Constructor
	singleton   bool
	impl        any
}

type Container struct {
	types map[string]TypeDescriptor
}

func NewContainer() *Container {
	return &Container{
		types: make(map[string]TypeDescriptor),
	}
}

func (c *Container) RegisterType(name string, constructor Constructor) {
	c.types[name] = TypeDescriptor{
		constructor: constructor,
	}
}

func (c *Container) RegisterSingleton(name string, constructor Constructor) {
	c.types[name] = TypeDescriptor{
		constructor: constructor,
		singleton:   true,
	}
}

func (c *Container) Resolve(name string) (any, error) {
	t, ok := c.types[name]
	if !ok {
		return nil, fmt.Errorf("Can't find type %s", name)
	}

	if !t.singleton {
		return t.constructor(), nil
	}

	if t.impl == nil {
		t.impl = t.constructor()
		c.types[name] = t
	}

	return t.impl, nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() any {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() any {
		return &MessageService{}
	})
	container.RegisterSingleton("SingletonService", func() any {
		return &UserService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)
	assert.False(t, userService1 == userService2)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	singletonService1, err := container.Resolve("SingletonService")
	assert.NoError(t, err)
	singletonService2, err := container.Resolve("SingletonService")
	assert.NoError(t, err)
	assert.True(t, singletonService1 == singletonService2)

	s1 := singletonService1.(*UserService)
	s2 := singletonService2.(*UserService)
	assert.True(t, s1 == s2)
}
