package list

import "sync"

// Options contains configuration for the Store
type Options struct {
	*sync.Mutex
	Name string
}

// Option sets values in Options
type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{Mutex: &sync.Mutex{}}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
