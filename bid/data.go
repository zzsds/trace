package bid

import (
	"container/list"
	"sync"
)

// DataServer ...
type DataServer interface {
	sync.Locker
	ListServer
	Init() error
	Name() string
	Sort() Sort
	Add(Unit) error
}

// Data ...
type Data struct {
	*sync.RWMutex
	*List
	sort Sort
	name string
}

// DataOption ...
type DataOption func(*Data)

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
	h.List = &List{list.New()}
	h.RWMutex = &sync.RWMutex{}
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
func (h *Data) Add(u Unit) error {
	h.Lock()
	defer h.Unlock()
	var node *Node
	if h.Len() <= 0 {
		node = h.PushFront(u)
	} else {
		for n := h.Front(); n != nil; n = n.Next() {
			v, ok := n.Value.(Unit)
			if !ok {
				break
			}
			if v.Price == u.Price {
				if v.UID == u.UID {
					v.Amount += u.Amount
					n.Value = v
					// node = n
					break
				}
				if v.CreateAt.After(u.CreateAt) {
					node = h.InsertBefore(u, n)
					break
				} else {
					node = h.InsertAfter(u, n)
					break
				}
			}

			// 降序，按照价格高优先，时间优先  买
			if h.sort == Sort_Desc && u.Price > v.Price {
				node = h.InsertBefore(u, n)
				break
			}

			// 升序，按照价格高优先，时间优先  卖
			if h.sort == Sort_Asc && u.Price < v.Price {
				node = h.InsertBefore(u, n)
				break
			}
		}
	}

	if node != nil {

	}

	return nil
}
