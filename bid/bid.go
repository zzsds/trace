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
	Add(queue.Server, *Unit) error
	Buffer() <-chan interface{}
}

// BufferMessage 缓冲消息
type BufferMessage struct {
	Queue queue.Server
	*Unit
}

// Unit ...
type Unit struct {
	Name    string
	Price   float64
	Number  int
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
	return bid
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
	return h.buy
}

// Add ...
func (h *Bid) Add(q queue.Server, u *Unit) error {
	data := queue.NewData(u)
	fmt.Println(u.Price, u.TradeID)
	if q.Len() <= 0 {
		q.Unshift(data)
		goto END
	}
	if h.buy == q {
		for n := q.Front(); n != nil; n = n.Next() {
			content := n.Data.Content.(*Unit)
			fmt.Println(content.Price, u.Price, u.TradeID)
			if content.UID == u.UID && content.Price == u.Price {
				content.Number += u.Number
				n.Data.Content = content
				continue
			}
			if content.Price < u.Price {
				q.InsertBefore(queue.NewNode(data), n)
				// q.Unshift(queue.NewData(u))
				break
			}
		}

		// fmt.Printf("%p ==== %p  %v \n\t", h.buy, q, u)
	}

END:
	h.opts.buffer <- &BufferMessage{
		q, u,
	}
	return nil
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
