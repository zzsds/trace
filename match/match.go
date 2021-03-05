package match

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/list"
)

// Server ...
type Server interface {
	Init(...Option)
	Options() Options
	Name() string
	Start() error
	Stop() error
	State() bool
	Run() error
	Close() error
	String() string
	Bid() bid.Server
	Queue() <-chan interface{}
}

// Result 撮合结果
type Result struct {
	Bid     bid.Server
	Amount  int
	Price   float64
	Trigger *bid.Unit
	Trades  []*bid.Unit
}

type match struct {
	opts  Options
	once  sync.Once
	bid   bid.Server
	queue chan interface{}
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
		m.bid = bid.NewBid(bid.Name(m.opts.Name))
		m.queue = make(chan interface{})
	})
}

// Queue ...
func (m *match) Queue() <-chan interface{} {
	return m.queue
}

//  Bid ...
func (m *match) Bid() bid.Server {
	return m.bid
}

// Options ...
func (m *match) Options() Options {
	return m.opts
}

// Close ...
func (m *match) Close() error {
	m.bid.Init(bid.Name(m.opts.Name))
	return m.Stop()
}

// Run ...
func (m *match) Start() error {
	if m.State() {
		return errors.New("match started")
	}
	var ctx context.Context
	ctx, m.opts.Cancel = context.WithCancel(context.Background())
	m.setState(true)
	return m.handle(ctx)
}

func (m *match) Stop() error {
	if !m.State() {
		return errors.New("match stopped")
	}
	var err error
	m.opts.Cancel()
	m.setState(false)
	return err
}

func (m *match) setState(b bool) {
	m.opts.mu.Lock()
	defer m.opts.mu.Unlock()
	m.opts.State = b
}

func (m *match) State() bool {
	m.opts.mu.Lock()
	defer m.opts.mu.Unlock()
	return m.opts.State
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

	return m.Close()
}

func (m *match) handle(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("goroutine exit")
				return
			case queue := <-m.bid.Queue():
				if err := m.match(queue); err != nil {
					break
				}
			}
		}
	}()
	return ctx.Err()
}

func (m *match) match(v interface{}) error {
	node, ok := v.(*list.Node)
	if !ok {
		return errors.New("Queue Parsing failed")
	}

	currentUnit, ok := node.Value.(*bid.Unit)
	if !ok {
		return errors.New("Unit Parsing failed")
	}

	current, object := m.bid.Buy(), m.bid.Sell()
	if currentUnit.Type == bid.Type_Sell {
		current, object = m.bid.Sell(), m.bid.Buy()
	}
	current.Lock()
	defer current.Unlock()

	result := Result{Bid: m.bid, Trigger: currentUnit}

	object.Lock()
	defer object.Unlock()
	// 撮合匹配
	for n := object.Front(); n != nil && currentUnit.Amount > result.Amount; n = n.Next() {
		objectUnit, ok := n.Value.(*bid.Unit)
		if !ok {
			log.Fatalln(fmt.Errorf("Parsing failed"))
			break
		}
		if objectUnit.Amount <= 0 {
			break
		}

		if currentUnit.Price >= objectUnit.Price {
			// 数量相等 全部匹配
			if currentUnit.Amount == objectUnit.Amount {
				current.Remove(node)
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				result.Price = currentUnit.Price
				result.Trades = append(result.Trades, objectUnit)
				break
			}

			if currentUnit.Amount < objectUnit.Amount {
				current.Remove(node)
				objectUnit.Amount -= currentUnit.Amount
				n.Value = objectUnit
				// 加入到撮合成功数量中
				result.Amount += currentUnit.Amount
				result.Price = currentUnit.Price
				result.Trades = append(result.Trades, objectUnit)
				break
			}

			if currentUnit.Amount > objectUnit.Amount {
				object.Remove(n)
				// 加入到撮合成功数量中
				result.Amount += objectUnit.Amount
				result.Price = currentUnit.Price
				// 减去购买数量
				currentUnit.Amount -= objectUnit.Amount
				n.Value = currentUnit
				result.Trades = append(result.Trades, objectUnit)
			}
		}
	}
	// 处理当前交易请求结果
	if result.Amount < result.Trigger.Amount {
		node.Value = currentUnit
	}
	// 成交后进行缓冲推送
	if result.Amount > 0 {
		m.queue <- result
	}

	return nil
}
