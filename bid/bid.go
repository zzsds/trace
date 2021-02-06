package bid

import (
	"fmt"
	"sort"

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
	buy, sell *Data
}

// Data ...
type Data struct {
	Sort  Sort
	Queue queue.Server
	Map   map[float64]*queue.Node
	Array []float64
}

// NewData ...
func NewData(name string, sort Sort) *Data {
	return &Data{
		Sort:  Sort_Asc,
		Queue: queue.NewQueue(queue.Name(name)),
		Map:   make(map[float64]*queue.Node),
		Array: make([]float64, 0, 100),
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
	h.buy = NewData(Type_Buy.String(), Sort_Desc)
	h.sell = NewData(Type_Sell.String(), Sort_Asc)
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
	return h.buy.Queue
}

// Sell ...
func (h *Bid) Sell() queue.Server {
	return h.sell.Queue
}

// Amount ...
func (h *Bid) Amount() int {
	return h.opts.amount
}

// Add ...
func (h *Bid) Add(u *Unit) (Unit, error) {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	t := h.buy
	if u.Type == Type_Sell {
		t = h.sell
	}

	q := t.Queue
	var newNode *queue.Node
	if q.Len() <= 0 {
		newNode = q.PushFront(u)
		t.Map[u.Price] = newNode
		t.Array = append(t.Array, u.Price)
		return *u, nil
	}

	if n, ok := t.Map[u.Price]; ok {
		v, ok := n.Value.(*Unit)
		if !ok {
			return Unit{}, nil
		}
		if v.UID == u.UID {
			v.Amount += u.Amount
			n.Value = v
			newNode = n
			return *u, nil
		}

		if v.CreateAt.After(u.CreateAt) {
			newNode = q.InsertBefore(u, n)
		} else {
			newNode = q.InsertAfter(u, n)
		}
		return *u, nil
	}

	for _, p := range t.Array {
		//如果是买家队列，按照价格高优先，时间优先
		if h.buy.Queue == q && u.Price > p {
			//价格高者优先
			newNode = q.InsertBefore(u, t.Map[p])
			break
		} else if h.sell.Queue == q && p > u.Price {
			//价格高者优先
			newNode = q.InsertBefore(u, t.Map[p])
			break
		}
	}

	t.Array = append(t.Array, u.Price)
	t.Map[u.Price] = newNode
	h.opts.amount++

	sort.Float64Slice(t.Array).Sort()
	return *u, nil

	message := Message{Queue: q}
	if q.Len() == 0 {
		message.Node = q.PushFront(u)
	} else {
		for n := q.Front(); n != nil; n = n.Next() {
			v, ok := n.Value.(*Unit)
			if !ok {
				break
			}
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
				} else {
					message.Node = q.InsertAfter(u, n)
					break
				}
			}

			//如果是买家队列，按照价格高优先，时间优先
			if h.buy.Queue == q && u.Price > v.Price {
				//价格高者优先
				message.Node = q.InsertBefore(u, n)
				break
			} else if h.sell.Queue == q && v.Price > u.Price {
				//价格高者优先
				message.Node = q.InsertBefore(u, n)
				break
			}
		}
	}
	h.opts.amount++
	// h.opts.buffer <- message

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
