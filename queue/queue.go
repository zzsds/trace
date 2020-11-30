package queue

import (
	"time"

	"github.com/google/uuid"
)

// Server 对外服务接口
type Server interface {
	Len() int
	Name() string
	Get(index uint) *Node
	Push(node *Node) bool
	Unshift(node *Node) bool
	Shift() *Node
	Pop() *Node
	BeforeAdd(*Node, *Node) error
	AfterAdd(*Node, *Node) error
	Header() *Node
	Tailed() *Node
	Reverse()
	Buffer() <-chan Node
}

// Node ...
type Node struct {
	prev, next *Node
	queue      *Queue
	Data       *Data
}

// NewNode ...
func NewNode(data *Data) *Node {
	node := new(Node).init()
	node.Data = data
	return node
}

// Next returns the next list Node or nil.
func (e *Node) init() *Node {
	e.Data = nil
	e.prev = nil
	e.next = nil
	e.queue = nil
	return e
}

// Next returns the next list Node or nil.
func (e *Node) Next() *Node {
	if p := e.next; e.queue != nil && p != e.queue.head {
		return p
	}
	return nil
}

// Prev returns the previous list Node or nil.
func (e *Node) Prev() *Node {
	if p := e.prev; e.queue != nil && p != e.queue.head {
		return p
	}
	return nil
}

// Queue ...
type Queue struct {
	opts       Options
	len        int
	head, tail *Node
}

// NewQueue new queue
func NewQueue(opts ...Option) Server {
	queue := new(Queue).Init()
	queue.opts = newOptions(opts...)
	// 启动过期监听
	go queue.expireListen()
	return queue
}

// Init initializes or clears list l.
func (h *Queue) Init() *Queue {
	h.head = nil
	h.tail = nil
	h.len = 0
	return h
}

// Len ...
func (h *Queue) Len() int { return h.len }

// Name ...
func (h *Queue) Name() string {
	return h.opts.name
}

// Front returns the first element of list l or nil if the list is empty.
func (h *Queue) Front() *Node {
	if h.len == 0 {
		return nil
	}
	return h.head.next
}

// Header ...
func (h *Queue) Header() *Node {
	return h.head
}

// Tailed ...
func (h *Queue) Tailed() *Node {
	return h.tail
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
	if h.len == 0 || h.len < int(index) {
		return nil
	}
	if index == 0 {
		return h.head
	}
	node := h.head

	return node
}

// Unshift 开头插入
func (h *Queue) Unshift(node *Node) bool {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	node.next = h.head
	h.head = node
	h.len++
	return true
}

// Shift 移出开始第一个
func (h *Queue) Shift() *Node {
	if h.head == nil {
		return nil
	}
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	head := h.head
	h.head = head.next
	h.len--
	return head
}

// Push 压入末尾
func (h *Queue) Push(node *Node) bool {
	if node == nil {
		return false
	}
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()

	if h.len == 0 {
		h.head = node
		h.tail = node
		node.next = nil
		node.prev = nil
	} else {
		node.prev = h.tail
		node.next = nil
		h.tail.next = node
		h.tail = node
	}
	h.len++
	return true
}

// Pop 弹出结尾最后一个
func (h *Queue) Pop() *Node {
	if h.len <= 0 {
		return nil
	}
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	tail := h.tail
	h.tail = tail.prev
	h.tail.next = nil

	return tail
}

// BeforeAdd 之前插入
func (h *Queue) BeforeAdd(before *Node, value *Node) error {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	h.len++
	value.prev = before.prev
	value.next = before
	return nil
}

// AfterAdd 之后插入
func (h *Queue) AfterAdd(after *Node, value *Node) error {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	h.len++
	value.prev = after
	value.next = after.next
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

// Remove 替换
func (h *Queue) Remove(node *Node) error {
	h.opts.mutex.Lock()
	defer h.opts.mutex.Unlock()
	if h.len <= 0 {
		return nil
	}

	if node.prev == nil {
		h.head = node.next
		node.next.prev = nil
		return nil
	}

	if node.next == nil {
		node.prev.next = nil
		h.tail = node.prev
		return nil
	}

	h.len--
	node.prev.next = node.next
	node.next.prev = node.prev
	return nil
}

// expireListen ...
func (h *Queue) expireListen() {
	go func() {
		for {
			node := h.head
			for node != nil {
				if node.Data.ExpireAt != nil && time.Now().After(*node.Data.ExpireAt) {
					h.opts.buffer <- *node
					h.Remove(node)
				}
				node = node.next
			}
		}
	}()
}

// Data ...
type Data struct {
	UUID     string
	CreateAt time.Time
	ExpireAt *time.Time
	Content  interface{}
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
