package queue

import (
	"time"

	"github.com/google/uuid"
)

// Server ...
type Server interface {
	Length() int
	Name() string
	Get(index uint) *Node
	Push(node *Node) bool
	Unshift(node *Node) bool
	Shift() *Node
	Pop() *Node
	Header() *Node
	Tailed() *Node
	Reverse()
}

type (
	// Queue ...
	Queue struct {
		opts Options
		Size int
		Head *Node
		Tail *Node
	}
	// Node ...
	Node struct {
		Prev *Node
		Next *Node
		Data *Data
		Sort uint
	}
	// Data ...
	Data struct {
		UUID     string
		CreateAt time.Time
		Content  interface{}
	}
)

// NewQueue new queue
func NewQueue(opts ...Option) Server {
	return &Queue{opts: newOptions(opts...), Size: 0, Head: nil, Tail: nil}
}

// Length ...
func (h *Queue) Length() int {
	return h.Size
}

// Name ...
func (h *Queue) Name() string {
	return h.opts.Name
}

// Header ...
func (h *Queue) Header() *Node {
	return h.Head
}

// Tailed ...
func (h *Queue) Tailed() *Node {
	return h.Tail
}

// Reverse ...
func (h *Queue) Reverse() {

}

// Get ...
func (h *Queue) Get(index uint) *Node {
	if h.Size == 0 || h.Size < int(index) {
		return nil
	}
	if index == 0 {
		return h.Head
	}
	node := h.Head

	for {
		if node.Sort == index {
			break
		}
		node = node.Next
	}
	return node
}

// Unshift 开头插入
func (h *Queue) Unshift(node *Node) bool {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	node.Sort = 0
	head := h.Head
	for head != nil {
		head.Sort++
		head = head.Next
	}
	node.Next = h.Head
	h.Head = node
	h.Size++
	return true
}

// Shift 移出开始第一个
func (h *Queue) Shift() *Node {
	if h.Head == nil {
		return nil
	}
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	head := h.Head
	h.Head = h.Head.Next
	for head := h.Head; head != nil; head = head.Next {
		head.Sort--
	}
	h.Size--
	return head
}

// Push 压入末尾
func (h *Queue) Push(node *Node) bool {
	if node == nil {
		return false
	}
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()

	node.Sort = uint(h.Size)
	if h.Size == 0 {
		h.Head = node
		h.Tail = node
		node.Next = nil
		node.Prev = nil
	} else {
		node.Prev = h.Tail
		node.Next = nil
		h.Tail.Next = node
		h.Tail = node
	}
	h.Size++
	return true
}

// Pop 弹出结尾最后一个
func (h *Queue) Pop() *Node {
	tail := h.Tail
	for head := h.Head; head != nil; head = head.Next {
		if head.Next == tail {
			head.Next = nil
		}
	}
	h.Tail = h.Tail.Prev
	return tail
}

// BeforeAdd 之前插入
func (h *Queue) BeforeAdd(before, value interface{}) error {
	return nil
}

// AfterAdd 之后插入
func (h *Queue) AfterAdd(after, value interface{}) error {
	return nil
}

// Unique 去重
func (h *Queue) Unique() error {
	return nil
}

// Replace 替换
func (h *Queue) Replace(old, value interface{}) error {
	return nil
}

// NewData ...
func NewData(content interface{}) *Data {
	uuid, _ := uuid.NewUUID()
	return &Data{
		UUID:     uuid.String(),
		CreateAt: time.Now(),
		Content:  content,
	}
}
