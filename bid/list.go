package bid

import (
	"container/list"
	sync "sync"
)

// Node ...
type Node = list.Element

// ListServer ...
type ListServer interface {
	sync.Locker
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

// CallOption ...
type CallOption func(*Node) bool

// List ...
type List struct {
	*sync.RWMutex
	*list.List
}

// Len ...
func (l *List) Len() int {
	l.Lock()
	defer l.Unlock()
	return l.List.Len()
}

// PushFront ...
func (l *List) PushFront(v interface{}) *Node {
	l.Lock()
	defer l.Unlock()
	return l.List.PushFront(v)
}

// NodeList ...
func (l *List) NodeList() []*Node {
	l.Lock()
	defer l.Unlock()
	var marks []*Node
	for n := l.Front(); n != nil; n = n.Next() {
		marks = append(marks, n)
	}
	return marks
}

// CallList ...
func (l *List) CallList(fn CallOption) {
	l.Lock()
	defer l.Unlock()
	for n := l.Front(); n != nil; n = n.Next() {
		if !fn(n) {
			break
		}
	}
}
