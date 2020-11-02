package queue

import "sync"

// Options ...
type Options struct {
	Name  string
	mutex *sync.RWMutex
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
