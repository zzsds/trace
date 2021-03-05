package bid

import (
	"errors"
	"log"

	"github.com/zzsds/trade/list"
)

// Server ...
type Server interface {
	Init(...Option)
	Name() string
	String() string
	Add(*Unit) error
	Cancel(*Unit) error
	Buy() list.Server
	Sell() list.Server
	Queue() <-chan interface{}
}

// Bid ...
type bid struct {
	ops       options
	queue     chan interface{}
	buy, sell list.Server
}

// NewBid ...
func NewBid(opts ...Option) Server {
	b := new(bid)
	b.ops = newOptions(opts...)
	b.queue = make(chan interface{})
	b.Init()
	return b
}

// Init ...
func (b *bid) Init(opts ...Option) {
	// Buy
	if b.buy != nil {
		b.buy.Lock()
		defer b.buy.Unlock()
	}
	b.buy = list.NewList(func(o *list.Options) {
		o.Name = Type_Buy.String()
	})
	// Sell
	if b.sell != nil {
		b.sell.Lock()
		defer b.sell.Unlock()
	}
	b.sell = list.NewList(func(o *list.Options) {
		o.Name = Type_Sell.String()
	})
}

// Name ...
func (b *bid) Name() string {
	return b.ops.name
}

// Name ...
func (b *bid) String() string {
	return "bid"
}

// Queue ...
func (b *bid) Queue() <-chan interface{} {
	return b.queue
}

// Buy ...
func (b *bid) Buy() list.Server {
	return b.buy
}

// Sell ...
func (b *bid) Sell() list.Server {
	return b.sell
}

// Add ...
func (b *bid) Add(unit *Unit) error {
	var object list.Server
	if unit.Type == Type_Buy {
		object = b.buy
	} else {
		object = b.sell
	}
	if n := b.add(object, unit); n != nil {
		b.queue <- n
	}
	return nil
}

// AddBuy ...
func (b *bid) add(l list.Server, u *Unit) *list.Node {
	l.Lock()
	defer l.Unlock()
	if l.Len() <= 0 {
		return l.PushFront(u)
	}

	var node *list.Node
	for n := l.Front(); n != nil; n = n.Next() {
		v, ok := n.Value.(*Unit)
		if !ok {
			log.Println("Parsing failed")
			return nil
		}
		if v.Price == u.Price {
			if v.UID == u.UID {
				v.Amount += u.Amount
				n.Value = v
				break
			}
			if v.CreateAt < u.CreateAt {
				node = l.InsertBefore(u, n)
			} else {
				node = l.InsertAfter(u, n)
			}
			return node
		}

		var t bool
		if l == b.buy {
			// 降序，按照价格高优先，时间优先  买
			t = u.Price > v.Price
		} else {
			// 升序，按照价格高优先，时间优先  卖
			t = u.Price < v.Price
		}

		if t {
			return l.InsertBefore(u, n)
		}
	}
	return nil
}

func (b *bid) Cancel(unit *Unit) error {
	var object list.Server
	if unit.Type == Type_Buy {
		object = b.buy
	} else {
		object = b.sell
	}

	return b.cancel(object, unit)
}

func (b *bid) cancel(l list.Server, u *Unit) error {
	l.Lock()
	defer l.Unlock()
	for n := l.Front(); n != nil; n = n.Next() {
		v, ok := n.Value.(*Unit)
		if !ok {
			return errors.New("Parsing failed")
		}
		if v.Price == u.Price && v.UID == u.UID {
			if v.Amount > u.Amount {
				v.Amount -= u.Amount
				n.Value = v
				break
			}

			if v.Amount <= u.Amount {
				l.Remove(n)
				break
			}
		}
	}

	return nil
}
