package bid

import (
	"time"

	"github.com/google/uuid"
)

// Options contains configuration for the Store
type options struct {
	name string
}

// Name ...
func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// Option sets values in Options
type Option func(o *options)

func newOptions(opts ...Option) options {
	opt := options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Unit ...
type Unit struct {
	ID       int     `json:"id"`
	UID      int     `json:"uid"`
	UUID     string  `json:"uuid"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Amount   int     `json:"amount"`
	Type     Type    `json:"type"`
	CreateAt int64   `json:"createAt"`
}

// UnitOption ...
type UnitOption func(*Unit)

// NewUnit ...
func NewUnit(unit ...UnitOption) *Unit {
	uuid, _ := uuid.NewUUID()
	u := &Unit{
		CreateAt: time.Now().Unix(),
		UUID:     uuid.String(),
	}
	for _, o := range unit {
		o(u)
	}
	return u
}
