package match

import (
	"log"
	"os"
	"os/signal"

	"github.com/zzsds/trade/bid"
)

// Server ...
type Server interface {
	Init() Server
	Name() string
	Bid(bid.Server) Server
	Run() error
	Suspend() error
}

// Match Match
type Match struct {
	opts options
	bid  bid.Server
}

// NewMatch new match
func NewMatch(opts ...Option) Server {
	m := new(Match)
	m.Init()
	m.opts = newOptions(opts...)
	return m
}

// Bid ...
func (h *Match) Bid(p bid.Server) Server {
	h.bid = p
	return h
}

// Init 初始化队列
func (h *Match) Init() Server {
	return h
}

// Name ..
func (h *Match) Name() string {
	return h.opts.name
}

// Suspend ...
func (h *Match) Suspend() error {
	return nil
}

// Stop ...
func (h *Match) Stop() error {
	h.bid.Init()
	return nil
}

// Run ...
func (h *Match) Run() error {
	go func() {
		for {
			select {
			case buy := <-h.bid.Buy().Buffer():
				log.Println(buy)
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	if h.opts.signal {
		signal.Notify(ch, os.Kill)
	}

	// wait on kill signal
	<-ch
	return h.Stop()
}
