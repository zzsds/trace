package trade

import (
	"errors"
	"os"
	"os/signal"
)

var (
	// errMissingName is returned by service.Run when a service is run
	// prior to it's name being set.
	errMissingName = errors.New("missing service name")
)

// trade 交易
type Service struct {
	opts Options
}

// Newtrade 初始化
func Newtrade(opts ...Option) *Service {
	return &Service{
		opts: newOptions(opts...),
	}
}

// Name of the service
func (s *Service) Name() string {
	return s.opts.Name
}

// Version of the service
func (s *Service) Version() string {
	return s.opts.Version
}

// Init ...
func (s *Service) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}
}

// Options ...
func (s *Service) Options() Options {
	return s.opts
}

func (s *Service) String() string {
	return "trade"
}

// Start ...
func (s *Service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
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

// Stop ...
func (s *Service) Stop() error {
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
func (s *Service) Run() error {
	// ensure service's have a name, this is injected by the runtime manager
	if len(s.Name()) == 0 {
		return errMissingName
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, os.Kill)
	}

	// wait on kill signal
	<-ch
	return s.Stop()
}
