package list

import (
	"context"
	"sync"
)

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

// Handler ...
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware ...
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
