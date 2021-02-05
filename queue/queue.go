package queue

import (
	"fmt"
	"os"
)

// Server 对外服务接口
type Server interface {
	Len() int
	Name() string
	Front() *Node
	Back() *Node
	MoveAfter(e, mark *Node)
	MoveBefore(e, mark *Node)
	MoveToBack(e *Node)
	MoveToFront(e *Node)
	PushBack(v interface{}) *Node
	PushBackQueue(other *Queue)
	PushFront(v interface{}) *Node
	PushFrontQueue(other *Queue)
	InsertBefore(v interface{}, mark *Node) *Node
	InsertAfter(v interface{}, mark *Node) *Node
	Remove(*Node) interface{}
	List() []interface{}
	Loop(Call) error
	Get(int) interface{}
}

// Call ...
type Call func(*Node) error

// NodeServer ...
type NodeServer interface {
	Current() *Node
	Next() *Node
	Prev() *Node
	Value() interface{}
}

// Node is an element of a linked list.
type Node struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Node

	// The list to which this element belongs.
	list *Queue

	// The value stored with this element.
	Value interface{}
}

// NewNode ...
func NewNode(v interface{}) *Node {
	node := new(Node).init()
	node.Value = v
	return node
}

// Next returns the next list Node or nil.
func (e *Node) init() *Node {
	e.prev = nil
	e.next = nil
	e.list = nil
	return e
}

// Next returns the next list element or nil.
func (e *Node) Next() *Node {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Node) Prev() *Node {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Queue represents a doubly linked list.
// The zero value for Queue is an empty list ready to use.
type Queue struct {
	opts options
	root Node // sentinel list element, only &root, root.prev, and root.next are used
	len  int  // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *Queue) Init() *Queue {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewQueue returns an initialized list.
func NewQueue(opts ...Option) Server {
	queue := new(Queue).Init()
	queue.opts = newOptions(opts...)
	return queue
}

// Name 名称
func (l *Queue) Name() string {
	return l.opts.name
}

// Len returns the Amount of elements of list l.
// The complexity is O(1).
func (l *Queue) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *Queue) Front() *Node {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *Queue) Back() *Node {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero Queue value.
func (l *Queue) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *Queue) insert(e, at *Node) *Node {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%#v \n\t 换行", e)
			fmt.Println("\n\t 隔断")
			fmt.Printf("%#v \n\t 结束", at)
			os.Exit(1)
		}
	}()
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Node{Value: v}, at).
func (l *Queue) insertValue(v interface{}, at *Node) *Node {
	return l.insert(&Node{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *Queue) remove(e *Node) *Node {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// move moves e to next to at and returns e.
func (l *Queue) move(e, at *Node) *Node {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e

	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *Queue) Remove(e *Node) interface{} {
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Node) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *Queue) PushFront(v interface{}) *Node {
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *Queue) PushBack(v interface{}) *Node {
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *Queue) InsertBefore(v interface{}, mark *Node) *Node {
	if mark.list != l {
		return nil
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	// see comment in Queue.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *Queue) InsertAfter(v interface{}, mark *Node) *Node {
	if mark.list != l {
		return nil
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	// see comment in Queue.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *Queue) MoveToFront(e *Node) {
	if e.list != l || l.root.next == e {
		return
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	// see comment in Queue.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *Queue) MoveToBack(e *Node) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	// see comment in Queue.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Queue) MoveBefore(e, mark *Node) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *Queue) MoveAfter(e, mark *Node) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.move(e, mark)
}

// PushBackQueue inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Queue) PushBackQueue(other *Queue) {
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontQueue inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *Queue) PushFrontQueue(other *Queue) {
	l.opts.mutex.Lock()
	defer l.opts.mutex.Unlock()
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// Loop 单次循环
func (l *Queue) Loop(f Call) error {
	for node := l.Front(); node != nil; node = node.Next() {
		if err := f(node); err != nil {
			break
		}
	}
	return nil
}

// Get 根据索引查询
func (l *Queue) Get(index int) interface{} {
	if l.len == 0 || l.len < int(index) {
		return nil
	}
	i := 0
	for node := l.Front(); node != nil; node = node.Next() {
		if i == index {
			return node.Value
		}
		i++
	}
	return nil
}

// List ...
func (l *Queue) List() []interface{} {
	list := make([]interface{}, 0, l.len)
	for node := l.Front(); node != nil; node = node.Next() {
		list = append(list, node.Value)
	}
	return list
}
