package match

import (
	"context"
)

// Options contains configuration for the Store
type Options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context
	Signal  bool
	Name    string
}

// Option sets values in Options
type Option func(o *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
		Signal:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
