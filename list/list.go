package list

import "container/list"

// Node ...
type Node = list.Element

// CallOption ...
type CallOption func(*Node) bool

// Server ...
type Server interface {
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
	CallList(CallOption)
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
func (l List) Name() string {
	return l.opts.Name
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

// CallList ...
func (l *List) CallList(fn CallOption) {
	l.opts.Lock()
	defer l.opts.Unlock()
	for n := l.Front(); n != nil; n = n.Next() {
		if !fn(n) {
			break
		}
	}
}
