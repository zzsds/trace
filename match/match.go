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
	ctx, h.opts.cancel = context.WithCancel(h.opts.ctx)
	h.opts.state = true
	return h.handle(ctx)
}

// Stop ...
func (h *Match) Stop() error {
	if h.opts.cancel == nil {
		return fmt.Errorf("cancel undefind")
	}
	h.opts.state = false
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
				fmt.Println("退出协程")
				h.opts.state = false
				return
			case message := <-h.bid.Buffer():
				if err := h.match(ctx, message); err != nil {
					break
				}
			}
		}
	}()

	return ctx.Err()
}

// 撮合买卖委托交易
func (h *Match) match(ctx context.Context, message bid.Message) error {
	node := message.Node
	if node == nil {
		fmt.Println("停止")
	}
	currentUnit := bid.NewUnit()
	if err := node.Data().ParseContent(currentUnit); err != nil {
		return err
	}
	result := Result{Bid: h.bid, Trigger: bid.UnitType{Type: bid.Type_Buy, Unit: *currentUnit}}
	unitType := bid.UnitType{Type: bid.Type_Sell}

	current, object := h.bid.Buy(), h.bid.Sell()
	if message.Queue == h.bid.Sell() {
		current, object = h.bid.Sell(), h.bid.Buy()
		// 指定当前交易类型
		result.Trigger.Type = bid.Type_Sell
		// 反向指定交易对象
		object = h.bid.Buy()
		// 反向指定交易对象类型
		unitType.Type = bid.Type_Buy
	}

	for n := object.Front(); n != nil && currentUnit.Amount > result.Amount; n = n.Next() {
		objectUnit, ok := n.Data().Content.(*bid.Unit)
		if !ok {
			break
		}
		current.Remove(node.Current())
		unitType.Unit = *objectUnit
		if currentUnit.Price >= objectUnit.Price {
			// 数量相等 全部匹配
			if currentUnit.Amount == objectUnit.Amount {
				current.Remove(node.Current())
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				result.Trades = append(result.Trades, unitType)
				break
			}

			if currentUnit.Amount < objectUnit.Amount {
				current.Remove(node.Current())
				objectUnit.Amount -= currentUnit.Amount
				n.Data().UpdateContent(objectUnit)
				// 加入到撮合成功数量中
				result.Amount += currentUnit.Amount
				result.Trades = append(result.Trades, unitType)
				break
			}

			if currentUnit.Amount > objectUnit.Amount {
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				// 减去购买数量
				currentUnit.Amount -= objectUnit.Amount
				result.Trades = append(result.Trades, unitType)
				continue
			}
		}
	}
	// 扣减当前交易结果
	if result.Amount < result.Trigger.Unit.Amount {
		node.Data().UpdateContent(currentUnit)
	}
	// 成交后进行缓冲推送
	if result.Amount > 0 {
		h.opts.buffer <- result
	}
	return nil
}

// Run ...
func (h *Match) Run() error {
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
