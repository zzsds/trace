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
	Buffer() <-chan Node
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
	}
	// Data ...
	Data struct {
		UUID     string
		CreateAt time.Time
		ExpireAt *time.Time
		Content  interface{}
	}
)

// NewQueue new queue
func NewQueue(opts ...Option) Server {
	queue := &Queue{opts: newOptions(opts...), Size: 0, Head: nil, Tail: nil}
	// 启动过期监听
	go queue.expireListen()
	return queue
}

// Length ...
func (h *Queue) Length() int {
	return h.Size
}

// Name ...
func (h *Queue) Name() string {
	return h.opts.name
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
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
}

// Buffer 获取缓存
func (h *Queue) Buffer() <-chan Node {
	return h.opts.buffer
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

	return node
}

// Unshift 开头插入
func (h *Queue) Unshift(node *Node) bool {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
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
func (h *Queue) BeforeAdd(before *Node, value *Node) error {
	return nil
}

// AfterAdd 之后插入
func (h *Queue) AfterAdd(after *Node, value *Node) error {
	return nil
}

// Unique 去重
func (h *Queue) Unique() error {
	return nil
}

// Replace 替换
func (h *Queue) Replace(old *Node, new *Node) error {
	return nil
}

// Delete 替换
func (h *Queue) Delete(node *Node) error {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	if h.Size <= 0 {
		return nil
	}

	if node.Prev == nil {
		h.Head = node.Next
		node.Next.Prev = nil
		return nil
	}

	if node.Next == nil {
		node.Prev.Next = nil
		h.Tail = node.Prev
		return nil
	}

	h.Size--
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	return nil
}

// expireListen ...
func (h *Queue) expireListen() {
	go func() {
		for {
			node := h.Head
			for node != nil {
				if node.Data.ExpireAt != nil && time.Now().After(*node.Data.ExpireAt) {
					h.opts.buffer <- *node
					h.Delete(node)
				}
				node = node.Next
			}
		}
	}()
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

// NewExpireData ...
func NewExpireData(content interface{}, expire *time.Time) *Data {
	uuid, _ := uuid.NewUUID()
	return &Data{
		UUID:     uuid.String(),
		CreateAt: time.Now(),
		Content:  content,
		ExpireAt: expire,
	}
}

// AddData 插入数据到队列
func (h *Queue) AddData(content interface{}) {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
}
