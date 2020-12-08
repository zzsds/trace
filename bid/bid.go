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
	Buffer() <-chan interface{}
	ID() int
	Name() string
}

// BufferMessage 缓冲消息
type BufferMessage struct {
	Queue queue.Server
	*queue.Node
}

// Unit ...
type Unit struct {
	Name    string
	Price   float64
	Amount  int
	UID     int
	TradeID int
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
	h.buy = queue.NewQueue(queue.Name("Buy"))
	h.sell = queue.NewQueue(queue.Name("Sell"))
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

// Add ...
func (h *Bid) Add(q queue.Server, u *Unit) (queue.Data, error) {
	buffer := BufferMessage{Queue: q}
	data := queue.NewData(u)
	//如果是买家队列，按照价格高优先，时间优先
	if h.buy == q {
		for n := q.Front(); ; n = n.Next() {
			if n == nil {
				buffer.Node = q.Push(data)
				break
			}
			content := n.Data.Content.(*Unit)
			if content.UID == u.UID && content.Price == u.Price {
				content.Amount += u.Amount
				n.Data.Content = content
				buffer.Node = n
				break
			}

			//价格高者优先
			if u.Price > content.Price {
				buffer.Node = q.InsertBefore(data, n)
				break
			}

			//时间优先
			if u.Price == content.Price && n.Data.CreateAt.After(data.CreateAt) {
				buffer.Node = q.InsertBefore(data, n)
				break
			}
		}
	}

	if h.sell == q {
		for n := q.Front(); ; n = n.Next() {
			if n == nil {
				buffer.Node = q.Push(data)
				break
			}
			content := n.Data.Content.(*Unit)
			if content.UID == u.UID && content.Price == u.Price {
				content.Amount += u.Amount
				n.Data.Content = content
				buffer.Node = n
				break
			}

			//价格高者优先
			if content.Price > u.Price {
				buffer.Node = q.InsertBefore(data, n)
				break
			}

			//时间优先
			if content.Price == u.Price && n.Data.CreateAt.After(data.CreateAt) {
				buffer.Node = q.InsertBefore(data, n)
				break
			}
		}
	}

	h.opts.buffer <- &buffer
	return *data, nil
}

// Cancel ...
func (h *Bid) Cancel(q queue.Server, tradeID int) error {
	q.Loop(func(n *queue.Node) error {
		content := n.Data.Content
		if content != nil && content.(*Unit).TradeID == tradeID {
			q.Remove(n)
		}
		return nil
	})
	return nil
}

// Buffer 获取缓冲
func (h *Bid) Buffer() <-chan interface{} {
	return h.opts.buffer
}
