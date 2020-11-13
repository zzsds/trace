package queue

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueue(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	buy  Server
	sell Server
)

func TestMain(t *testing.M) {
	buy = NewQueue(func(o *Options) {
		o.Name = "BUY"
	})
	sell = NewQueue(func(o *Options) {
		o.Name = "SELL"
	})
	t.Run()
}

func TestQueueHandle(t *testing.T) {
	type push struct {
		Name   string
		Price  float64
		Number uint
	}
	t.Run("Push", func(t *testing.T) {
		rand.Seed(time.Now().Unix())
		for i := 0; i < 3; i++ {
			buy.Push(&Node{Data: NewData(push{
				Name:   fmt.Sprintf("%s-%d", t.Name(), i),
				Number: uint(i),
				Price:  1,
			})})
		}
	})
	t.Run("Unshift", func(t *testing.T) {
		for i := 0; i < 6; i++ {
			buy.Unshift(&Node{Data: NewData(push{
				Name:   fmt.Sprintf("%s-%d", t.Name(), i),
				Number: uint(i),
				Price:  1,
			})})
		}
	})
	t.Run("Print", TestQueuePrint)
	t.Run("Shift", func(t *testing.T) {
		fmt.Println(*buy.Shift().Data)
	})
	t.Run("Pop", func(t *testing.T) {
		fmt.Println(*buy.Pop().Data)
	})
	// t.Run("Print", TestQueuePrint)
	t.Run("Reverse", func(t *testing.T) {
		buy.Reverse()
	})
	t.Log("success")
}

func TestQueuePrint(t *testing.T) {
	head := buy.Header()
	for head != nil {
		t.Log(head.Data, head.Sort)
		head = head.Next
	}
}

func BenchmarkQueueUnshift(t *testing.B) {
	i := 0
	rand.Seed(time.Now().Unix())
	for i < t.N {
		buy.Unshift(&Node{Data: NewData(rand.Intn(10))})
		// buy.Push(&Node{Data: NewData(rand.Intn(10))})
		i++
	}

	t.Log(buy.Length(), "success")
}

func TestQueuePush(t *testing.T) {
	q := NewQueue()
	i := 0
	rand.Seed(time.Now().Unix())
	for i < 100000000 {
		i++
		q.Push(&Node{Data: NewData(rand.Intn(8))})
	}
	i = 0
	for i < 10 {
		t.Log(*q.Get(uint(i)).Data)
		i++
	}
}

type account struct {
	name  string
	age   int
	total float64
}

func TestArrayPointer(t *testing.T) {
	arr := make([]*account, 0, 20)
	for i := 0; i < 20; i++ {
		arr = append(arr, &account{
			name:  strconv.Itoa(i),
			age:   i,
			total: 0.25,
		})
	}
	index := 0
	demo := &account{}
	for k, v := range arr {
		if k == 5 {
			demo = v
			index = k
			continue
		}
	}
	demo.name = "jayden"
	arr[index].total -= 0.02

	for _, v := range arr {
		fmt.Println(*v)
	}
}
