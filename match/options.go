package match

import (
	"context"
	"sync"
)

// options ...
type options struct {
	ctx    context.Context
	mutex  *sync.RWMutex
	name   string
	buffer chan Result
	signal bool
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	ch := make(chan Result, 100)
	opt := options{
		ctx:    context.Background(),
		mutex:  &sync.RWMutex{},
		buffer: ch,
		signal: true,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name ...
func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
