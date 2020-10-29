package queue

type Server interface {
	Size() int
}

type (
	Queue struct {
		opts Options
		size int
		head *Node
		tail *Node
	}
	Node struct {
		prve  *Node
		next  *Node
		value interface{}
	}
)

func NewQueue(opts ...Option) Server {
	return &Queue{opts: newOptions(opts), size: 0}
}

func (h *Queue) Size() int {
	return h.size
}
