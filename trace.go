package trade

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/match"
)

var (
	// errMissingName is returned by service.Run when a service is run
	// prior to it's name being set.
	errMissingName = errors.New("missing service name")
)

// Trade 交易
type Trade struct {
	opts Options
	bids map[int]bid.Server
}

// Newtrade 初始化
func Newtrade(opts ...Option) *Trade {
	return &Trade{
		opts: newOptions(opts...),
		bids: make(map[int]bid.Server),
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

// Init ...
func (s *Trade) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}
}

// Options ...
func (s *Trade) Options() Options {
	return s.opts
}

func (s *Trade) String() string {
	return "trade"
}

// Start ...
func (s *Trade) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	for _, v := range s.bids {
		m := match.NewMatch()
		m.Bid(v)
		go m.Run()
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// Stop ...
func (s *Trade) Stop() error {
	var err error

	for _, fn := range s.opts.BeforeStop {
		if e := fn(); e != nil {
			err = e
		}
	}

	for _, v := range s.bids {
		match := match.NewMatch()
		match.Bid(v)
		go match.Run()
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

// RegisterBid a bid
func (s *Trade) RegisterBid(id int, p bid.Server) error {
	if _, ok := s.bids[id]; ok {
		return fmt.Errorf("bid %d already exists", id)
	}
	s.bids[id] = p
	return nil
}

// LoadBid a bid
func (s *Trade) LoadBid(id int) (bid.Server, error) {
	v, ok := s.bids[id]
	if !ok {
		return nil, fmt.Errorf("bid %d does not exist", id)
	}
	return v, nil
}
