package bid

import (
	"fmt"

	"github.com/zzsds/trade/queue"
)

// Server ...
type Server interface {
	Init() Server
	Buy() queue.Server
	Sell() queue.Server
	Cancel(queue.Server, int) error
	Add(queue.Server, *Unit) (queue.Data, error)
	Buffer() <-chan Message
	ID() int
	Name() string
}

// Message 缓冲消息
type Message struct {
	Queue queue.Server
	Node  queue.NodeServer
}

// UnitType ...
type UnitType struct {
	Type
	Unit
}

// Unit ...
type Unit struct {
	ID     int
	UID    int
	Name   string
	Price  float64
	Amount int
}

// NewUnit ...
func NewUnit() *Unit {
	return new(Unit)
}

var bids = make(map[int]Server)

// Bid bid
type Bid struct {
	opts      options
	buy, sell queue.Server
}

// Register a bid
func Register(id int, p Server) error {
	if _, ok := bids[id]; ok {
		return fmt.Errorf("bid %d already exists", id)
	}
	bids[id] = p
	return nil
}

// Load a bid
func Load(id int) (Server, error) {
	v, ok := bids[id]
	if !ok {
		return nil, fmt.Errorf("bid %d does not exist", id)
	}
	return v, nil
}

// NewBid ...
func NewBid(opts ...Option) Server {
	bid := new(Bid)
	bid.Init()
	bid.opts = newOptions(opts...)
	bids[bid.opts.id] = bid
	return bid
}

// ID ...
func (h *Bid) ID() int {
	return h.opts.id
}

// Init 初始化交易对
func (h *Bid) Init() Server {
	h.buy = queue.NewQueue(queue.Name(Type_Buy.String()))
	h.sell = queue.NewQueue(queue.Name(Type_Sell.String()))
	return h
}

// Name ...
func (h *Bid) Name() string {
	return h.opts.name
}

// Buy ...
func (h *Bid) Buy() queue.Server {
	return h.buy
}

// Sell ...
func (h *Bid) Sell() queue.Server {
	return h.sell
}

// Amount ...
func (h *Bid) Amount() int {
	return h.opts.amount
}

// Add ...
func (h *Bid) Add(q queue.Server, u *Unit) (queue.Data, error) {
	message := Message{Queue: q}
	data := queue.NewData(u)
	for n := q.Front(); ; n = n.Next() {
		if n == nil {
			message.Node = q.Push(data)
			break
		}
		content := n.Data().Content.(*Unit)
		if content.Price == u.Price {
			if content.UID == u.UID {
				content.Amount += u.Amount
				n.Data().UpdateContent(content)
				message.Node = n
				break
			}
			if n.Data().CreateAt.After(data.CreateAt) {
				message.Node = q.InsertBefore(data, n)
				break
			}
		}

		//如果是买家队列，按照价格高优先，时间优先
		if h.buy == q {
			//价格高者优先
			if u.Price > content.Price {
				message.Node = q.InsertBefore(data, n)
				break
			}
		} else if h.sell == q {
			//价格高者优先
			if content.Price > u.Price {
				message.Node = q.InsertBefore(data, n)
				break
			}
		}
	}
	h.opts.amount++
	h.opts.buffer <- message
	return *data, nil
}

// Cancel ...
func (h *Bid) Cancel(q queue.Server, ID int) error {
	q.Loop(func(n *queue.Node) error {
		content := n.Data().Content
		if content != nil && content.(*Unit).ID == ID {
			q.Remove(n)
		}
		return nil
	})
	return nil
}

// Buffer 获取缓冲
func (h *Bid) Buffer() <-chan Message {
	return h.opts.buffer
}
