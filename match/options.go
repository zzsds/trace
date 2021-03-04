package match

import (
	"context"
	"sync"
)

// Options contains configuration for the Store
type options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	context context.Context
	cancel  context.CancelFunc
	state   bool
	signal  bool
	name    string
	mu      *sync.Mutex
}

// Option sets values in Options
type Option func(o *options)

func newOptions(opts ...Option) options {
	opt := options{
		context: context.Background(),
		signal:  true,
		mu:      &sync.Mutex{},
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
		o.context = ctx
	}
}
