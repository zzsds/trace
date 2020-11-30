package queue

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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

func TestPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rand.Seed(time.Now().Unix())
	node := NewNode(NewData(&Unit{
		Name:    "qwe",
		Number:  uint(rand.Intn(1000)),
		Price:   1.0,
		UID:     0,
		TradeID: 0,
	}))
	if ok := buy.Push(node); !ok {
		t.Fatal("压入末尾失败")
	}
}

func TestUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	node := NewNode(NewData(&Unit{
		Name:    "asd",
		Number:  uint(rand.Intn(1000)),
		Price:   2.0,
		UID:     1,
		TradeID: 1,
	}))
	if ok := buy.Unshift(node); !ok {
		t.Fatal("开头插入失败")
	}
}

func TestExpireUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	expire := time.Now().Add(3 * time.Second)
	node := NewNode(NewExpireData(&Unit{
		Name:    "xlj",
		Number:  uint(rand.Intn(1000)),
		Price:   2.0,
		UID:     1,
		TradeID: 1,
	}, &expire))
	if ok := buy.Unshift(node); !ok {
		t.Fatal("开头插入延时队列失败")
	}
}

func TestListen(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case buf := <-buy.Buffer():
				log.Println(*buf.Data)
				wg.Done()
			case <-time.After(10 * time.Second):
				log.Fatal("10 超时")
				wg.Done()
			}
		}
	}()
	wg.Wait()
}

func TestQueueHandle(t *testing.T) {
	t.Run("Push", TestPush)
	t.Run("Print", TestQueuePrint)
	t.Run("Unshift", TestUnshift)
	t.Run("Print", TestQueuePrint)
	t.Run("Shift", func(t *testing.T) {
		fmt.Println(buy.Shift().Data, buy.Shift().Data.Content)
	})
	t.Run("Pop", func(t *testing.T) {
		fmt.Println(buy.Pop().Data, buy.Pop().Data.Content)
	})

	t.Run("ExpireUnshift", TestExpireUnshift)

	t.Run("PrintExpire", TestQueuePrint)

	t.Run("Listen", TestListen)

	time.Sleep(4 * time.Second)

	t.Run("Print", TestQueuePrint)
}

func TestQueuePrint(t *testing.T) {
	for head := buy.Header(); head != nil; head = head.next {
		data := head.Data
		t.Log(data, data.Content)
	}
}

func BenchmarkQueueUnshift(t *testing.B) {
	i := 0
	rand.Seed(time.Now().Unix())
	for i < t.N {
		// buy.Unshift(&Node{Data: NewData(rand.Intn(10))})
		buy.Push(&Node{Data: NewData(rand.Intn(10))})
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
		q.Push(&Node{Data: NewData(rand.Intn(8))})
	}
	i = 0
	for i < 10 {
		t.Log(*q.Get(uint(i)).Data)
		i++
	}
}
