package match

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/queue"
)

// Server ...
type Server interface {
	Init() Server
	Name() string
	Bid(bid.Server) Server
	Run() error
	Suspend() error
	Buffer() <-chan Result
}

// Result 撮合结果
type Result struct {
	Bid     bid.Server
	Amount  int
	Price   float64
	Trigger bid.UnitType
	Trades  []bid.UnitType
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

// Buffer ...
func (h *Match) Buffer() <-chan Result {
	return h.opts.buffer
}

// Stop ...
func (h *Match) Stop() error {
	h.bid.Init()
	return nil
}

func (h *Match) matchBuy(q queue.Server, n queue.NodeServer) error {
	unit := bid.NewUnit()
	if err := n.Content(unit); err != nil {
		return err
	}
	// unit, ok := n.Data().Content.(*bid.Unit)
	// if !ok {
	// 	return fmt.Errorf("buffer data type fail %v", n.Data().Content)
	// }

	result := Result{Bid: h.bid, Trigger: bid.UnitType{Type: bid.Type_Buy, Unit: *unit}}

	bizBid := h.bid.Sell()
	for n := bizBid.Front(); n != nil; n = n.Next() {
		content, ok := n.Data().Content.(*bid.Unit)
		if !ok {
			break
		}
		if unit.Price >= content.Price && result.Amount < unit.Amount {
			result.Price = content.Price
			if unit.Amount <= content.Amount {
				q.Remove(n)
				result.Amount += unit.Amount
				// 减去购买数量
				unit.Amount -= unit.Amount
				content.Amount -= unit.Amount
				if unit.Amount == content.Amount {
					bizBid.Remove(n)
				}

				result.Trades = append(result.Trades, bid.UnitType{Type: bid.Type_Sell, Unit: *content})
				break
			}
			if unit.Amount > content.Amount {
				bizBid.Remove(n)
				// 减去购买数量
				unit.Amount -= content.Amount
				// 加入到撮合成功数量中
				result.Amount += content.Amount
				result.Trades = append(result.Trades, bid.UnitType{Type: bid.Type_Sell, Unit: *content})
				continue
			}
		}
	}

	// 扣减当前交易结果
	if result.Amount < result.Trigger.Unit.Amount {
		n.Data().Update(unit)
	}
	// 成交后进行缓冲推送
	if result.Amount > 0 {
		h.opts.buffer <- result
	}
	return nil
}

func (h *Match) matchSell(q queue.Server, n queue.NodeServer) error {
	unit, ok := n.Data().Content.(*bid.Unit)
	if !ok {
		return fmt.Errorf("buffer data type fail %v", n.Data().Content)
	}

	result := Result{Bid: h.bid, Trigger: bid.UnitType{Unit: *unit}}

	bizBid := h.bid.Buy()
	for n := bizBid.Front(); n != nil; n = n.Next() {
		content := n.Data().Content.(*bid.Unit)
		if unit.Price <= content.Price {
			if unit.Amount == content.Amount {
				bizBid.Remove(n)
				q.Remove(n)
				result.Amount = unit.Amount
				result.Price = unit.Price
				result.Trades = append(result.Trades, bid.UnitType{Type: bid.Type_Buy, Unit: *content})
				// log.Println(q.Name(), unit.Price)
				break
			}
		}
	}
	// 扣减当前交易结果
	if result.Amount < result.Trigger.Unit.Amount {
		n.Data().Update(unit)
	}
	// 成交后进行缓冲推送
	if result.Amount > 0 {
		h.opts.buffer <- result
	}
	return nil
}

// Run ...
func (h *Match) Run() error {
	// 处理撮合
	go func() {
		for {
			select {
			case message := <-h.bid.Buffer():

				switch q := message.Queue; q {
				case h.bid.Buy():
					if err := h.matchBuy(q, message.Node); err != nil {
						break
					}
				case h.bid.Sell():
					h.matchBuy(q, message.Node)
				}
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

func (h *Match) buy() {

}

func (h *Match) sell() {

}
