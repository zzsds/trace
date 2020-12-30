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
	buffer    chan Message
	signal    bool
}

// Option ...
type Option func(*options)

func newOptions(opts ...Option) options {
	opt := options{
		mutex:  &sync.RWMutex{},
		buffer: make(chan Message),
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

// DefaultUnit ...
type DefaultUnit func(*Unit)

// Unit ...
type Unit struct {
	ID       int
	UID      string
	Name     string
	Price    float64
	Amount   int
	Type     Type
	CreateAt time.Time
}

// NewUnit ...
func NewUnit(unit ...DefaultUnit) *Unit {
	uuid, _ := uuid.NewUUID()
	u := &Unit{
		CreateAt: time.Now(),
		UID:      uuid.String(),
	}
	for _, o := range unit {
		o(u)
	}
	return u
}
