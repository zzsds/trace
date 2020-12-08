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

// Result 撮合结果
type Result struct {
	Bid     bid.Server
	Amount  int
	Price   float64
	Trigger bid.Unit
	Trades  []bid.Unit
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

// Handle ...
func (h *Match) handle() error {
	go func() {
		for {
			select {
			case buf := <-h.bid.Buffer():
				result := Result{Bid: h.bid}
				msg := buf.(*bid.BufferMessage)
				q := msg.Queue
				if q == h.bid.Buy() {
					sell := h.bid.Sell()
					for n := sell.Front(); n != nil; n = n.Next() {
						content := n.Data.Content.(*bid.Unit)
						unit := msg.Data.Content.(*bid.Unit)
						if unit.Price >= content.Price {
							if unit.Amount == content.Amount {
								sell.Remove(n)
								q.Remove(msg.Node)
								log.Println(q.Name(), unit.Price, "--------------")
								result.Amount = unit.Amount
								result.Price = unit.Price
								result.Trigger = *unit
								result.Trades = append(result.Trades, *content)
								break
							}
						}
					}
				} else if q == h.bid.Sell() {
					buy := h.bid.Buy()
					for n := h.bid.Buy().Front(); n != nil; n = n.Next() {
						content := n.Data.Content.(*bid.Unit)
						unit := msg.Data.Content.(*bid.Unit)
						if unit.Price <= content.Price {
							if unit.Amount == content.Amount {
								buy.Remove(n)
								q.Remove(msg.Node)
								result.Amount = unit.Amount
								result.Price = unit.Price
								result.Trigger = *unit
								result.Trades = append(result.Trades, *content)
								log.Println(msg.Queue.Name(), unit.Price)
								break
							}
						}
					}
				}

				// 成交后进行缓冲推送
				if result.Amount > 0 {
					h.opts.buffer <- &result
				}
			}
		}
	}()
	return nil
}

// Run ...
func (h *Match) Run() error {

	// 处理撮合
	h.handle()

	ch := make(chan os.Signal, 1)
	if h.opts.signal {
		signal.Notify(ch, os.Kill)
	}

	// wait on kill signal
	<-ch
	return h.Stop()
}
