package trade

import (
	"errors"
	"os"
	"os/signal"

	"github.com/zzsds/trade/match"
)

var (
	// errMissingName is returned by service.Run when a service is run
	// prior to it's name being set.
	errMissingName = errors.New("missing service name")
)

// Server ...
type Server interface {
	Name() string
	Version() string
	String() string
	Register(match.Server) Server
	Load(string) (match.Server, error)
	Run() error
}

// Trade 交易
type Trade struct {
	opts  Options
	match map[string]match.Server
}

// Newtrade 初始化
func Newtrade(opts ...Option) Server {
	return &Trade{
		opts:  newOptions(opts...),
		match: make(map[string]match.Server),
	}
}

// Name of the service
func (s *Trade) Name() string {
	return s.opts.Name
}

// Version of the service
func (s *Trade) Version() string {
	return s.opts.Version
}

// Options ...
func (s *Trade) Options() Options {
	return s.opts
}

func (s *Trade) String() string {
	return "trade"
}

// Load ...
func (s *Trade) Load(name string) (match.Server, error) {
	m, ok := s.match[name]
	if !ok {
		return nil, errors.New("Trade match non-existent")
	}
	return m, nil
}

// Register ...
func (s *Trade) Register(match match.Server) Server {
	if _, ok := s.match[match.Name()]; !ok {
		match.Start()
		s.match[match.Name()] = match
	}
	return s
}

// start ...
func (s *Trade) start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}
	for _, m := range s.match {
		if err := m.Start(); err != nil {
			return err
		}
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// stop ...
func (s *Trade) stop() error {
	var err error

	for _, fn := range s.opts.BeforeStop {
		if e := fn(); e != nil {
			err = e
		}
	}

	for _, fn := range s.opts.AfterStop {
		if e := fn(); e != nil {
			err = e
		}
	}

	return err
}

// Run the service
func (s *Trade) Run() error {
	// ensure service's have a name, this is injected by the runtime manager
	if len(s.Name()) == 0 {
		return errMissingName
	}

	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, os.Kill)
	}

	// wait on kill signal
	<-ch
	return s.stop()
}
