package bid

import (
	"fmt"

	"github.com/zzsds/trade/queue"
)

// Server ...
type Server interface {
	ID() int
	Name() string
	Init() Server
	Buy() queue.Server
	Sell() queue.Server
	Cancel(queue.Server, int) error
	Add(*Unit) (Unit, error)
	Buffer() <-chan Message
}

// Message 缓冲消息
type Message struct {
	Queue queue.Server
	*queue.Node
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

// Init 初始化交易对
func (h *Bid) Init() Server {
	h.buy = queue.NewQueue(queue.Name(Type_Buy.String()))
	h.sell = queue.NewQueue(queue.Name(Type_Sell.String()))
	return h
}

// ID ...
func (h *Bid) ID() int {
	return h.opts.id
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
func (h *Bid) Add(u *Unit) (Unit, error) {
	q := h.buy
	if u.Type == Type_Sell {
		q = h.sell
	}
	message := Message{Queue: q}
	for n := q.Front(); ; n = n.Next() {
		if n == nil {
			message.Node = q.PushBack(u)
			break
		}
		v := n.Value.(*Unit)
		if v.Price == u.Price {
			if v.UID == u.UID {
				v.Amount += u.Amount
				n.Value = v
				message.Node = n
				break
			}
			if v.CreateAt.After(u.CreateAt) {
				message.Node = q.InsertBefore(u, n)
				break
			}
		}

		//如果是买家队列，按照价格高优先，时间优先
		if h.buy == q {
			//价格高者优先
			if u.Price > v.Price {
				message.Node = q.InsertBefore(u, n)
				break
			}
		} else if h.sell == q {
			//价格高者优先
			if v.Price > u.Price {
				message.Node = q.InsertBefore(u, n)
				break
			}
		}
	}
	h.opts.amount++
	h.opts.buffer <- message
	return *u, nil
}

// Cancel ...
func (h *Bid) Cancel(q queue.Server, ID int) error {
	q.Loop(func(n *queue.Node) error {
		if n != nil && n.Value.(*Unit).ID == ID {
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
