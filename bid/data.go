package bid

import (
	"container/list"
	"fmt"
	"log"
	"sync"
)

// DataServer ...
type DataServer interface {
	ListServer
	Init() error
	Name() string
	Sort() Sort
	Add(Unit) (*Node, error)
}

// Data ...
type Data struct {
	*List
	sort   Sort
	name   string
	buffer chan interface{}
}

// DataOption ...
type DataOption func(*Data)

// WithBuffer ...
func WithBuffer(ch chan interface{}) DataOption {
	return func(d *Data) {
		d.buffer = ch
	}
}

// WithSort 设置排序方式
func WithSort(sort Sort) DataOption {
	return func(d *Data) {
		d.sort = sort
	}
}

// WithName 设置数据名称
func WithName(name string) DataOption {
	return func(d *Data) {
		d.name = name
	}
}

// NewData ...
func NewData(opts ...DataOption) DataServer {
	opt := Data{}
	for _, o := range opts {
		o(&opt)
	}
	if opt.Init() != nil {
		panic("Init fail")
	}
	return &opt
}

// Init ...
func (h *Data) Init() error {
	h.List = &List{&sync.RWMutex{}, list.New()}
	return nil
}

// Name ...
func (h Data) Name() string {
	return h.name
}

// Sort ...
func (h Data) Sort() Sort {
	return h.sort
}

// Add ...
func (h *Data) Add(u Unit) (*Node, error) {
	if h.Len() <= 0 {
		return h.PushFront(u), nil
	}

	var node *Node
	h.CallList(func(n *Node) bool {
		v, ok := n.Value.(Unit)
		if !ok {
			log.Fatalln(fmt.Errorf("Parsing failed"))
			return false
		}
		if v.Price == u.Price {
			if v.UID == u.UID {
				v.Amount += u.Amount
				n.Value = v
				// node = n
				// break
				return true
			}
			if v.CreateAt.After(u.CreateAt) {
				node = h.InsertBefore(u, n)
			} else {
				node = h.InsertAfter(u, n)
			}
			return false
		}

		// 降序，按照价格高优先，时间优先  买
		if h.sort == Sort_Desc && u.Price > v.Price {
			node = h.InsertBefore(u, n)
			return false
		}

		// 升序，按照价格高优先，时间优先  卖
		if h.sort == Sort_Asc && u.Price < v.Price {
			node = h.InsertBefore(u, n)
			return false
		}
		return true
	})

	return node, nil
}
