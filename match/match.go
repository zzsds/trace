package match

// DefaultServer ...
var DefaultServer Server

// Server ...
type Server interface {
	Run() error
}

// Unit ...
type Unit struct {
	UID    uint
	Number uint
	Amount float64
}

// QueueNode 队列
type QueueNode struct {
	Number uint
	Value  *Unit
	Prev   *Queue
	Next   *Queue
}

// NewMatch ...
func NewMatch() Server {
	var (
		buy = NewBuy()
	)
}

// Run ...
func (h *Queue) Run() error {

	return h.Stop()
}

// Start ...
func (h *Queue) Start() error {
	return nil
}

// Stop ...
func (h *Queue) Stop() error {
	return nil
}
