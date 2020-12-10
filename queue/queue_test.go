package queue

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var que Server

type Unit struct {
	Name   string
	Price  float64
	Amount int
	UID    int
	ID     int
}

func TestMain(t *testing.M) {
	que = NewQueue(Name("BUY"))

	go func() {
		for {
			select {
			case buf := <-que.Buffer():
				switch buf.(type) {
				case *Unit:
					log.Println(buf.(*Unit).ID)
				}
			case <-time.After(10 * time.Second):
				log.Fatal("10 超时")
			}
		}
	}()

	t.Run()
}

func TestListen(t *testing.T) {
	que.Listen(func(n *Node) error {
		if n.Data.ExpireAt == nil {
			que.WriteBuffer(n.Data.Content)
			que.Remove(n)
		} else if n.Data.ExpireAt.Before(time.Now()) {
			n.Data.ExpireAt = nil
		}
		return nil
	})
}

func TestPush(t *testing.T) {
	rand.Seed(time.Now().Unix())
	node := NewData(&Unit{
		Name:   "qwe",
		Amount: rand.Intn(1000),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	que.Push(node)

	t.Run("TestPrint", TestQueuePrint)
}

func TestUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	data := NewData(&Unit{
		Name:   "asd",
		Amount: int(rand.Intn(1000)),
		Price:  2.0,
		UID:    1,
		ID:     1,
	})
	que.Unshift(data)
	t.Run("TestPrint", TestQueuePrint)
}
func TestList(t *testing.T) {
	for k, v := range que.List() {
		log.Println(k, v, v.Content)
	}
}

func TestExpireUnshift(t *testing.T) {
	rand.Seed(time.Now().Unix())
	expire := time.Now().Add(3 * time.Second)
	node := NewExpireData(&Unit{
		Name:   "xlj",
		Amount: int(rand.Intn(1000)),
		Price:  2.0,
		UID:    1,
		ID:     1,
	}, expire)
	que.Unshift(node)
	t.Run("TestPrint", TestQueuePrint)
}

func TestBuffer(t *testing.T) {
	go func() {
		for {
			select {
			case buf := <-que.Buffer():
				log.Println(buf.(*Unit))
			case <-time.After(10 * time.Second):
				log.Fatal("10 超时")
			}
		}
	}()
}

func TestQueueHandle(t *testing.T) {
	// t.Run("Listen", TestListen)
	t.Run("Push", TestPush)
	t.Run("Unshift", TestUnshift)

	t.Run("ExpireUnshift", TestExpireUnshift)

	// t.Run("Buffer", TestBuffer)

	time.Sleep(5 * time.Second)

	t.Run("List", TestList)
}

func TestQueuePrint(t *testing.T) {
	for node := que.Front(); node != nil; node = node.Next() {
		data := node.Data
		t.Log(data, data.Content)
	}
}

func BenchmarkQueue(t *testing.B) {
	t.Run("Unshift", BenchmarkQueueUnshift)
	t.Run("Push", BenchmarkQueuePush)
}

func BenchmarkQueueUnshift(t *testing.B) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		que.Unshift(NewData(&Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(10000000000)),
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		}))
	}
	t.Log(t.N, que.Len(), "success")
}

func BenchmarkQueuePush(t *testing.B) {
	for i := 0; i < t.N; i++ {
		price, _ := strconv.ParseFloat(strconv.Itoa(i), 64)
		que.Push(NewData(&Unit{
			Name:   "xlj",
			Amount: int(rand.Intn(10000000000)),
			Price:  price,
			UID:    int(i),
			ID:     int(i),
		}))
	}
	t.Log(t.N, que.Len(), "success")
}
