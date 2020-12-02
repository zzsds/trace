package queue

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

//go:generate mockgen . Unit

type Unit struct {
	Name    string
	Price   float64
	Number  uint
	UID     uint
	TradeID uint
}

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

func TestListen(t *testing.T) {
	buy.Listen(func(n *Node) error {
		now := time.Now()
		expireAt := n.Data.ExpireAt
		if expireAt != nil && expireAt.Before(now) {
			buy.WriteBuffer(n.Data.Content)
			buy.Remove(n)
		}
		return nil
	})
}

func TestPush(t *testing.T) {
	rand.Seed(time.Now().Unix())
	node := NewData(&Unit{
		Name:    "qwe",
		Number:  uint(rand.Intn(1000)),
		Price:   1.0,
		UID:     0,
		TradeID: 0,
	})
	buy.Push(node)

	t.Run("TestPrint", TestQueuePrint)
}

func TestUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	data := NewData(&Unit{
		Name:    "asd",
		Number:  uint(rand.Intn(1000)),
		Price:   2.0,
		UID:     1,
		TradeID: 1,
	})
	buy.Unshift(data)
	t.Run("TestPrint", TestQueuePrint)
}

func TestExpireUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	expire := time.Now().Add(3 * time.Second)
	node := NewExpireData(&Unit{
		Name:    "xlj",
		Number:  uint(rand.Intn(1000)),
		Price:   2.0,
		UID:     1,
		TradeID: 1,
	}, &expire)
	buy.Unshift(node)
	t.Run("TestPrint", TestQueuePrint)
}

func TestBuffer(t *testing.T) {
	go func() {
		for {
			select {
			case buf := <-buy.Buffer():
				log.Println(buf.(*Unit))
			case <-time.After(10 * time.Second):
				log.Fatal("10 超时")
			}
		}
	}()
}

func TestQueueHandle(t *testing.T) {
	t.Run("Listen", TestListen)
	t.Run("Push", TestPush)
	t.Run("Unshift", TestUnshift)
	// t.Run("Shift", func(t *testing.T) {
	// 	fmt.Println(buy.Shift().Data, buy.Shift().Data.Content)
	// })
	// t.Run("Pop", func(t *testing.T) {
	// 	fmt.Println(buy.Pop().Data, buy.Pop().Data.Content)
	// })

	t.Run("ExpireUnshift", TestExpireUnshift)

	t.Run("Buffer", TestBuffer)

	time.Sleep(5 * time.Second)

	t.Run("Print", TestQueuePrint)
}

func TestQueuePrint(t *testing.T) {
	for node := buy.Front(); node != nil; node = node.Next() {
		data := node.Data
		t.Log(data, data.Content)
	}
}

func BenchmarkQueueUnshift(t *testing.B) {
	i := 0
	rand.Seed(time.Now().Unix())
	for i < t.N {
		// buy.Unshift(&Node{Data: NewData(rand.Intn(10))})
		buy.Push(NewData(rand.Intn(10)))
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
		q.Push(NewData(rand.Intn(8)))
	}
	i = 0
	for i < 10 {
		t.Log(*q.Get(i))
		i++
	}
}
