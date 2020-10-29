package match

import "time"

// Buy ...
type Buy struct {
	ID       uint
	Number   uint
	Amount   uint
	CreateAt *time.Time
}

// NewBuy buy queue
func NewBuy() *Queue {
	return &Queue{}
}
