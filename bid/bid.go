package bid

import (
	"fmt"
	sync "sync"

	"github.com/zzsds/trade/queue"
)

// Server ...
type Server interface {
	ID() int
	Name() string
	Init() Server
	Buy() *Data
	Sell() *Data
	Add(*Unit) error
	Remove(*Data, int, float64) error
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
	buy, sell *Data
}

// Data ...
type Data struct {
	*sync.RWMutex
	queue.Server
}

// NewData ...
func NewData(name string) *Data {
	return &Data{
		&sync.RWMutex{},
		queue.NewQueue(queue.Name(name)),
	}
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
	h.buy = NewData("BUY")
	h.sell = NewData("SELL")
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

// Amount ...
func (h *Bid) Amount() int {
	return h.opts.amount
}

// Buy ...
func (h *Bid) Buy() *Data {
	return h.buy
}

// Sell ...
func (h *Bid) Sell() *Data {
	return h.sell
}

func (h *Bid) add(q queue.Server, u *Unit) error {
	var newNode *queue.Node
	if q.Len() <= 0 {
		newNode = q.PushFront(u)
		return nil
	}

	if q.Len() == 0 {
		newNode = q.PushFront(u)
	} else {
		for n := q.Front(); n != nil; n = n.Next() {
			v, ok := n.Value.(*Unit)
			if !ok {
				return fmt.Errorf("asset fail")
			}

			if v.Price == u.Price {
				if v.UID == u.UID {
					v.Amount += u.Amount
					n.Value = v
					newNode = n
					break
				}
				if v.CreateAt.After(u.CreateAt) {
					newNode = q.InsertBefore(u, n)
					break
				} else {
					newNode = q.InsertAfter(u, n)
					break
				}
			}

			//如果是买家队列，按照价格高优先，时间优先
			if h.buy == q && u.Price > v.Price {
				//价格高者优先
				newNode = q.InsertBefore(u, n)
				break
			} else if h.sell == q && v.Price > u.Price {
				//价格高者优先
				newNode = q.InsertBefore(u, n)
				break
			}
		}
	}
	h.opts.amount++

	fmt.Println(newNode)
	defer func(n *queue.Node) {
		if n == nil {
			return
		}
		h.opts.buffer <- Message{
			Queue: q,
			Node:  n,
		}
	}(newNode)
	return nil
}

// Add ...
func (h *Bid) Add(u *Unit) error {
	q := h.buy
	if u.Type == Type_Sell {
		q = h.sell
	}
	q.Lock()
	defer q.Unlock()

	return h.add(q, u)
}

// Remove ...
func (h *Bid) Remove(d *Data, uid int, price float64) error {
	d.Lock()
	defer d.Unlock()
	return d.Loop(func(n *queue.Node) error {
		u := n.Value.(*Unit)
		if n != nil && u.UID == uid && u.Price == price {
			d.Remove(n)
		}
		return nil
	})
}

// Buffer 获取缓冲
func (h *Bid) Buffer() <-chan Message {
	return h.opts.buffer
}
