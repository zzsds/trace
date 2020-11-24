package queue

import "sync"

// Options ...
type Options struct {
	mutex *sync.RWMutex
	name  string
	size  int
}

// Option ...
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		mutex: &sync.RWMutex{},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name ...
func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}
