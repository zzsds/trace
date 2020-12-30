package queue

import (
	"sync"
)

// options ...
type options struct {
	mutex *sync.RWMutex
	name  string
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		mutex: &sync.RWMutex{},
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
