package match

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server ...
type Server interface {
	Init(...Option)
	Options() Options
	Name() string
	Start() error
	Run() error
	Close() error
	String() string
}

type match struct {
	opts Options
	once sync.Once
}

// NewMatch ...
func NewMatch(opts ...Option) Server {
	m := new(match)
	// set opts
	m.opts = newOptions(opts...)
	m.Init()

	return m
}

// String ...
func (m *match) String() string {
	return "match"
}

func (m *match) Name() string {
	return m.opts.Name
}

// Init ...
func (m *match) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&m.opts)
	}

	m.once.Do(func() {

	})
}

// Options ...
func (m *match) Options() Options {
	return m.opts
}

// Close ...
func (m *match) Close() error {
	return nil
}

// Run ...
func (m *match) Start() error {
	return nil
}

func (m *match) Stop() error {
	var err error

	return err
}

// Run ...
func (m *match) Run() error {

	if err := m.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if m.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-m.opts.Context.Done():
	}

	return m.Stop()
}
