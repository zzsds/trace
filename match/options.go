package match

import (
	"context"
	"sync"
)

// Options contains configuration for the Store
type Options struct {
	// Context should contain all implementation specific Options, using context.WithValue.
	Context context.Context
	Cancel  context.CancelFunc
	State   bool
	Signal  bool
	Name    string
	mu      *sync.Mutex
}

// Option sets values in Options
type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
		Signal:  true,
		mu:      &sync.Mutex{},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name ...
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// WithContext ...
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
