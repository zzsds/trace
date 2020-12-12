package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/zzsds/trade"
	"github.com/zzsds/trade/bid"
	"github.com/zzsds/trade/queue"
)

//go:generate go run main.go
func main() {
	t := trade.Newtrade(func(o *trade.Options) {
		o.Name = "New Product"
	})

	log.Println(t.Name())

	t.RegisterBid(1, bid.NewBid(bid.Name("product")))
	b, _ := t.LoadBid(1)

	data := queue.NewData(&bid.Unit{
		Name:   "xlj",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	b.Buy().Push(data)

	data = queue.NewData(&bid.Unit{
		Name:   "qwe",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	b.Sell().Push(data)
	// fmt.Println(b)
	t.Run()
}

func queueTest() {
	que := queue.NewQueue(queue.Name("Buy"))
	log.Println(que.Name(), que.Len())

	que.Listen(func(n *queue.Node) error {
		if n.Data().ExpireAt == nil {
			que.WriteBuffer(*n.Data())
			que.Remove(n)
		} else if n.Data().ExpireAt.Before(time.Now()) {
			n.Data().ExpireAt = nil
		}
		return nil
	})

	go func() {
		for {
			select {
			case buff := <-que.Buffer():
				log.Println(buff, "出")
			}
		}
	}()

	rand.Seed(time.Now().Unix())
	data := queue.NewExpireData(&bid.Unit{
		Name:   "qwe",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	}, time.Now().Add(3*time.Second))
	que.Push(data)

	time.Sleep(5 * time.Second)
	data = queue.NewData(&bid.Unit{
		Name:   "xlj",
		Amount: int(rand.Intn(1000)),
		Price:  1.0,
		UID:    0,
		ID:     0,
	})
	que.Unshift(data)

	time.Sleep(2 * time.Second)
	fmt.Println("End")
	for k, v := range que.List() {
		log.Println(k, v.Content, "list")
	}
}
