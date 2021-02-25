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
	Edit(*Node) interface{}
	Remove(*Node) interface{}
	PushFront(interface{}) *Node
	PushBack(interface{}) *Node
	InsertBefore(interface{}, *Node) *Node
	InsertAfter(interface{}, *Node) *Node
}

// List ...
type List struct {
	*sync.RWMutex
	*list.List
}

// Edit ...
func (l *List) Edit(node *Node) interface{} {
	return nil
}

// Remove ...
func (l *List) Remove(e *Node) interface{} {
	l.Lock()
	defer l.Unlock()
	return l.List.Remove(e)
}

// PushFront ...
func (l *List) PushFront(v interface{}) *Node {
	l.Lock()
	defer l.Unlock()
	return l.List.PushFront(v)
}

// PushBack ...
func (l *List) PushBack(v interface{}) *Node {
	l.Lock()
	defer l.Unlock()
	return l.List.PushBack(v)
}

// InsertBefore ...
func (l *List) InsertBefore(v interface{}, mark *Node) *Node {
	l.Lock()
	defer l.Unlock()
	return l.List.InsertBefore(v, mark)
}

// InsertAfter ...
func (l *List) InsertAfter(v interface{}, mark *Node) *Node {
	l.Lock()
	defer l.Unlock()
	return l.List.InsertAfter(v, mark)
}
