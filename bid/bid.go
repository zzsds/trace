package bid

import (
	"fmt"
)

// Server 交易对接口
type Server interface {
	Init() error
	ID() int
	Name() string
	Buy() DataServer
	Sell() DataServer
	Add(*Unit) error
	Buffer() <-chan Message
}

var bids = make(map[int]Server)

// Bid bid
type Bid struct {
	opts      options
	buy, sell DataServer
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
func (h *Bid) Init() error {
	h.buy = NewData(WithName("BUY"), WithSort(Sort_Desc))
	h.sell = NewData(WithName("SELL"), WithSort(Sort_Asc))
	return nil
}

// ID ...
func (h Bid) ID() int {
	return h.opts.id
}

// Name ...
func (h Bid) Name() string {
	return h.opts.name
}

// Amount ...
func (h Bid) Amount() int {
	return h.opts.amount
}

// Buy ...
func (h *Bid) Buy() DataServer {
	return h.buy
}

// Sell ...
func (h *Bid) Sell() DataServer {
	return h.sell
}

// Buffer ...
func (h *Bid) Buffer() <-chan Message {
	return h.opts.buffer
}

// Add ...
func (h *Bid) Add(unit *Unit) error {
	var t DataServer
	if unit.Type == Type_Buy {
		t = h.buy
	} else {
		t = h.sell
	}
	return t.Add(*unit)
}
