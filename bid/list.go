package bid

import (
	"container/list"
)

// Node ...
type Node = list.Element

// ListServer ...
type ListServer interface {
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
	*list.List
}

// Edit ...
func (l *List) Edit(node *Node) interface{} {
	return nil
}
