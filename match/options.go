package match

import (
	"context"
	"sync"
)

// options ...
type options struct {
	name   string
	mutex  *sync.RWMutex
	ctx    context.Context
	cancel func()
	state  bool
	buffer chan Result
	signal bool
	exit   chan bool
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		ctx:    context.Background(),
		mutex:  &sync.RWMutex{},
		buffer: make(chan Result, 100),
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

// WithContext ...
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}
