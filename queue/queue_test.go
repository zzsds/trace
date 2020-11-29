package queue

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func ExampleQueue() {
	queue := NewQueue()
	fmt.Println(queue.Name())
}

var (
	buy  Server
	sell Server
)

func TestMain(t *testing.M) {
	buy = NewQueue(Name("BUY"))
	sell = NewQueue(Name("SELL"))
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
			buy.Push(&Node{data: NewData(push{
				Name:   fmt.Sprintf("%s-%d", t.Name(), i),
				Number: uint(i),
				Price:  1,
			})})
		}
	})
	t.Run("Unshift", func(t *testing.T) {
		for i := 0; i < 6; i++ {
			buy.Unshift(&Node{data: NewData(push{
				Name:   fmt.Sprintf("%s-%d", t.Name(), i),
				Number: uint(i),
				Price:  1,
			})})
		}
	})
	t.Run("Print", TestQueuePrint)
	t.Run("Shift", func(t *testing.T) {
		fmt.Println(*buy.Shift().data)
	})
	t.Run("Pop", func(t *testing.T) {
		fmt.Println(*buy.Pop().data)
	})
	// t.Run("Print", TestQueuePrint)
	t.Run("Reverse", func(t *testing.T) {
		buy.Reverse()
	})
	t.Log("success")

	now := time.Now().Add(3 * time.Second)
	t.Log(now)
	expire := &Node{data: NewExpireData(
		push{
			Name:   "jayden",
			Number: 100000,
			Price:  10000,
		}, &now)}
	buy.Unshift(expire)

	now = now.Add(4 * time.Second)
	buy.BeforeAdd(expire, &Node{data: NewExpireData(
		push{
			Name:   "wangjuan",
			Number: 1000,
			Price:  1000,
		}, &now)})

	t.Run("PrintExpire", TestQueuePrint)
	go func() {
		for {
			select {
			case buf := <-buy.Buffer():
				fmt.Println(buf.data.ExpireAt)
			case <-time.After(10 * time.Second):
				log.Fatal("10 超时")
			}
		}
	}()
	time.Sleep(4 * time.Second)

	t.Run("PrintDelete", TestQueuePrint)
}

func TestQueuePrint(t *testing.T) {
	head := buy.Header()
	for head != nil {
		t.Log(head.data)
		head = head.next
	}
}

func BenchmarkQueueUnshift(t *testing.B) {
	i := 0
	rand.Seed(time.Now().Unix())
	for i < t.N {
		// buy.Unshift(&Node{data: NewData(rand.Intn(10))})
		buy.Push(&Node{data: NewData(rand.Intn(10))})
		i++
	}
	t.Log(buy.Len(), t.N, "success")
}

func TestQueuePush(t *testing.T) {
	q := NewQueue()
	i := 0
	rand.Seed(time.Now().Unix())
	for i < 1000000 {
		i++
		q.Push(&Node{data: NewData(rand.Intn(8))})
	}
	i = 0
	for i < 10 {
		t.Log(*q.Get(uint(i)).data)
		i++
	}
}
