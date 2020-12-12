package match

import (
	"os"
	"os/signal"
	"time"

	"github.com/zzsds/trade/bid"
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

// 撮合卖
func (h *Match) match(message bid.Message) error {
	node := message.Node
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
	// 处理撮合
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
			case message := <-h.bid.Buffer():
				if err := h.match(message); err != nil {
					break
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
