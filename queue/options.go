package queue

import "sync"

// options ...
type options struct {
	mutex  *sync.RWMutex
	name   string
	size   int
	buffer chan interface{}
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		mutex:  &sync.RWMutex{},
		buffer: make(chan interface{}),
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
