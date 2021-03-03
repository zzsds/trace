package bid

import (
	"fmt"

	"github.com/zzsds/trade/store"
)

// Server ...
type Server interface {
	Init(...Option)
	Name() string
	String() string
	Add(*Unit) error
}

// Bid ...
type Bid struct {
	ops       Options
	Queue     store.Store
	Buy, Sell store.Store
}

// NewBid ...
func NewBid(opts ...Option) Server {
	b := new(Bid)
	b.ops = newOptions(opts...)
	b.Init()
	return b
}

// Init ...
func (b *Bid) Init(opts ...Option) {

}

// Name ...
func (b *Bid) Name() string {
	return b.ops.Name
}

// Name ...
func (b *Bid) String() string {
	return "bid"
}

// Add ...
func (b *Bid) Add(unit *Unit) error {
	fmt.Println(unit)
	return nil
}
