package match

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/zzsds/trade/bid"
)

// Server ...
type Server interface {
	Init() Server
	Name() string
	Register(bid.Server) Server
	Bid() bid.Server
	State() bool
	Start() error
	Stop() error
	// Run 启动
	Run() error
	Buffer() <-chan Result
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

// Register ...
func (h *Match) Register(p bid.Server) Server {
	h.bid = p
	return h
}

// Bid ...
func (h *Match) Bid() bid.Server {
	return h.bid
}

// Init 初始化队列
func (h *Match) Init() Server {
	return h
}

// Name ..
func (h *Match) Name() string {
	return h.opts.name
}

// State State
func (h *Match) State() bool {
	return h.opts.state
}

// Start ...
func (h *Match) Start() error {
	var ctx context.Context
	if h.opts.state {
		return fmt.Errorf("Cannot start repeatedly")
	}
	ctx, h.opts.cancel = context.WithCancel(h.opts.ctx)
	h.opts.state = true
	return h.handle(ctx)
}

// Stop ...
func (h *Match) Stop() error {
	if !h.opts.state || h.opts.cancel == nil {
		return fmt.Errorf("Cannot stop repeatedly")
	}
	h.opts.cancel()
	return nil
}

// Buffer ...
func (h *Match) Buffer() <-chan Result {
	return h.opts.buffer
}

// handle 处理委托队列执行撮合交易
func (h *Match) handle(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				h.opts.state = false
				log.Println("goroutine exit")
				return
			case message := <-h.bid.Buffer():
				_ = message
				// fmt.Println(message.(*bid.Node))
				err := h.match(message.(*bid.Node))
				if err != nil {
					break
				}

			}
		}
	}()

	return ctx.Err()
}

// 撮合买卖委托交易
func (h *Match) match(node *bid.Node) error {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	if !h.opts.state {
		return nil
	}
	currentUnit, ok := node.Value.(bid.Unit)
	if !ok {
		return fmt.Errorf("Parsing failed")
	}

	result := Result{Bid: h.bid, Trigger: currentUnit}
	current, object := h.bid.Buy(), h.bid.Sell()
	if currentUnit.Type == bid.Type_Sell {
		current, object = h.bid.Sell(), h.bid.Buy()
	}

	if object.Len() <= 0 {
		return nil
	}

	for n := object.Front(); n != nil && currentUnit.Amount > result.Amount; n = n.Next() {
		objectUnit, ok := n.Value.(bid.Unit)
		if !ok {
			break
		}
		if objectUnit.Amount <= 0 {
			object.Remove(n)
			break
		}

		if currentUnit.Price >= objectUnit.Price {
			// 数量相等 全部匹配
			if currentUnit.Amount == objectUnit.Amount {
				current.Remove(node)
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				result.Trades = append(result.Trades, objectUnit)
				break
			}

			if currentUnit.Amount < objectUnit.Amount {
				current.Remove(node)
				objectUnit.Amount -= currentUnit.Amount
				n.Value = objectUnit
				// 加入到撮合成功数量中
				result.Amount += currentUnit.Amount
				result.Trades = append(result.Trades, objectUnit)
				break
			}

			if currentUnit.Amount > objectUnit.Amount {
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				// 减去购买数量
				currentUnit.Amount -= objectUnit.Amount
				n.Value = currentUnit
				result.Trades = append(result.Trades, objectUnit)
				continue
			}
		}
	}
	// 处理当前交易请求结果
	if result.Amount < result.Trigger.Amount {
		node.Value = currentUnit
	}
	// 成交后进行缓冲推送
	if result.Amount > 0 {
		h.opts.buffer <- result
	}
	return nil
}

// Run ...
func (h *Match) Run() error {
	if h.opts.state {
		return fmt.Errorf("runing")
	}
	// 开始撮合
	if err := h.Start(); err != nil {
		log.Fatalf("Run Match Fail：%v", err)
	}

	ch := make(chan os.Signal, 1)
	if h.opts.signal {
		signal.Notify(ch, os.Kill)
	}

	// wait on kill signal
	<-ch
	return h.Stop()
}
