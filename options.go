package trace

// Options options
type Options struct {
	Name    string
	Version string
	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Signal bool
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:   "Trace",
		Signal: true,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Option ...
type Option func(o *Options)

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

// Name ...
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Version ...
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}
