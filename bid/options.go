package bid

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// options ...
type options struct {
	mutex     *sync.RWMutex
	id        int
	name      string
	amount    int
	maxAmount int
	buffer    chan interface{}
	signal    bool
}

// Message ...
type Message struct{}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		mutex:  &sync.RWMutex{},
		buffer: make(chan interface{}),
		signal: true,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name ...
func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// OptionUnit ...
type OptionUnit func(*Unit)

// Unit ...
type Unit struct {
	ID       int       `json:"id"`
	UID      int       `json:"uid"`
	UUID     string    `json:"uuid"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Amount   int       `json:"amount"`
	Type     Type      `json:"type"`
	CreateAt time.Time `json:"createAt"`
}

// NewUnit ...
func NewUnit(unit ...OptionUnit) *Unit {
	uuid, _ := uuid.NewUUID()
	u := &Unit{
		CreateAt: time.Now(),
		UUID:     uuid.String(),
	}
	for _, o := range unit {
		o(u)
	}
	return u
}
