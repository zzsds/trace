package bid

import (
	"sync"
)

// options ...
type options struct {
	mutex  *sync.RWMutex
	id     int
	name   string
	buffer chan Message
	signal bool
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		mutex:  &sync.RWMutex{},
		buffer: make(chan Message, 1000),
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
