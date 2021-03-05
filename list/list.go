package list

import (
	"container/list"
	"sync"
)

// Node ...
type Node = list.Element

// NodeFunc ...
type NodeFunc func(*Node) (*Node, error)

// Server ...
type Server interface {
	sync.Locker
	Name() string
	Len() int
	Front() *Node
	Back() *Node
	Remove(*Node) interface{}
	PushFront(interface{}) *Node
	PushBack(interface{}) *Node
	InsertBefore(interface{}, *Node) *Node
	InsertAfter(interface{}, *Node) *Node
	NodeList() []*Node
	CallListHandler(NodeFunc) error
}

// List ...
type List struct {
	opts Options
	*list.List
}

// NewList ...
func NewList(opts ...Option) Server {
	return &List{newOptions(opts...), list.New()}
}

// Name ...
func (l *List) Name() string {
	return l.opts.Name
}

// Lock ...
func (l *List) Lock() {
	l.opts.Lock()
}

// Unlock ...
func (l *List) Unlock() {
	l.opts.Unlock()
}

// NodeList ...
func (l *List) NodeList() []*Node {
	l.opts.Lock()
	defer l.opts.Unlock()
	var marks []*Node
	for n := l.Front(); n != nil; n = n.Next() {
		marks = append(marks, n)
	}
	return marks
}

// CallListHandler ...
func (l *List) CallListHandler(call NodeFunc) error {
	l.opts.Lock()
	defer l.opts.Unlock()
	for n := l.Front(); n != nil; n = n.Next() {
		c, err := call(n)
		if err != nil {
			return err
		}
		if c == nil {
			break
		}
	}
	return nil
}
