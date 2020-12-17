package match

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// options ...
type options struct {
	ctx    context.Context
	mutex  *sync.RWMutex
	name   string
	buffer chan Result
	signal bool
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	ch := make(chan Result, 100)
	ctx, cancel := context.WithCancel(context.Background())

	opt := options{
		ctx:    ctx,
		mutex:  &sync.RWMutex{},
		buffer: ch,
		signal: true,
	}
	for _, o := range opts {
		o(&opt)
	}

	timer := time.NewTimer(5 * time.Second)
	go func() {
		select {
		case <-timer.C:
			cancel()
			fmt.Println("结束")
		}
	}()
	return opt
}

// Name ...
func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
